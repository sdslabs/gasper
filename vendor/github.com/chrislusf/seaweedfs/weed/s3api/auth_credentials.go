package s3api

import (
	"bytes"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/s3api/s3err"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/iam_pb"
)

type Action string

const (
	ACTION_READ  = "Read"
	ACTION_WRITE = "Write"
	ACTION_ADMIN = "Admin"
)

type Iam interface {
	Check(f http.HandlerFunc, actions ...Action) http.HandlerFunc
}

type IdentityAccessManagement struct {
	identities []*Identity
	domain     string
}

type Identity struct {
	Name        string
	Credentials []*Credential
	Actions     []Action
}

type Credential struct {
	AccessKey string
	SecretKey string
}

func NewIdentityAccessManagement(fileName string, domain string) *IdentityAccessManagement {
	iam := &IdentityAccessManagement{
		domain: domain,
	}
	if fileName == "" {
		return iam
	}
	if err := iam.loadS3ApiConfiguration(fileName); err != nil {
		glog.Fatalf("fail to load config file %s: %v", fileName, err)
	}
	return iam
}

func (iam *IdentityAccessManagement) loadS3ApiConfiguration(fileName string) error {

	s3ApiConfiguration := &iam_pb.S3ApiConfiguration{}

	rawData, readErr := ioutil.ReadFile(fileName)
	if readErr != nil {
		glog.Warningf("fail to read %s : %v", fileName, readErr)
		return fmt.Errorf("fail to read %s : %v", fileName, readErr)
	}

	glog.V(1).Infof("load s3 config: %v", fileName)
	if err := jsonpb.Unmarshal(bytes.NewReader(rawData), s3ApiConfiguration); err != nil {
		glog.Warningf("unmarshal error: %v", err)
		return fmt.Errorf("unmarshal %s error: %v", fileName, err)
	}

	for _, ident := range s3ApiConfiguration.Identities {
		t := &Identity{
			Name:        ident.Name,
			Credentials: nil,
			Actions:     nil,
		}
		for _, action := range ident.Actions {
			t.Actions = append(t.Actions, Action(action))
		}
		for _, cred := range ident.Credentials {
			t.Credentials = append(t.Credentials, &Credential{
				AccessKey: cred.AccessKey,
				SecretKey: cred.SecretKey,
			})
		}
		iam.identities = append(iam.identities, t)
	}

	return nil
}

func (iam *IdentityAccessManagement) isEnabled() bool {

	return len(iam.identities) > 0
}

func (iam *IdentityAccessManagement) lookupByAccessKey(accessKey string) (identity *Identity, cred *Credential, found bool) {

	for _, ident := range iam.identities {
		for _, cred := range ident.Credentials {
			if cred.AccessKey == accessKey {
				return ident, cred, true
			}
		}
	}
	return nil, nil, false
}

func (iam *IdentityAccessManagement) lookupAnonymous() (identity *Identity, found bool) {

	for _, ident := range iam.identities {
		if ident.Name == "anonymous" {
			return ident, true
		}
	}
	return nil, false
}

func (iam *IdentityAccessManagement) Auth(f http.HandlerFunc, action Action) http.HandlerFunc {

	if !iam.isEnabled() {
		return f
	}

	return func(w http.ResponseWriter, r *http.Request) {
		errCode := iam.authRequest(r, action)
		if errCode == s3err.ErrNone {
			f(w, r)
			return
		}
		writeErrorResponse(w, errCode, r.URL)
	}
}

// check whether the request has valid access keys
func (iam *IdentityAccessManagement) authRequest(r *http.Request, action Action) s3err.ErrorCode {
	var identity *Identity
	var s3Err s3err.ErrorCode
	var found bool
	switch getRequestAuthType(r) {
	case authTypeStreamingSigned:
		return s3err.ErrNone
	case authTypeUnknown:
		glog.V(3).Infof("unknown auth type")
		return s3err.ErrAccessDenied
	case authTypePresignedV2, authTypeSignedV2:
		glog.V(3).Infof("v2 auth type")
		identity, s3Err = iam.isReqAuthenticatedV2(r)
	case authTypeSigned, authTypePresigned:
		glog.V(3).Infof("v4 auth type")
		identity, s3Err = iam.reqSignatureV4Verify(r)
	case authTypePostPolicy:
		glog.V(3).Infof("post policy auth type")
		return s3err.ErrNone
	case authTypeJWT:
		glog.V(3).Infof("jwt auth type")
		return s3err.ErrNotImplemented
	case authTypeAnonymous:
		identity, found = iam.lookupAnonymous()
		if !found {
			return s3err.ErrAccessDenied
		}
	default:
		return s3err.ErrNotImplemented
	}

	glog.V(3).Infof("auth error: %v", s3Err)
	if s3Err != s3err.ErrNone {
		return s3Err
	}

	glog.V(3).Infof("user name: %v actions: %v", identity.Name, identity.Actions)

	bucket, _ := getBucketAndObject(r)

	if !identity.canDo(action, bucket) {
		return s3err.ErrAccessDenied
	}

	return s3err.ErrNone

}

func (identity *Identity) canDo(action Action, bucket string) bool {
	for _, a := range identity.Actions {
		if a == "Admin" {
			return true
		}
	}
	for _, a := range identity.Actions {
		if a == action {
			return true
		}
	}
	if bucket == "" {
		return false
	}
	limitedByBucket := string(action) + ":" + bucket
	for _, a := range identity.Actions {
		if string(a) == limitedByBucket {
			return true
		}
	}
	return false
}
