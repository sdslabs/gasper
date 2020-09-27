package weed_server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/util"
)

func (fs *FilerServer) LookupDirectoryEntry(ctx context.Context, req *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {

	glog.V(4).Infof("LookupDirectoryEntry %s", filepath.Join(req.Directory, req.Name))

	entry, err := fs.filer.FindEntry(ctx, util.JoinPath(req.Directory, req.Name))
	if err == filer_pb.ErrNotFound {
		return &filer_pb.LookupDirectoryEntryResponse{}, err
	}
	if err != nil {
		glog.V(3).Infof("LookupDirectoryEntry %s: %+v, ", filepath.Join(req.Directory, req.Name), err)
		return nil, err
	}

	return &filer_pb.LookupDirectoryEntryResponse{
		Entry: &filer_pb.Entry{
			Name:        req.Name,
			IsDirectory: entry.IsDirectory(),
			Attributes:  filer.EntryAttributeToPb(entry),
			Chunks:      entry.Chunks,
			Extended:    entry.Extended,
			HardLinkId:  entry.HardLinkId,
			HardLinkCounter: entry.HardLinkCounter,
		},
	}, nil
}

func (fs *FilerServer) ListEntries(req *filer_pb.ListEntriesRequest, stream filer_pb.SeaweedFiler_ListEntriesServer) error {

	glog.V(4).Infof("ListEntries %v", req)

	limit := int(req.Limit)
	if limit == 0 {
		limit = fs.option.DirListingLimit
	}

	paginationLimit := filer.PaginationSize
	if limit < paginationLimit {
		paginationLimit = limit
	}

	lastFileName := req.StartFromFileName
	includeLastFile := req.InclusiveStartFrom
	for limit > 0 {
		entries, err := fs.filer.ListDirectoryEntries(stream.Context(), util.FullPath(req.Directory), lastFileName, includeLastFile, paginationLimit, req.Prefix)

		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return nil
		}

		includeLastFile = false

		for _, entry := range entries {

			lastFileName = entry.Name()

			if err := stream.Send(&filer_pb.ListEntriesResponse{
				Entry: &filer_pb.Entry{
					Name:        entry.Name(),
					IsDirectory: entry.IsDirectory(),
					Chunks:      entry.Chunks,
					Attributes:  filer.EntryAttributeToPb(entry),
					Extended:    entry.Extended,
					HardLinkId:  entry.HardLinkId,
					HardLinkCounter: entry.HardLinkCounter,
				},
			}); err != nil {
				return err
			}

			limit--
			if limit == 0 {
				return nil
			}
		}

		if len(entries) < paginationLimit {
			break
		}

	}

	return nil
}

func (fs *FilerServer) LookupVolume(ctx context.Context, req *filer_pb.LookupVolumeRequest) (*filer_pb.LookupVolumeResponse, error) {

	resp := &filer_pb.LookupVolumeResponse{
		LocationsMap: make(map[string]*filer_pb.Locations),
	}

	for _, vidString := range req.VolumeIds {
		vid, err := strconv.Atoi(vidString)
		if err != nil {
			glog.V(1).Infof("Unknown volume id %d", vid)
			return nil, err
		}
		var locs []*filer_pb.Location
		locations, found := fs.filer.MasterClient.GetLocations(uint32(vid))
		if !found {
			continue
		}
		for _, loc := range locations {
			locs = append(locs, &filer_pb.Location{
				Url:       loc.Url,
				PublicUrl: loc.PublicUrl,
			})
		}
		resp.LocationsMap[vidString] = &filer_pb.Locations{
			Locations: locs,
		}
	}

	return resp, nil
}

func (fs *FilerServer) lookupFileId(fileId string) (targetUrl string, err error) {
	fid, err := needle.ParseFileIdFromString(fileId)
	if err != nil {
		return "", err
	}
	locations, found := fs.filer.MasterClient.GetLocations(uint32(fid.VolumeId))
	if !found || len(locations) == 0 {
		return "", fmt.Errorf("not found volume %d in %s", fid.VolumeId, fileId)
	}
	return fmt.Sprintf("http://%s/%s", locations[0].Url, fileId), nil
}

