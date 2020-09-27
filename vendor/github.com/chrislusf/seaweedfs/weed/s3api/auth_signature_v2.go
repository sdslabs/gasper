/*
 * The following code tries to reverse engineer the Amazon S3 APIs,
 * and is mostly copied from minio implementation.
 */

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package s3api

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/s3api/s3err"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Whitelist resource list that will be used in query string for signature-V2 calculation.
// The list should be alphabetically sorted
var resourceList = []string{
	"acl",
	"delete",
	"lifecycle",
	"location",
	"logging",
	"notification",
	"partNumber",
	"policy",
	"requestPayment",
	"response-cache-control",
	"response-content-disposition",
	"response-content-encoding",
	"response-content-language",
	"response-content-type",
	"response-expires",
	"torrent",
	"uploadId",
	"uploads",
	"versionId",
	"versioning",
	"versions",
	"website",
}

// Verify if request has valid AWS Signature Version '2'.
func (iam *IdentityAccessManagement) isReqAuthenticatedV2(r *http.Request) (*Identity, s3err.ErrorCode) {
	if isRequestSignatureV2(r) {
		return iam.doesSignV2Match(r)
	}
	return iam.doesPresignV2SignatureMatch(r)
}

func (iam *IdentityAccessManagement) doesPolicySignatureV2Match(formValues http.Header) s3err.ErrorCode {
	accessKey := formValues.Get("AWSAccessKeyId")
	_, cred, found := iam.lookupByAccessKey(accessKey)
	if !found {
		return s3err.ErrInvalidAccessKeyID
	}
	policy := formValues.Get("Policy")
	signature := formValues.Get("Signature")
	if !compareSignatureV2(signature, calculateSignatureV2(policy, cred.SecretKey)) {
		return s3err.ErrSignatureDoesNotMatch
	}
	return s3err.ErrNone
}

// Authorization = "AWS" + " " + AWSAccessKeyId + ":" + Signature;
// Signature = Base64( HMAC-SHA1( YourSecretKey, UTF-8-Encoding-Of( StringToSign ) ) );
//
// StringToSign = HTTP-Verb + "\n" +
//  	Content-Md5 + "\n" +
//  	Content-Type + "\n" +
//  	Date + "\n" +
//  	CanonicalizedProtocolHeaders +
//  	CanonicalizedResource;
//
// CanonicalizedResource = [ "/" + Bucket ] +
//  	<HTTP-Request-URI, from the protocol name up to the query string> +
//  	[ subresource, if present. For example "?acl", "?location", "?logging", or "?torrent"];
//
// CanonicalizedProtocolHeaders = <described below>

// doesSignV2Match - Verify authorization header with calculated header in accordance with
//     - http://docs.aws.amazon.com/AmazonS3/latest/dev/auth-request-sig-v2.html
// returns true if matches, false otherwise. if error is not nil then it is always false

func validateV2AuthHeader(v2Auth string) (accessKey string, errCode s3err.ErrorCode) {
	if v2Auth == "" {
		return "", s3err.ErrAuthHeaderEmpty
	}
	// Verify if the header algorithm is supported or not.
	if !strings.HasPrefix(v2Auth, signV2Algorithm) {
		return "", s3err.ErrSignatureVersionNotSupported
	}

	// below is V2 Signed Auth header format, splitting on `space` (after the `AWS` string).
	// Authorization = "AWS" + " " + AWSAccessKeyId + ":" + Signature
	authFields := strings.Split(v2Auth, " ")
	if len(authFields) != 2 {
		return "", s3err.ErrMissingFields
	}

	// Then will be splitting on ":", this will seprate `AWSAccessKeyId` and `Signature` string.
	keySignFields := strings.Split(strings.TrimSpace(authFields[1]), ":")
	if len(keySignFields) != 2 {
		return "", s3err.ErrMissingFields
	}

	return keySignFields[0], s3err.ErrNone
}

