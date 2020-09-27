package redis

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/util"
)

const (
	DIR_LIST_MARKER = "\x00"
)

type UniversalRedisStore struct {
	Client redis.UniversalClient
}

func (store *UniversalRedisStore) BeginTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (store *UniversalRedisStore) CommitTransaction(ctx context.Context) error {
	return nil
}
func (store *UniversalRedisStore) RollbackTransaction(ctx context.Context) error {
	return nil
}

func (store *UniversalRedisStore) InsertEntry(ctx context.Context, entry *filer.Entry) (err error) {

	value, err := entry.EncodeAttributesAndChunks()
	if err != nil {
		return fmt.Errorf("encoding %s %+v: %v", entry.FullPath, entry.Attr, err)
	}

	if len(entry.Chunks) > 50 {
		value = util.MaybeGzipData(value)
	}

	_, err = store.Client.Set(string(entry.FullPath), value, time.Duration(entry.TtlSec)*time.Second).Result()

	if err != nil {
		return fmt.Errorf("persisting %s : %v", entry.FullPath, err)
	}

	dir, name := entry.FullPath.DirAndName()
	if name != "" {
		_, err = store.Client.SAdd(genDirectoryListKey(dir), name).Result()
		if err != nil {
			return fmt.Errorf("persisting %s in parent dir: %v", entry.FullPath, err)
		}
	}

	return nil
}

func (store *UniversalRedisStore) UpdateEntry(ctx context.Context, entry *filer.Entry) (err error) {

	return store.InsertEntry(ctx, entry)
}

func (store *UniversalRedisStore) FindEntry(ctx context.Context, fullpath util.FullPath) (entry *filer.Entry, err error) {

	data, err := store.Client.Get(string(fullpath)).Result()
	if err == redis.Nil {
		return nil, filer_pb.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get %s : %v", fullpath, err)
	}

	entry = &filer.Entry{
		FullPath: fullpath,
	}
	err = entry.DecodeAttributesAndChunks(util.MaybeDecompressData([]byte(data)))
	if err != nil {
		return entry, fmt.Errorf("decode %s : %v", entry.FullPath, err)
	}

	return entry, nil
}

func (store *UniversalRedisStore) DeleteEntry(ctx context.Context, fullpath util.FullPath) (err error) {

	_, err = store.Client.Del(string(fullpath)).Result()

	if err != nil {
		return fmt.Errorf("delete %s : %v", fullpath, err)
	}

	dir, name := fullpath.DirAndName()
	if name != "" {
		_, err = store.Client.SRem(genDirectoryListKey(dir), name).Result()
		if err != nil {
			return fmt.Errorf("delete %s in parent dir: %v", fullpath, err)
		}
	}

	return nil
}

func (store *UniversalRedisStore) DeleteFolderChildren(ctx context.Context, fullpath util.FullPath) (err error) {

	members, err := store.Client.SMembers(genDirectoryListKey(string(fullpath))).Result()
	if err != nil {
		return fmt.Errorf("delete folder %s : %v", fullpath, err)
	}

	for _, fileName := range members {
		path := util.NewFullPath(string(fullpath), fileName)
		_, err = store.Client.Del(string(path)).Result()
		if err != nil {
			return fmt.Errorf("delete %s in parent dir: %v", fullpath, err)
		}
	}

	return nil
}

func (store *UniversalRedisStore) ListDirectoryPrefixedEntries(ctx context.Context, fullpath util.FullPath, startFileName string, inclusive bool, limit int, prefix string) (entries []*filer.Entry, err error) {
	return nil, filer.ErrUnsupportedListDirectoryPrefixed
}

func (store *UniversalRedisStore) ListDirectoryEntries(ctx context.Context, fullpath util.FullPath, startFileName string, inclusive bool,
	limit int) (entries []*filer.Entry, err error) {

	dirListKey := genDirectoryListKey(string(fullpath))
	members, err := store.Client.SMembers(dirListKey).Result()
	if err != nil {
		return nil, fmt.Errorf("list %s : %v", fullpath, err)
	}

	// skip
	if startFileName != "" {
		var t []string
		for _, m := range members {
			if strings.Compare(m, startFileName) >= 0 {
				if m == startFileName {
					if inclusive {
						t = append(t, m)
					}
				} else {
					t = append(t, m)
				}
			}
		}
		members = t
	}

	// sort
	sort.Slice(members, func(i, j int) bool {
		return strings.Compare(members[i], members[j]) < 0
	})

	// limit
	if limit < len(members) {
		members = members[:limit]
	}

	// fetch entry meta
	for _, fileName := range members {
		path := util.NewFullPath(string(fullpath), fileName)
		entry, err := store.FindEntry(ctx, path)
		if err != nil {
			glog.V(0).Infof("list %s : %v", path, err)
		} else {
			if entry.TtlSec > 0 {
				if entry.Attr.Crtime.Add(time.Duration(entry.TtlSec) * time.Second).Before(time.Now()) {
					store.Client.Del(string(path)).Result()
					store.Client.SRem(dirListKey, fileName).Result()
					continue
				}
			}
			entries = append(entries, entry)
		}
	}

	return entries, err
}

func genDirectoryListKey(dir string) (dirList string) {
	return dir + DIR_LIST_MARKER
}

func (store *UniversalRedisStore) Shutdown() {
	store.Client.Close()
}
