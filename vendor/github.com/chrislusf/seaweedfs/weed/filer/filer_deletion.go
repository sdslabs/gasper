package filer

import (
	"strings"
	"time"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/wdclient"
)

func LookupByMasterClientFn(masterClient *wdclient.MasterClient) func(vids []string) (map[string]operation.LookupResult, error) {
	return func(vids []string) (map[string]operation.LookupResult, error) {
		m := make(map[string]operation.LookupResult)
		for _, vid := range vids {
			locs, _ := masterClient.GetVidLocations(vid)
			var locations []operation.Location
			for _, loc := range locs {
				locations = append(locations, operation.Location{
					Url:       loc.Url,
					PublicUrl: loc.PublicUrl,
				})
			}
			m[vid] = operation.LookupResult{
				VolumeId:  vid,
				Locations: locations,
			}
		}
		return m, nil
	}
}

func (f *Filer) loopProcessingDeletion() {

	lookupFunc := LookupByMasterClientFn(f.MasterClient)

	DeletionBatchSize := 100000 // roughly 20 bytes cost per file id.

	var deletionCount int
	for {
		deletionCount = 0
		f.fileIdDeletionQueue.Consume(func(fileIds []string) {
			for len(fileIds) > 0 {
				var toDeleteFileIds []string
				if len(fileIds) > DeletionBatchSize {
					toDeleteFileIds = fileIds[:DeletionBatchSize]
					fileIds = fileIds[DeletionBatchSize:]
				} else {
					toDeleteFileIds = fileIds
					fileIds = fileIds[:0]
				}
				deletionCount = len(toDeleteFileIds)
				_, err := operation.DeleteFilesWithLookupVolumeId(f.GrpcDialOption, toDeleteFileIds, lookupFunc)
				if err != nil {
					if !strings.Contains(err.Error(), "already deleted") {
						glog.V(0).Infof("deleting fileIds len=%d error: %v", deletionCount, err)
					}
				} else {
					glog.V(1).Infof("deleting fileIds len=%d", deletionCount)
				}
			}
		})

		if deletionCount == 0 {
			time.Sleep(1123 * time.Millisecond)
		}
	}
}

func (f *Filer) DeleteChunks(chunks []*filer_pb.FileChunk) {
	for _, chunk := range chunks {
		if !chunk.IsChunkManifest {
			f.fileIdDeletionQueue.EnQueue(chunk.GetFileIdString())
			continue
		}
		dataChunks, manifestResolveErr := ResolveOneChunkManifest(f.MasterClient.LookupFileId, chunk)
		if manifestResolveErr != nil {
			glog.V(0).Infof("failed to resolve manifest %s: %v", chunk.FileId, manifestResolveErr)
		}
		for _, dChunk := range dataChunks {
			f.fileIdDeletionQueue.EnQueue(dChunk.GetFileIdString())
		}
		f.fileIdDeletionQueue.EnQueue(chunk.GetFileIdString())
	}
}

func (f *Filer) deleteChunksIfNotNew(oldEntry, newEntry *Entry) {

	if oldEntry == nil {
		return
	}
	if newEntry == nil {
		f.DeleteChunks(oldEntry.Chunks)
	}

	var toDelete []*filer_pb.FileChunk
	newChunkIds := make(map[string]bool)
	for _, newChunk := range newEntry.Chunks {
		newChunkIds[newChunk.GetFileIdString()] = true
	}

	for _, oldChunk := range oldEntry.Chunks {
		if _, found := newChunkIds[oldChunk.GetFileIdString()]; !found {
			toDelete = append(toDelete, oldChunk)
		}
	}
	f.DeleteChunks(toDelete)
}