func (iam *IdentityAccessManagement) doesSignV2Match(r *http.Request) (*Identity, s3err.ErrorCode) {
	v2Auth := r.Header.Get("Authorization")

	accessKey, apiError := validateV2AuthHeader(v2Auth)
	if apiError != s3err.ErrNone {
		return nil, apiError
	}

	// Access credentials.
	// Validate if access key id same.
	ident, cred, found := iam.lookupByAccessKey(accessKey)
	if !found {
		return nil, s3err.ErrInvalidAccessKeyID
	}

	// r.RequestURI will have raw encoded URI as sent by the client.
	tokens := strings.SplitN(r.RequestURI, "?", 2)
	encodedResource := tokens[0]
	encodedQuery := ""
	if len(tokens) == 2 {
		encodedQuery = tokens[1]
	}

	unescapedQueries, err := unescapeQueries(encodedQuery)
	if err != nil {
		return nil, s3err.ErrInvalidQueryParams
	}

	encodedResource, err = getResource(encodedResource, r.Host, iam.domain)
	if err != nil {
		return nil, s3err.ErrInvalidRequest
	}

	prefix := fmt.Sprintf("%s %s:", signV2Algorithm, cred.AccessKey)
	if !strings.HasPrefix(v2Auth, prefix) {
		return nil, s3err.ErrSignatureDoesNotMatch
	}
	v2Auth = v2Auth[len(prefix):]
	expectedAuth := signatureV2(cred, r.Method, encodedResource, strings.Join(unescapedQueries, "&"), r.Header)
	if !compareSignatureV2(v2Auth, expectedAuth) {
		return nil, s3err.ErrSignatureDoesNotMatch
	}
	return ident, s3err.ErrNone
}

// doesPresignV2SignatureMatch - Verify query headers with presigned signature
//     - http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html#RESTAuthenticationQueryStringAuth
// returns ErrNone if matches. S3 errors otherwise.
func (iam *IdentityAccessManagement) doesPresignV2SignatureMatch(r *http.Request) (*Identity, s3err.ErrorCode) {

	// r.RequestURI will have raw encoded URI as sent by the client.
	tokens := strings.SplitN(r.RequestURI, "?", 2)
	encodedResource := tokens[0]
	encodedQuery := ""
	if len(tokens) == 2 {
		encodedQuery = tokens[1]
	}

	var (
		filteredQueries []string
		gotSignature    string
		expires         string
		accessKey       string
		err             error
	)

	var unescapedQueries []string
	unescapedQueries, err = unescapeQueries(encodedQuery)
	if err != nil {
		return nil, s3err.ErrInvalidQueryParams
	}

	// Extract the necessary values from presigned query, construct a list of new filtered queries.
	for _, query := range unescapedQueries {
		keyval := strings.SplitN(query, "=", 2)
		if len(keyval) != 2 {
			return nil, s3err.ErrInvalidQueryParams
		}
		switch keyval[0] {
		case "AWSAccessKeyId":
			accessKey = keyval[1]
		case "Signature":
			gotSignature = keyval[1]
		case "Expires":
			expires = keyval[1]
		default:
			filteredQueries = append(filteredQueries, query)
		}
	}

	// Invalid values returns error.
	if accessKey == "" || gotSignature == "" || expires == "" {
		return nil, s3err.ErrInvalidQueryParams
	}

	// Validate if access key id same.
	ident, cred, found := iam.lookupByAccessKey(accessKey)
	if !found {
		return nil, s3err.ErrInvalidAccessKeyID
	}

	// Make sure the request has not expired.
	expiresInt, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return nil, s3err.ErrMalformedExpires
	}

	// Check if the presigned URL has expired.
	if expiresInt < time.Now().UTC().Unix() {
		return nil, s3err.ErrExpiredPresignRequest
	}

	encodedResource, err = getResource(encodedResource, r.Host, iam.domain)
	if err != nil {
		return nil, s3err.ErrInvalidRequest
	}

	expectedSignature := preSignatureV2(cred, r.Method, encodedResource, strings.Join(filteredQueries, "&"), r.Header, expires)
	if !compareSignatureV2(gotSignature, expectedSignature) {
		return nil, s3err.ErrSignatureDoesNotMatch
	}

	return ident, s3err.ErrNone
}

// Escape encodedQuery string into unescaped list of query params, returns error
// if any while unescaping the values.
func unescapeQueries(encodedQuery string) (unescapedQueries []string, err error) {
	for _, query := range strings.Split(encodedQuery, "&") {
		var unescapedQuery string
		unescapedQuery, err = url.QueryUnescape(query)
		if err != nil {
			return nil, err
		}
		unescapedQueries = append(unescapedQueries, unescapedQuery)
	}
	return unescapedQueries, nil
}

// Returns "/bucketName/objectName" for path-style or virtual-host-style requests.
func getResource(path string, host string, domain string) (string, error) {
	if domain == "" {
		return path, nil
	}
	// If virtual-host-style is enabled construct the "resource" properly.
	if strings.Contains(host, ":") {
		// In bucket.mydomain.com:9000, strip out :9000
		var err error
		if host, _, err = net.SplitHostPort(host); err != nil {
			return "", err
		}
	}
	if !strings.HasSuffix(host, "."+domain) {
		return path, nil
	}
	bucket := strings.TrimSuffix(host, "."+domain)
	return "/" + pathJoin(bucket, path), nil
}