func (fs *FilerServer) CreateEntry(ctx context.Context, req *filer_pb.CreateEntryRequest) (resp *filer_pb.CreateEntryResponse, err error) {

	glog.V(4).Infof("CreateEntry %v", req)

	resp = &filer_pb.CreateEntryResponse{}

	chunks, garbage, err2 := fs.cleanupChunks(nil, req.Entry)
	if err2 != nil {
		return &filer_pb.CreateEntryResponse{}, fmt.Errorf("CreateEntry cleanupChunks %s %s: %v", req.Directory, req.Entry.Name, err2)
	}

	createErr := fs.filer.CreateEntry(ctx, &filer.Entry{
		FullPath: util.JoinPath(req.Directory, req.Entry.Name),
		Attr:     filer.PbToEntryAttribute(req.Entry.Attributes),
		Chunks:   chunks,
		Extended: req.Entry.Extended,
		HardLinkId: filer.HardLinkId(req.Entry.HardLinkId),
		HardLinkCounter: req.Entry.HardLinkCounter,
	}, req.OExcl, req.IsFromOtherCluster, req.Signatures)

	if createErr == nil {
		fs.filer.DeleteChunks(garbage)
	} else {
		glog.V(3).Infof("CreateEntry %s: %v", filepath.Join(req.Directory, req.Entry.Name), createErr)
		resp.Error = createErr.Error()
	}

	return
}

func (fs *FilerServer) UpdateEntry(ctx context.Context, req *filer_pb.UpdateEntryRequest) (*filer_pb.UpdateEntryResponse, error) {

	glog.V(4).Infof("UpdateEntry %v", req)

	fullpath := util.Join(req.Directory, req.Entry.Name)
	entry, err := fs.filer.FindEntry(ctx, util.FullPath(fullpath))
	if err != nil {
		return &filer_pb.UpdateEntryResponse{}, fmt.Errorf("not found %s: %v", fullpath, err)
	}

	chunks, garbage, err2 := fs.cleanupChunks(entry, req.Entry)
	if err2 != nil {
		return &filer_pb.UpdateEntryResponse{}, fmt.Errorf("UpdateEntry cleanupChunks %s: %v", fullpath, err2)
	}

	newEntry := &filer.Entry{
		FullPath: util.JoinPath(req.Directory, req.Entry.Name),
		Attr:     entry.Attr,
		Extended: req.Entry.Extended,
		Chunks:   chunks,
		HardLinkId: filer.HardLinkId(req.Entry.HardLinkId),
		HardLinkCounter: req.Entry.HardLinkCounter,
	}

	glog.V(3).Infof("updating %s: %+v, chunks %d: %v => %+v, chunks %d: %v, extended: %v => %v",
		fullpath, entry.Attr, len(entry.Chunks), entry.Chunks,
		req.Entry.Attributes, len(req.Entry.Chunks), req.Entry.Chunks,
		entry.Extended, req.Entry.Extended)

	if req.Entry.Attributes != nil {
		if req.Entry.Attributes.Mtime != 0 {
			newEntry.Attr.Mtime = time.Unix(req.Entry.Attributes.Mtime, 0)
		}
		if req.Entry.Attributes.FileMode != 0 {
			newEntry.Attr.Mode = os.FileMode(req.Entry.Attributes.FileMode)
		}
		newEntry.Attr.Uid = req.Entry.Attributes.Uid
		newEntry.Attr.Gid = req.Entry.Attributes.Gid
		newEntry.Attr.Mime = req.Entry.Attributes.Mime
		newEntry.Attr.UserName = req.Entry.Attributes.UserName
		newEntry.Attr.GroupNames = req.Entry.Attributes.GroupName

	}

	if filer.EqualEntry(entry, newEntry) {
		return &filer_pb.UpdateEntryResponse{}, err
	}

	if err = fs.filer.UpdateEntry(ctx, entry, newEntry); err == nil {
		fs.filer.DeleteChunks(garbage)

		fs.filer.NotifyUpdateEvent(ctx, entry, newEntry, true, req.IsFromOtherCluster, req.Signatures)

	} else {
		glog.V(3).Infof("UpdateEntry %s: %v", filepath.Join(req.Directory, req.Entry.Name), err)
	}

	return &filer_pb.UpdateEntryResponse{}, err
}

