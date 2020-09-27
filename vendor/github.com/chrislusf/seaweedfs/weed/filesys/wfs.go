package filesys

import (
	"context"
	"fmt"
	"math"
	"os"
	"path"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/chrislusf/seaweedfs/weed/util/grace"

	"github.com/seaweedfs/fuse"
	"github.com/seaweedfs/fuse/fs"

	"github.com/chrislusf/seaweedfs/weed/filesys/meta_cache"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/util"
	"github.com/chrislusf/seaweedfs/weed/util/chunk_cache"
)

type Option struct {
	FilerGrpcAddress   string
	GrpcDialOption     grpc.DialOption
	FilerMountRootPath string
	Collection         string
	Replication        string
	TtlSec             int32
	ChunkSizeLimit     int64
	CacheDir           string
	CacheSizeMB        int64
	DataCenter         string
	EntryCacheTtl      time.Duration
	Umask              os.FileMode

	MountCtime time.Time
	MountMtime time.Time

	OutsideContainerClusterMode bool // whether the mount runs outside SeaweedFS containers
	Cipher                      bool // whether encrypt data on volume server
	UidGidMapper                *meta_cache.UidGidMapper
}

var _ = fs.FS(&WFS{})
var _ = fs.FSStatfser(&WFS{})

type WFS struct {
	option *Option

	// contains all open handles, protected by handlesLock
	handlesLock sync.Mutex
	handles     map[uint64]*FileHandle

	bufPool sync.Pool

	stats statsCache

	root        fs.Node
	fsNodeCache *FsCache

	chunkCache *chunk_cache.TieredChunkCache
	metaCache  *meta_cache.MetaCache
	signature  int32
}
type statsCache struct {
	filer_pb.StatisticsResponse
	lastChecked int64 // unix time in seconds
}

func NewSeaweedFileSystem(option *Option) *WFS {
	wfs := &WFS{
		option:  option,
		handles: make(map[uint64]*FileHandle),
		bufPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, option.ChunkSizeLimit)
			},
		},
		signature: util.RandomInt32(),
	}
	cacheUniqueId := util.Md5String([]byte(option.FilerGrpcAddress + option.FilerMountRootPath + util.Version()))[0:4]
	cacheDir := path.Join(option.CacheDir, cacheUniqueId)
	if option.CacheSizeMB > 0 {
		os.MkdirAll(cacheDir, os.FileMode(0777)&^option.Umask)
		wfs.chunkCache = chunk_cache.NewTieredChunkCache(256, cacheDir, option.CacheSizeMB)
	}

	wfs.metaCache = meta_cache.NewMetaCache(path.Join(cacheDir, "meta"), option.UidGidMapper)
	startTime := time.Now()
	go meta_cache.SubscribeMetaEvents(wfs.metaCache, wfs.signature, wfs, wfs.option.FilerMountRootPath, startTime.UnixNano())
	grace.OnInterrupt(func() {
		wfs.metaCache.Shutdown()
	})

	entry, _ := filer_pb.GetEntry(wfs, util.FullPath(wfs.option.FilerMountRootPath))
	wfs.root = &Dir{name: wfs.option.FilerMountRootPath, wfs: wfs, entry: entry}
	wfs.fsNodeCache = newFsCache(wfs.root)

	return wfs
}

func (wfs *WFS) Root() (fs.Node, error) {
	return wfs.root, nil
}

func (wfs *WFS) AcquireHandle(file *File, uid, gid uint32) (fileHandle *FileHandle) {

	fullpath := file.fullpath()
	glog.V(4).Infof("AcquireHandle %s uid=%d gid=%d", fullpath, uid, gid)

	wfs.handlesLock.Lock()
	defer wfs.handlesLock.Unlock()

	inodeId := file.fullpath().AsInode()
	existingHandle, found := wfs.handles[inodeId]
	if found && existingHandle != nil {
		file.isOpen++
		return existingHandle
	}

	fileHandle = newFileHandle(file, uid, gid)
	file.maybeLoadEntry(context.Background())
	file.isOpen++

	wfs.handles[inodeId] = fileHandle
	fileHandle.handle = inodeId

	return
}

func (wfs *WFS) ReleaseHandle(fullpath util.FullPath, handleId fuse.HandleID) {
	wfs.handlesLock.Lock()
	defer wfs.handlesLock.Unlock()

	glog.V(4).Infof("%s ReleaseHandle id %d current handles length %d", fullpath, handleId, len(wfs.handles))

	delete(wfs.handles, fullpath.AsInode())

	return
}

// Statfs is called to obtain file system metadata. Implements fuse.FSStatfser
func (wfs *WFS) Statfs(ctx context.Context, req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error {

	glog.V(4).Infof("reading fs stats: %+v", req)

	if wfs.stats.lastChecked < time.Now().Unix()-20 {

		err := wfs.WithFilerClient(func(client filer_pb.SeaweedFilerClient) error {

			request := &filer_pb.StatisticsRequest{
				Collection:  wfs.option.Collection,
				Replication: wfs.option.Replication,
				Ttl:         fmt.Sprintf("%ds", wfs.option.TtlSec),
			}

			glog.V(4).Infof("reading filer stats: %+v", request)
			resp, err := client.Statistics(context.Background(), request)
			if err != nil {
				glog.V(0).Infof("reading filer stats %v: %v", request, err)
				return err
			}
			glog.V(4).Infof("read filer stats: %+v", resp)

			wfs.stats.TotalSize = resp.TotalSize
			wfs.stats.UsedSize = resp.UsedSize
			wfs.stats.FileCount = resp.FileCount
			wfs.stats.lastChecked = time.Now().Unix()

			return nil
		})
		if err != nil {
			glog.V(0).Infof("filer Statistics: %v", err)
			return err
		}
	}

	totalDiskSize := wfs.stats.TotalSize
	usedDiskSize := wfs.stats.UsedSize
	actualFileCount := wfs.stats.FileCount

	// Compute the total number of available blocks
	resp.Blocks = totalDiskSize / blockSize

	// Compute the number of used blocks
	numBlocks := uint64(usedDiskSize / blockSize)

	// Report the number of free and available blocks for the block size
	resp.Bfree = resp.Blocks - numBlocks
	resp.Bavail = resp.Blocks - numBlocks
	resp.Bsize = uint32(blockSize)

	// Report the total number of possible files in the file system (and those free)
	resp.Files = math.MaxInt64
	resp.Ffree = math.MaxInt64 - actualFileCount

	// Report the maximum length of a name and the minimum fragment size
	resp.Namelen = 1024
	resp.Frsize = uint32(blockSize)

	return nil
}

func (wfs *WFS) mapPbIdFromFilerToLocal(entry *filer_pb.Entry) {
	if entry.Attributes == nil {
		return
	}
	entry.Attributes.Uid, entry.Attributes.Gid = wfs.option.UidGidMapper.FilerToLocal(entry.Attributes.Uid, entry.Attributes.Gid)
}
func (wfs *WFS) mapPbIdFromLocalToFiler(entry *filer_pb.Entry) {
	if entry.Attributes == nil {
		return
	}
	entry.Attributes.Uid, entry.Attributes.Gid = wfs.option.UidGidMapper.LocalToFiler(entry.Attributes.Uid, entry.Attributes.Gid)
}