// pathJoin - like path.Join() but retains trailing "/" of the last element
func pathJoin(elem ...string) string {
	trailingSlash := ""
	if len(elem) > 0 {
		if strings.HasSuffix(elem[len(elem)-1], "/") {
			trailingSlash = "/"
		}
	}
	return path.Join(elem...) + trailingSlash
}

// Return the signature v2 of a given request.
func signatureV2(cred *Credential, method string, encodedResource string, encodedQuery string, headers http.Header) string {
	stringToSign := getStringToSignV2(method, encodedResource, encodedQuery, headers, "")
	signature := calculateSignatureV2(stringToSign, cred.SecretKey)
	return signature
}

// Return string to sign under two different conditions.
// - if expires string is set then string to sign includes date instead of the Date header.
// - if expires string is empty then string to sign includes date header instead.
func getStringToSignV2(method string, encodedResource, encodedQuery string, headers http.Header, expires string) string {
	canonicalHeaders := canonicalizedAmzHeadersV2(headers)
	if len(canonicalHeaders) > 0 {
		canonicalHeaders += "\n"
	}

	date := expires // Date is set to expires date for presign operations.
	if date == "" {
		// If expires date is empty then request header Date is used.
		date = headers.Get("Date")
	}

	// From the Amazon docs:
	//
	// StringToSign = HTTP-Verb + "\n" +
	// 	 Content-Md5 + "\n" +
	//	 Content-Type + "\n" +
	//	 Date/Expires + "\n" +
	//	 CanonicalizedProtocolHeaders +
	//	 CanonicalizedResource;
	stringToSign := strings.Join([]string{
		method,
		headers.Get("Content-MD5"),
		headers.Get("Content-Type"),
		date,
		canonicalHeaders,
	}, "\n")

	return stringToSign + canonicalizedResourceV2(encodedResource, encodedQuery)
}

// Return canonical resource string.
func canonicalizedResourceV2(encodedResource, encodedQuery string) string {
	queries := strings.Split(encodedQuery, "&")
	keyval := make(map[string]string)
	for _, query := range queries {
		key := query
		val := ""
		index := strings.Index(query, "=")
		if index != -1 {
			key = query[:index]
			val = query[index+1:]
		}
		keyval[key] = val
	}

	var canonicalQueries []string
	for _, key := range resourceList {
		val, ok := keyval[key]
		if !ok {
			continue
		}
		if val == "" {
			canonicalQueries = append(canonicalQueries, key)
			continue
		}
		canonicalQueries = append(canonicalQueries, key+"="+val)
	}

	// The queries will be already sorted as resourceList is sorted, if canonicalQueries
	// is empty strings.Join returns empty.
	canonicalQuery := strings.Join(canonicalQueries, "&")
	if canonicalQuery != "" {
		return encodedResource + "?" + canonicalQuery
	}
	return encodedResource
}

// Return canonical headers.
func canonicalizedAmzHeadersV2(headers http.Header) string {
	var keys []string
	keyval := make(map[string]string)
	for key := range headers {
		lkey := strings.ToLower(key)
		if !strings.HasPrefix(lkey, "x-amz-") {
			continue
		}
		keys = append(keys, lkey)
		keyval[lkey] = strings.Join(headers[key], ",")
	}
	sort.Strings(keys)
	var canonicalHeaders []string
	for _, key := range keys {
		canonicalHeaders = append(canonicalHeaders, key+":"+keyval[key])
	}
	return strings.Join(canonicalHeaders, "\n")
}

func calculateSignatureV2(stringToSign string, secret string) string {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

// compareSignatureV2 returns true if and only if both signatures
// are equal. The signatures are expected to be base64 encoded strings
// according to the AWS S3 signature V2 spec.
func compareSignatureV2(sig1, sig2 string) bool {
	// Decode signature string to binary byte-sequence representation is required
	// as Base64 encoding of a value is not unique:
	// For example "aGVsbG8=" and "aGVsbG8=\r" will result in the same byte slice.
	signature1, err := base64.StdEncoding.DecodeString(sig1)
	if err != nil {
		return false
	}
	signature2, err := base64.StdEncoding.DecodeString(sig2)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(signature1, signature2) == 1
}

// Return signature-v2 for the presigned request.
func preSignatureV2(cred *Credential, method string, encodedResource string, encodedQuery string, headers http.Header, expires string) string {
	stringToSign := getStringToSignV2(method, encodedResource, encodedQuery, headers, expires)
	return calculateSignatureV2(stringToSign, cred.SecretKey)
}
