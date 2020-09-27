package weed_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/util"
)

// handling single chunk POST or PUT upload
func (fs *FilerServer) encrypt(ctx context.Context, w http.ResponseWriter, r *http.Request,
	replication string, collection string, dataCenter string, ttlSeconds int32, ttlString string, fsync bool) (filerResult *FilerPostResult, err error) {

	fileId, urlLocation, auth, err := fs.assignNewFileInfo(replication, collection, dataCenter, ttlString, fsync)

	if err != nil || fileId == "" || urlLocation == "" {
		return nil, fmt.Errorf("fail to allocate volume for %s, collection:%s, datacenter:%s", r.URL.Path, collection, dataCenter)
	}

	glog.V(4).Infof("write %s to %v", r.URL.Path, urlLocation)

	// Note: encrypt(gzip(data)), encrypt data first, then gzip

	sizeLimit := int64(fs.option.MaxMB) * 1024 * 1024

	pu, err := needle.ParseUpload(r, sizeLimit)
	uncompressedData := pu.Data
	if pu.IsGzipped {
		uncompressedData = pu.UncompressedData
	}
	if pu.MimeType == "" {
		pu.MimeType = http.DetectContentType(uncompressedData)
		// println("detect2 mimetype to", pu.MimeType)
	}

	uploadResult, uploadError := operation.UploadData(urlLocation, pu.FileName, true, uncompressedData, false, pu.MimeType, pu.PairMap, auth)
	if uploadError != nil {
		return nil, fmt.Errorf("upload to volume server: %v", uploadError)
	}

	// Save to chunk manifest structure
	fileChunks := []*filer_pb.FileChunk{uploadResult.ToPbFileChunk(fileId, 0)}

	// fmt.Printf("uploaded: %+v\n", uploadResult)

	path := r.URL.Path
	if strings.HasSuffix(path, "/") {
		if pu.FileName != "" {
			path += pu.FileName
		}
	}

	entry := &filer.Entry{
		FullPath: util.FullPath(path),
		Attr: filer.Attr{
			Mtime:       time.Now(),
			Crtime:      time.Now(),
			Mode:        0660,
			Uid:         OS_UID,
			Gid:         OS_GID,
			Replication: replication,
			Collection:  collection,
			TtlSec:      ttlSeconds,
			Mime:        pu.MimeType,
			Md5:         util.Base64Md5ToBytes(pu.ContentMd5),
		},
		Chunks: fileChunks,
	}

	filerResult = &FilerPostResult{
		Name: pu.FileName,
		Size: int64(pu.OriginalDataSize),
	}

	if dbErr := fs.filer.CreateEntry(ctx, entry, false, false, nil); dbErr != nil {
		fs.filer.DeleteChunks(entry.Chunks)
		err = dbErr
		filerResult.Error = dbErr.Error()
		return
	}

	return
}