func (fs *FilerServer) cleanupChunks(existingEntry *filer.Entry, newEntry *filer_pb.Entry) (chunks, garbage []*filer_pb.FileChunk, err error) {

	// remove old chunks if not included in the new ones
	if existingEntry != nil {
		garbage, err = filer.MinusChunks(fs.lookupFileId, existingEntry.Chunks, newEntry.Chunks)
		if err != nil {
			return newEntry.Chunks, nil, fmt.Errorf("MinusChunks: %v", err)
		}
	}

	// files with manifest chunks are usually large and append only, skip calculating covered chunks
	manifestChunks, nonManifestChunks := filer.SeparateManifestChunks(newEntry.Chunks)

	chunks, coveredChunks := filer.CompactFileChunks(fs.lookupFileId, nonManifestChunks)
	garbage = append(garbage, coveredChunks...)

	if newEntry.Attributes != nil {
		chunks, err = filer.MaybeManifestize(fs.saveAsChunk(
			newEntry.Attributes.Replication,
			newEntry.Attributes.Collection,
			"",
			needle.SecondsToTTL(newEntry.Attributes.TtlSec),
			false), chunks)
		if err != nil {
			// not good, but should be ok
			glog.V(0).Infof("MaybeManifestize: %v", err)
		}
	}

	chunks = append(chunks, manifestChunks...)

	return
}

func (fs *FilerServer) AppendToEntry(ctx context.Context, req *filer_pb.AppendToEntryRequest) (*filer_pb.AppendToEntryResponse, error) {

	glog.V(4).Infof("AppendToEntry %v", req)

	fullpath := util.NewFullPath(req.Directory, req.EntryName)
	var offset int64 = 0
	entry, err := fs.filer.FindEntry(ctx, util.FullPath(fullpath))
	if err == filer_pb.ErrNotFound {
		entry = &filer.Entry{
			FullPath: fullpath,
			Attr: filer.Attr{
				Crtime: time.Now(),
				Mtime:  time.Now(),
				Mode:   os.FileMode(0644),
				Uid:    OS_UID,
				Gid:    OS_GID,
			},
		}
	} else {
		offset = int64(filer.TotalSize(entry.Chunks))
	}

	for _, chunk := range req.Chunks {
		chunk.Offset = offset
		offset += int64(chunk.Size)
	}

	entry.Chunks = append(entry.Chunks, req.Chunks...)

	entry.Chunks, err = filer.MaybeManifestize(fs.saveAsChunk(
		entry.Replication,
		entry.Collection,
		"",
		needle.SecondsToTTL(entry.TtlSec),
		false), entry.Chunks)
	if err != nil {
		// not good, but should be ok
		glog.V(0).Infof("MaybeManifestize: %v", err)
	}

	err = fs.filer.CreateEntry(context.Background(), entry, false, false, nil)

	return &filer_pb.AppendToEntryResponse{}, err
}

