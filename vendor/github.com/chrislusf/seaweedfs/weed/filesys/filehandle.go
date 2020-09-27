package filesys

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/seaweedfs/fuse"
	"github.com/seaweedfs/fuse/fs"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
)

type FileHandle struct {
	// cache file has been written to
	dirtyPages  *ContinuousDirtyPages
	contentType string
	handle      uint64
	sync.RWMutex

	f         *File
	RequestId fuse.RequestID // unique ID for request
	NodeId    fuse.NodeID    // file or directory the request is about
	Uid       uint32         // user ID of process making request
	Gid       uint32         // group ID of process making request

}

func newFileHandle(file *File, uid, gid uint32) *FileHandle {
	fh := &FileHandle{
		f:          file,
		dirtyPages: newDirtyPages(file),
		Uid:        uid,
		Gid:        gid,
	}
	if fh.f.entry != nil {
		fh.f.entry.Attributes.FileSize = filer.FileSize(fh.f.entry)
	}

	return fh
}

var _ = fs.Handle(&FileHandle{})

// var _ = fs.HandleReadAller(&FileHandle{})
var _ = fs.HandleReader(&FileHandle{})
var _ = fs.HandleFlusher(&FileHandle{})
var _ = fs.HandleWriter(&FileHandle{})
var _ = fs.HandleReleaser(&FileHandle{})

func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {

	glog.V(4).Infof("%s read fh %d: [%d,%d) size %d resp.Data cap=%d", fh.f.fullpath(), fh.handle, req.Offset, req.Offset+int64(req.Size), req.Size, cap(resp.Data))
	fh.RLock()
	defer fh.RUnlock()

	if req.Size <= 0 {
		return nil
	}

	buff := resp.Data[:cap(resp.Data)]
	if req.Size > cap(resp.Data) {
		// should not happen
		buff = make([]byte, req.Size)
	}

	totalRead, err := fh.readFromChunks(buff, req.Offset)
	if err == nil {
		maxStop := fh.readFromDirtyPages(buff, req.Offset)
		totalRead = max(maxStop-req.Offset, totalRead)
	}

	if err != nil {
		glog.Warningf("file handle read %s %d: %v", fh.f.fullpath(), totalRead, err)
		return nil
	}

	if totalRead > int64(len(buff)) {
		glog.Warningf("%s FileHandle Read %d: [%d,%d) size %d totalRead %d", fh.f.fullpath(), fh.handle, req.Offset, req.Offset+int64(req.Size), req.Size, totalRead)
		totalRead = min(int64(len(buff)), totalRead)
	}
	// resp.Data = buff[:totalRead]
	resp.Data = buff

	return err
}

func (fh *FileHandle) readFromDirtyPages(buff []byte, startOffset int64) (maxStop int64) {
	maxStop = fh.dirtyPages.ReadDirtyDataAt(buff, startOffset)
	return
}

func (fh *FileHandle) readFromChunks(buff []byte, offset int64) (int64, error) {

	fileSize := int64(filer.FileSize(fh.f.entry))

	if fileSize == 0 {
		glog.V(1).Infof("empty fh %v", fh.f.fullpath())
		return 0, io.EOF
	}

	var chunkResolveErr error
	if fh.f.entryViewCache == nil {
		fh.f.entryViewCache, chunkResolveErr = filer.NonOverlappingVisibleIntervals(filer.LookupFn(fh.f.wfs), fh.f.entry.Chunks)
		if chunkResolveErr != nil {
			return 0, fmt.Errorf("fail to resolve chunk manifest: %v", chunkResolveErr)
		}
		fh.f.reader = nil
	}

	if fh.f.reader == nil {
		chunkViews := filer.ViewFromVisibleIntervals(fh.f.entryViewCache, 0, math.MaxInt64)
		fh.f.reader = filer.NewChunkReaderAtFromClient(fh.f.wfs, chunkViews, fh.f.wfs.chunkCache, fileSize)
	}

	totalRead, err := fh.f.reader.ReadAt(buff, offset)

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		glog.Errorf("file handle read %s: %v", fh.f.fullpath(), err)
	}

	glog.V(4).Infof("file handle read %s [%d,%d] %d : %v", fh.f.fullpath(), offset, offset+int64(totalRead), totalRead, err)

	return int64(totalRead), err
}

