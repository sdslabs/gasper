package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/security"
	"github.com/chrislusf/seaweedfs/weed/util"
)

type UploadResult struct {
	Name       string `json:"name,omitempty"`
	Size       uint32 `json:"size,omitempty"`
	Error      string `json:"error,omitempty"`
	ETag       string `json:"eTag,omitempty"`
	CipherKey  []byte `json:"cipherKey,omitempty"`
	Mime       string `json:"mime,omitempty"`
	Gzip       uint32 `json:"gzip,omitempty"`
	ContentMd5 string `json:"contentMd5,omitempty"`
}

func (uploadResult *UploadResult) ToPbFileChunk(fileId string, offset int64) *filer_pb.FileChunk {
	fid, _ := filer_pb.ToFileIdObject(fileId)
	return &filer_pb.FileChunk{
		FileId:       fileId,
		Offset:       offset,
		Size:         uint64(uploadResult.Size),
		Mtime:        time.Now().UnixNano(),
		ETag:         uploadResult.ETag,
		CipherKey:    uploadResult.CipherKey,
		IsCompressed: uploadResult.Gzip > 0,
		Fid:          fid,
	}
}

// HTTPClient interface for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient HTTPClient
)

func init() {
	HttpClient = &http.Client{Transport: &http.Transport{
		MaxIdleConnsPerHost: 1024,
	}}
}

var fileNameEscaper = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")