func (fs *FilerServer) DeleteEntry(ctx context.Context, req *filer_pb.DeleteEntryRequest) (resp *filer_pb.DeleteEntryResponse, err error) {

	glog.V(4).Infof("DeleteEntry %v", req)

	err = fs.filer.DeleteEntryMetaAndData(ctx, util.JoinPath(req.Directory, req.Name), req.IsRecursive, req.IgnoreRecursiveError, req.IsDeleteData, req.IsFromOtherCluster, req.Signatures)
	resp = &filer_pb.DeleteEntryResponse{}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

func (fs *FilerServer) AssignVolume(ctx context.Context, req *filer_pb.AssignVolumeRequest) (resp *filer_pb.AssignVolumeResponse, err error) {

	ttlStr := ""
	if req.TtlSec > 0 {
		ttlStr = strconv.Itoa(int(req.TtlSec))
	}
	collection, replication, _ := fs.detectCollection(req.ParentPath, req.Collection, req.Replication)

	var altRequest *operation.VolumeAssignRequest

	dataCenter := req.DataCenter
	if dataCenter == "" {
		dataCenter = fs.option.DataCenter
	}

	assignRequest := &operation.VolumeAssignRequest{
		Count:       uint64(req.Count),
		Replication: replication,
		Collection:  collection,
		Ttl:         ttlStr,
		DataCenter:  dataCenter,
	}
	if dataCenter != "" {
		altRequest = &operation.VolumeAssignRequest{
			Count:       uint64(req.Count),
			Replication: replication,
			Collection:  collection,
			Ttl:         ttlStr,
			DataCenter:  "",
		}
	}
	assignResult, err := operation.Assign(fs.filer.GetMaster(), fs.grpcDialOption, assignRequest, altRequest)
	if err != nil {
		glog.V(3).Infof("AssignVolume: %v", err)
		return &filer_pb.AssignVolumeResponse{Error: fmt.Sprintf("assign volume: %v", err)}, nil
	}
	if assignResult.Error != "" {
		glog.V(3).Infof("AssignVolume error: %v", assignResult.Error)
		return &filer_pb.AssignVolumeResponse{Error: fmt.Sprintf("assign volume result: %v", assignResult.Error)}, nil
	}

	return &filer_pb.AssignVolumeResponse{
		FileId:      assignResult.Fid,
		Count:       int32(assignResult.Count),
		Url:         assignResult.Url,
		PublicUrl:   assignResult.PublicUrl,
		Auth:        string(assignResult.Auth),
		Collection:  collection,
		Replication: replication,
	}, nil
}

func (fs *FilerServer) DeleteCollection(ctx context.Context, req *filer_pb.DeleteCollectionRequest) (resp *filer_pb.DeleteCollectionResponse, err error) {

	glog.V(4).Infof("DeleteCollection %v", req)

	err = fs.filer.MasterClient.WithClient(func(client master_pb.SeaweedClient) error {
		_, err := client.CollectionDelete(context.Background(), &master_pb.CollectionDeleteRequest{
			Name: req.GetCollection(),
		})
		return err
	})

	return &filer_pb.DeleteCollectionResponse{}, err
}

func (fs *FilerServer) Statistics(ctx context.Context, req *filer_pb.StatisticsRequest) (resp *filer_pb.StatisticsResponse, err error) {

	var output *master_pb.StatisticsResponse

	err = fs.filer.MasterClient.WithClient(func(masterClient master_pb.SeaweedClient) error {
		grpcResponse, grpcErr := masterClient.Statistics(context.Background(), &master_pb.StatisticsRequest{
			Replication: req.Replication,
			Collection:  req.Collection,
			Ttl:         req.Ttl,
		})
		if grpcErr != nil {
			return grpcErr
		}

		output = grpcResponse
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &filer_pb.StatisticsResponse{
		TotalSize: output.TotalSize,
		UsedSize:  output.UsedSize,
		FileCount: output.FileCount,
	}, nil
}

func (fs *FilerServer) GetFilerConfiguration(ctx context.Context, req *filer_pb.GetFilerConfigurationRequest) (resp *filer_pb.GetFilerConfigurationResponse, err error) {

	t := &filer_pb.GetFilerConfigurationResponse{
		Masters:            fs.option.Masters,
		Collection:         fs.option.Collection,
		Replication:        fs.option.DefaultReplication,
		MaxMb:              uint32(fs.option.MaxMB),
		DirBuckets:         fs.filer.DirBucketsPath,
		Cipher:             fs.filer.Cipher,
		Signature:          fs.filer.Signature,
		MetricsAddress:     fs.metricsAddress,
		MetricsIntervalSec: int32(fs.metricsIntervalSec),
	}

	glog.V(4).Infof("GetFilerConfiguration: %v", t)

	return t, nil
}

func (fs *FilerServer) KeepConnected(stream filer_pb.SeaweedFiler_KeepConnectedServer) error {

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	clientName := fmt.Sprintf("%s:%d", req.Name, req.GrpcPort)
	m := make(map[string]bool)
	for _, tp := range req.Resources {
		m[tp] = true
	}
	fs.brokersLock.Lock()
	fs.brokers[clientName] = m
	glog.V(0).Infof("+ broker %v", clientName)
	fs.brokersLock.Unlock()

	defer func() {
		fs.brokersLock.Lock()
		delete(fs.brokers, clientName)
		glog.V(0).Infof("- broker %v: %v", clientName, err)
		fs.brokersLock.Unlock()
	}()

	for {
		if err := stream.Send(&filer_pb.KeepConnectedResponse{}); err != nil {
			glog.V(0).Infof("send broker %v: %+v", clientName, err)
			return err
		}
		// println("replied")

		if _, err := stream.Recv(); err != nil {
			glog.V(0).Infof("recv broker %v: %v", clientName, err)
			return err
		}
		// println("received")
	}

}

func (fs *FilerServer) LocateBroker(ctx context.Context, req *filer_pb.LocateBrokerRequest) (resp *filer_pb.LocateBrokerResponse, err error) {

	resp = &filer_pb.LocateBrokerResponse{}

	fs.brokersLock.Lock()
	defer fs.brokersLock.Unlock()

	var localBrokers []*filer_pb.LocateBrokerResponse_Resource

	for b, m := range fs.brokers {
		if _, found := m[req.Resource]; found {
			resp.Found = true
			resp.Resources = []*filer_pb.LocateBrokerResponse_Resource{
				{
					GrpcAddresses: b,
					ResourceCount: int32(len(m)),
				},
			}
			return
		}
		localBrokers = append(localBrokers, &filer_pb.LocateBrokerResponse_Resource{
			GrpcAddresses: b,
			ResourceCount: int32(len(m)),
		})
	}

	resp.Resources = localBrokers

	return resp, nil

}