// Write to the file handle
func (fh *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {

	fh.Lock()
	defer fh.Unlock()

	// write the request to volume servers
	data := make([]byte, len(req.Data))
	copy(data, req.Data)

	fh.f.entry.Attributes.FileSize = uint64(max(req.Offset+int64(len(data)), int64(fh.f.entry.Attributes.FileSize)))
	glog.V(4).Infof("%v write [%d,%d) %d", fh.f.fullpath(), req.Offset, req.Offset+int64(len(req.Data)), len(req.Data))

	chunks, err := fh.dirtyPages.AddPage(req.Offset, data)
	if err != nil {
		glog.Errorf("%v write fh %d: [%d,%d): %v", fh.f.fullpath(), fh.handle, req.Offset, req.Offset+int64(len(data)), err)
		return fuse.EIO
	}

	resp.Size = len(data)

	if req.Offset == 0 {
		// detect mime type
		fh.contentType = http.DetectContentType(data)
		fh.f.dirtyMetadata = true
	}

	if len(chunks) > 0 {

		fh.f.addChunks(chunks)

		fh.f.dirtyMetadata = true
	}

	return nil
}

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {

	glog.V(4).Infof("Release %v fh %d", fh.f.fullpath(), fh.handle)

	fh.Lock()
	defer fh.Unlock()

	fh.f.isOpen--

	if fh.f.isOpen < 0 {
		glog.V(0).Infof("Release reset %s open count %d => %d", fh.f.Name, fh.f.isOpen, 0)
		fh.f.isOpen = 0
		return nil
	}

	if fh.f.isOpen == 0 {
		fh.doFlush(ctx, req.Header)
		fh.f.wfs.ReleaseHandle(fh.f.fullpath(), fuse.HandleID(fh.handle))
	}

	return nil
}

func (fh *FileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) error {

	fh.Lock()
	defer fh.Unlock()

	return fh.doFlush(ctx, req.Header)
}

func (fh *FileHandle) doFlush(ctx context.Context, header fuse.Header) error {
	// fflush works at fh level
	// send the data to the OS
	glog.V(4).Infof("doFlush %s fh %d", fh.f.fullpath(), fh.handle)

	chunks, err := fh.dirtyPages.saveExistingPagesToStorage()
	if err != nil {
		glog.Errorf("flush %s: %v", fh.f.fullpath(), err)
		return fuse.EIO
	}

	if len(chunks) > 0 {

		fh.f.addChunks(chunks)
		fh.f.dirtyMetadata = true
	}

	if !fh.f.dirtyMetadata {
		return nil
	}

	err = fh.f.wfs.WithFilerClient(func(client filer_pb.SeaweedFilerClient) error {

		if fh.f.entry.Attributes != nil {
			fh.f.entry.Attributes.Mime = fh.contentType
			if fh.f.entry.Attributes.Uid == 0 {
				fh.f.entry.Attributes.Uid = header.Uid
			}
			if fh.f.entry.Attributes.Gid == 0 {
				fh.f.entry.Attributes.Gid = header.Gid
			}
			if fh.f.entry.Attributes.Crtime == 0 {
				fh.f.entry.Attributes.Crtime = time.Now().Unix()
			}
			fh.f.entry.Attributes.Mtime = time.Now().Unix()
			fh.f.entry.Attributes.FileMode = uint32(os.FileMode(fh.f.entry.Attributes.FileMode) &^ fh.f.wfs.option.Umask)
			fh.f.entry.Attributes.Collection = fh.dirtyPages.collection
			fh.f.entry.Attributes.Replication = fh.dirtyPages.replication
		}

		request := &filer_pb.CreateEntryRequest{
			Directory:  fh.f.dir.FullPath(),
			Entry:      fh.f.entry,
			Signatures: []int32{fh.f.wfs.signature},
		}

		glog.V(4).Infof("%s set chunks: %v", fh.f.fullpath(), len(fh.f.entry.Chunks))
		for i, chunk := range fh.f.entry.Chunks {
			glog.V(4).Infof("%s chunks %d: %v [%d,%d)", fh.f.fullpath(), i, chunk.GetFileIdString(), chunk.Offset, chunk.Offset+int64(chunk.Size))
		}

		manifestChunks, nonManifestChunks := filer.SeparateManifestChunks(fh.f.entry.Chunks)

		chunks, _ := filer.CompactFileChunks(filer.LookupFn(fh.f.wfs), nonManifestChunks)
		chunks, manifestErr := filer.MaybeManifestize(fh.f.wfs.saveDataAsChunk(fh.f.dir.FullPath()), chunks)
		if manifestErr != nil {
			// not good, but should be ok
			glog.V(0).Infof("MaybeManifestize: %v", manifestErr)
		}
		fh.f.entry.Chunks = append(chunks, manifestChunks...)
		fh.f.entryViewCache = nil

		fh.f.wfs.mapPbIdFromLocalToFiler(request.Entry)
		defer fh.f.wfs.mapPbIdFromFilerToLocal(request.Entry)

		if err := filer_pb.CreateEntry(client, request); err != nil {
			glog.Errorf("fh flush create %s: %v", fh.f.fullpath(), err)
			return fmt.Errorf("fh flush create %s: %v", fh.f.fullpath(), err)
		}

		fh.f.wfs.metaCache.InsertEntry(context.Background(), filer.FromPbEntry(request.Directory, request.Entry))

		return nil
	})

	if err == nil {
		fh.f.dirtyMetadata = false
	}

	if err != nil {
		glog.Errorf("%v fh %d flush: %v", fh.f.fullpath(), fh.handle, err)
		return fuse.EIO
	}

	return nil
}