// Upload sends a POST request to a volume server to upload the content with adjustable compression level
func UploadData(uploadUrl string, filename string, cipher bool, data []byte, isInputCompressed bool, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (uploadResult *UploadResult, err error) {
	uploadResult, err = retriedUploadData(uploadUrl, filename, cipher, data, isInputCompressed, mtype, pairMap, jwt)
	return
}

// Upload sends a POST request to a volume server to upload the content with fast compression
func Upload(uploadUrl string, filename string, cipher bool, reader io.Reader, isInputCompressed bool, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (uploadResult *UploadResult, err error, data []byte) {
	uploadResult, err, data = doUpload(uploadUrl, filename, cipher, reader, isInputCompressed, mtype, pairMap, jwt)
	return
}

func doUpload(uploadUrl string, filename string, cipher bool, reader io.Reader, isInputCompressed bool, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (uploadResult *UploadResult, err error, data []byte) {
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("read input: %v", err)
		return
	}
	uploadResult, uploadErr := retriedUploadData(uploadUrl, filename, cipher, data, isInputCompressed, mtype, pairMap, jwt)
	return uploadResult, uploadErr, data
}

func retriedUploadData(uploadUrl string, filename string, cipher bool, data []byte, isInputCompressed bool, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (uploadResult *UploadResult, err error) {
	for i := 0; i < 1; i++ {
		uploadResult, err = doUploadData(uploadUrl, filename, cipher, data, isInputCompressed, mtype, pairMap, jwt)
		if err == nil {
			return
		} else {
			glog.Warningf("uploading to %s: %v", uploadUrl, err)
		}
	}
	return
}

func doUploadData(uploadUrl string, filename string, cipher bool, data []byte, isInputCompressed bool, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (uploadResult *UploadResult, err error) {
	contentIsGzipped := isInputCompressed
	shouldGzipNow := false
	if !isInputCompressed {
		if mtype == "" {
			mtype = http.DetectContentType(data)
			// println("detect1 mimetype to", mtype)
			if mtype == "application/octet-stream" {
				mtype = ""
			}
		}
		if shouldBeCompressed, iAmSure := util.IsCompressableFileType(filepath.Base(filename), mtype); iAmSure && shouldBeCompressed {
			shouldGzipNow = true
		} else if !iAmSure && mtype == "" && len(data) > 128 {
			var compressed []byte
			compressed, err = util.GzipData(data[0:128])
			shouldGzipNow = len(compressed)*10 < 128*9 // can not compress to less than 90%
		}
	}

	var clearDataLen int

	// gzip if possible
	// this could be double copying
	clearDataLen = len(data)
	if shouldGzipNow {
		compressed, compressErr := util.GzipData(data)
		// fmt.Printf("data is compressed from %d ==> %d\n", len(data), len(compressed))
		if compressErr == nil {
			data = compressed
			contentIsGzipped = true
		}
	} else if isInputCompressed {
		// just to get the clear data length
		clearData, err := util.DecompressData(data)
		if err == nil {
			clearDataLen = len(clearData)
		}
	}

	if cipher {
		// encrypt(gzip(data))

		// encrypt
		cipherKey := util.GenCipherKey()
		encryptedData, encryptionErr := util.Encrypt(data, cipherKey)
		if encryptionErr != nil {
			err = fmt.Errorf("encrypt input: %v", encryptionErr)
			return
		}

		// upload data
		uploadResult, err = upload_content(uploadUrl, func(w io.Writer) (err error) {
			_, err = w.Write(encryptedData)
			return
		}, "", false, len(encryptedData), "", nil, jwt)
		if uploadResult != nil {
			uploadResult.Name = filename
			uploadResult.Mime = mtype
			uploadResult.CipherKey = cipherKey
		}
	} else {
		// upload data
		uploadResult, err = upload_content(uploadUrl, func(w io.Writer) (err error) {
			_, err = w.Write(data)
			return
		}, filename, contentIsGzipped, 0, mtype, pairMap, jwt)
	}

	if uploadResult == nil {
		return
	}

	uploadResult.Size = uint32(clearDataLen)
	if contentIsGzipped {
		uploadResult.Gzip = 1
	}

	return uploadResult, err
}

func upload_content(uploadUrl string, fillBufferFunction func(w io.Writer) error, filename string, isGzipped bool, originalDataSize int, mtype string, pairMap map[string]string, jwt security.EncodedJwt) (*UploadResult, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileNameEscaper.Replace(filename)))
	if mtype == "" {
		mtype = mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))
	}
	if mtype != "" {
		h.Set("Content-Type", mtype)
	}
	if isGzipped {
		h.Set("Content-Encoding", "gzip")
	}

	file_writer, cp_err := body_writer.CreatePart(h)
	if cp_err != nil {
		glog.V(0).Infoln("error creating form file", cp_err.Error())
		return nil, cp_err
	}
	if err := fillBufferFunction(file_writer); err != nil {
		glog.V(0).Infoln("error copying data", err)
		return nil, err
	}
	content_type := body_writer.FormDataContentType()
	if err := body_writer.Close(); err != nil {
		glog.V(0).Infoln("error closing body", err)
		return nil, err
	}

	req, postErr := http.NewRequest("POST", uploadUrl, body_buf)
	if postErr != nil {
		glog.V(1).Infof("create upload request %s: %v", uploadUrl, postErr)
		return nil, fmt.Errorf("create upload request %s: %v", uploadUrl, postErr)
	}
	req.Header.Set("Content-Type", content_type)
	for k, v := range pairMap {
		req.Header.Set(k, v)
	}
	if jwt != "" {
		req.Header.Set("Authorization", "BEARER "+string(jwt))
	}
	resp, post_err := HttpClient.Do(req)
	if post_err != nil {
		glog.Errorf("upload %s %d bytes to %v: %v", filename, originalDataSize, uploadUrl, post_err)
		debug.PrintStack()
		return nil, fmt.Errorf("upload %s %d bytes to %v: %v", filename, originalDataSize, uploadUrl, post_err)
	}
	defer util.CloseResponse(resp)

	var ret UploadResult
	etag := getEtag(resp)
	if resp.StatusCode == http.StatusNoContent {
		ret.ETag = etag
		return &ret, nil
	}

	resp_body, ra_err := ioutil.ReadAll(resp.Body)
	if ra_err != nil {
		return nil, fmt.Errorf("read response body %v: %v", uploadUrl, ra_err)
	}

	unmarshal_err := json.Unmarshal(resp_body, &ret)
	if unmarshal_err != nil {
		glog.Errorf("unmarshal %s: %v", uploadUrl, string(resp_body))
		return nil, fmt.Errorf("unmarshal %v: %v", uploadUrl, unmarshal_err)
	}
	if ret.Error != "" {
		return nil, fmt.Errorf("unmarshalled error %v: %v", uploadUrl, ret.Error)
	}
	ret.ETag = etag
	ret.ContentMd5 = resp.Header.Get("Content-MD5")
	return &ret, nil
}

func getEtag(r *http.Response) (etag string) {
	etag = r.Header.Get("ETag")
	if strings.HasPrefix(etag, "\"") && strings.HasSuffix(etag, "\"") {
		etag = etag[1 : len(etag)-1]
	}
	return
}
