package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/replication"
	"github.com/chrislusf/seaweedfs/weed/replication/sink/filersink"
	"github.com/chrislusf/seaweedfs/weed/replication/source"
	"github.com/chrislusf/seaweedfs/weed/security"
	"github.com/chrislusf/seaweedfs/weed/util"
	"github.com/chrislusf/seaweedfs/weed/util/grace"
	"google.golang.org/grpc"
	"io"
	"strings"
	"time"
)

type SyncOptions struct {
	isActivePassive *bool
	filerA          *string
	filerB          *string
	aPath           *string
	bPath           *string
	aReplication    *string
	bReplication    *string
	aCollection     *string
	bCollection     *string
	aTtlSec         *int
	bTtlSec         *int
	aDebug          *bool
	bDebug          *bool
}

var (
	syncOptions    SyncOptions
	syncCpuProfile *string
	syncMemProfile *string
)

func init() {
	cmdFilerSynchronize.Run = runFilerSynchronize // break init cycle
	syncOptions.isActivePassive = cmdFilerSynchronize.Flag.Bool("isActivePassive", false, "one directional follow if true")
	syncOptions.filerA = cmdFilerSynchronize.Flag.String("a", "", "filer A in one SeaweedFS cluster")
	syncOptions.filerB = cmdFilerSynchronize.Flag.String("b", "", "filer B in the other SeaweedFS cluster")
	syncOptions.aPath = cmdFilerSynchronize.Flag.String("a.path", "/", "directory to sync on filer A")
	syncOptions.bPath = cmdFilerSynchronize.Flag.String("b.path", "/", "directory to sync on filer B")
	syncOptions.aReplication = cmdFilerSynchronize.Flag.String("a.replication", "", "replication on filer A")
	syncOptions.bReplication = cmdFilerSynchronize.Flag.String("b.replication", "", "replication on filer B")
	syncOptions.aCollection = cmdFilerSynchronize.Flag.String("a.collection", "", "collection on filer A")
	syncOptions.bCollection = cmdFilerSynchronize.Flag.String("b.collection", "", "collection on filer B")
	syncOptions.aTtlSec = cmdFilerSynchronize.Flag.Int("a.ttlSec", 0, "ttl in seconds on filer A")
	syncOptions.bTtlSec = cmdFilerSynchronize.Flag.Int("b.ttlSec", 0, "ttl in seconds on filer B")
	syncOptions.aDebug = cmdFilerSynchronize.Flag.Bool("a.debug", false, "debug mode to print out filer A received files")
	syncOptions.bDebug = cmdFilerSynchronize.Flag.Bool("b.debug", false, "debug mode to print out filer B received files")
	syncCpuProfile = cmdFilerSynchronize.Flag.String("cpuprofile", "", "cpu profile output file")
	syncMemProfile = cmdFilerSynchronize.Flag.String("memprofile", "", "memory profile output file")
}

var cmdFilerSynchronize = &Command{
	UsageLine: "filer.sync -a=<oneFilerHost>:<oneFilerPort> -b=<otherFilerHost>:<otherFilerPort>",
	Short:     "continuously synchronize between two active-active or active-passive SeaweedFS clusters",
	Long: `continuously synchronize file changes between two active-active or active-passive filers

	filer.sync listens on filer notifications. If any file is updated, it will fetch the updated content,
	and write to the other destination. Different from filer.replicate:

	* filer.sync only works between two filers.
	* filer.sync does not need any special message queue setup.
	* filer.sync supports both active-active and active-passive modes.
	
	If restarted, the synchronization will resume from the previous checkpoints, persisted every minute.
	A fresh sync will start from the earliest metadata logs.

`,
}

func runFilerSynchronize(cmd *Command, args []string) bool {

	grpcDialOption := security.LoadClientTLS(util.GetViper(), "grpc.client")

	grace.SetupProfiling(*syncCpuProfile, *syncMemProfile)

	go func() {
		for {
			err := doSubscribeFilerMetaChanges(grpcDialOption, *syncOptions.filerA, *syncOptions.aPath, *syncOptions.filerB,
				*syncOptions.bPath, *syncOptions.bReplication, *syncOptions.bCollection, *syncOptions.bTtlSec, *syncOptions.bDebug)
			if err != nil {
				glog.Errorf("sync from %s to %s: %v", *syncOptions.filerA, *syncOptions.filerB, err)
				time.Sleep(1747 * time.Millisecond)
			}
		}
	}()

	if !*syncOptions.isActivePassive {
		go func() {
			for {
				err := doSubscribeFilerMetaChanges(grpcDialOption, *syncOptions.filerB, *syncOptions.bPath, *syncOptions.filerA,
					*syncOptions.aPath, *syncOptions.aReplication, *syncOptions.aCollection, *syncOptions.aTtlSec, *syncOptions.aDebug)
				if err != nil {
					glog.Errorf("sync from %s to %s: %v", *syncOptions.filerB, *syncOptions.filerA, err)
					time.Sleep(2147 * time.Millisecond)
				}
			}
		}()
	}

	select {}

	return true
}

func doSubscribeFilerMetaChanges(grpcDialOption grpc.DialOption, sourceFiler, sourcePath, targetFiler, targetPath string,
	replicationStr, collection string, ttlSec int, debug bool) error {

	// read source filer signature
	sourceFilerSignature, sourceErr := replication.ReadFilerSignature(grpcDialOption, sourceFiler)
	if sourceErr != nil {
		return sourceErr
	}
	// read target filer signature
	targetFilerSignature, targetErr := replication.ReadFilerSignature(grpcDialOption, targetFiler)
	if targetErr != nil {
		return targetErr
	}

	// if first time, start from now
	// if has previously synced, resume from that point of time
	sourceFilerOffsetTsNs, err := readSyncOffset(grpcDialOption, targetFiler, sourceFilerSignature)
	if err != nil {
		return err
	}

	glog.V(0).Infof("start sync %s(%d) => %s(%d) from %v(%d)", sourceFiler, sourceFilerSignature, targetFiler, targetFilerSignature, time.Unix(0, sourceFilerOffsetTsNs), sourceFilerOffsetTsNs)

	// create filer sink
	filerSource := &source.FilerSource{}
	filerSource.DoInitialize(pb.ServerToGrpcAddress(sourceFiler), sourcePath)
	filerSink := &filersink.FilerSink{}
	filerSink.DoInitialize(pb.ServerToGrpcAddress(targetFiler), targetPath, replicationStr, collection, ttlSec, grpcDialOption)
	filerSink.SetSourceFiler(filerSource)

	processEventFn := func(resp *filer_pb.SubscribeMetadataResponse) error {
		message := resp.EventNotification

		var sourceOldKey, sourceNewKey util.FullPath
		if message.OldEntry != nil {
			sourceOldKey = util.FullPath(resp.Directory).Child(message.OldEntry.Name)
		}
		if message.NewEntry != nil {
			sourceNewKey = util.FullPath(message.NewParentPath).Child(message.NewEntry.Name)
		}

		for _, sig := range message.Signatures {
			if sig == targetFilerSignature && targetFilerSignature != 0 {
				fmt.Printf("%s skipping %s change to %v\n", targetFiler, sourceFiler, message)
				return nil
			}
		}
		if debug {
			fmt.Printf("%s check %s change %s,%s sig %v, target sig: %v\n", targetFiler, sourceFiler, sourceOldKey, sourceNewKey, message.Signatures, targetFilerSignature)
		}

		if !strings.HasPrefix(resp.Directory, sourcePath) {
			return nil
		}

		// handle deletions
		if message.OldEntry != nil && message.NewEntry == nil {
			if !strings.HasPrefix(string(sourceOldKey), sourcePath) {
				return nil
			}
			key := util.Join(targetPath, string(sourceOldKey)[len(sourcePath):])
			return filerSink.DeleteEntry(key, message.OldEntry.IsDirectory, message.DeleteChunks, message.Signatures)
		}

		// handle new entries
		if message.OldEntry == nil && message.NewEntry != nil {
			if !strings.HasPrefix(string(sourceNewKey), sourcePath) {
				return nil
			}
			key := util.Join(targetPath, string(sourceNewKey)[len(sourcePath):])
			return filerSink.CreateEntry(key, message.NewEntry, message.Signatures)
		}

		// this is something special?
		if message.OldEntry == nil && message.NewEntry == nil {
			return nil
		}

		// handle updates
		if strings.HasPrefix(string(sourceOldKey), sourcePath) {
			// old key is in the watched directory
			if strings.HasPrefix(string(sourceNewKey), sourcePath) {
				// new key is also in the watched directory
				oldKey := util.Join(targetPath, string(sourceOldKey)[len(sourcePath):])
				message.NewParentPath = util.Join(targetPath, message.NewParentPath[len(sourcePath):])
				foundExisting, err := filerSink.UpdateEntry(string(oldKey), message.OldEntry, message.NewParentPath, message.NewEntry, message.DeleteChunks, message.Signatures)
				if foundExisting {
					return err
				}

				// not able to find old entry
				if err = filerSink.DeleteEntry(string(oldKey), message.OldEntry.IsDirectory, false, message.Signatures); err != nil {
					return fmt.Errorf("delete old entry %v: %v", oldKey, err)
				}

				// create the new entry
				newKey := util.Join(targetPath, string(sourceNewKey)[len(sourcePath):])
				return filerSink.CreateEntry(newKey, message.NewEntry, message.Signatures)

			} else {
				// new key is outside of the watched directory
				key := util.Join(targetPath, string(sourceOldKey)[len(sourcePath):])
				return filerSink.DeleteEntry(key, message.OldEntry.IsDirectory, message.DeleteChunks, message.Signatures)
			}
		} else {
			// old key is outside of the watched directory
			if strings.HasPrefix(string(sourceNewKey), sourcePath) {
				// new key is in the watched directory
				key := util.Join(targetPath, string(sourceNewKey)[len(sourcePath):])
				return filerSink.CreateEntry(key, message.NewEntry, message.Signatures)
			} else {
				// new key is also outside of the watched directory
				// skip
			}
		}

		return nil
	}

	return pb.WithFilerClient(sourceFiler, grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stream, err := client.SubscribeMetadata(ctx, &filer_pb.SubscribeMetadataRequest{
			ClientName: "syncTo_" + targetFiler,
			PathPrefix: sourcePath,
			SinceNs:    sourceFilerOffsetTsNs,
			Signature:  targetFilerSignature,
		})
		if err != nil {
			return fmt.Errorf("listen: %v", err)
		}

		var counter int64
		var lastWriteTime time.Time
		for {
			resp, listenErr := stream.Recv()
			if listenErr == io.EOF {
				return nil
			}
			if listenErr != nil {
				return listenErr
			}

			if err := processEventFn(resp); err != nil {
				return err
			}

			counter++
			if lastWriteTime.Add(3 * time.Second).Before(time.Now()) {
				glog.V(0).Infof("sync %s => %s progressed to %v %0.2f/sec", sourceFiler, targetFiler, time.Unix(0, resp.TsNs), float64(counter)/float64(3))
				counter = 0
				lastWriteTime = time.Now()
				if err := writeSyncOffset(grpcDialOption, targetFiler, sourceFilerSignature, resp.TsNs); err != nil {
					return err
				}
			}

		}

	})

}

const (
	SyncKeyPrefix = "sync."
)

func readSyncOffset(grpcDialOption grpc.DialOption, filer string, filerSignature int32) (lastOffsetTsNs int64, readErr error) {

	readErr = pb.WithFilerClient(filer, grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		syncKey := []byte(SyncKeyPrefix + "____")
		util.Uint32toBytes(syncKey[len(SyncKeyPrefix):len(SyncKeyPrefix)+4], uint32(filerSignature))

		resp, err := client.KvGet(context.Background(), &filer_pb.KvGetRequest{Key: syncKey})
		if err != nil {
			return err
		}

		if len(resp.Error) != 0 {
			return errors.New(resp.Error)
		}
		if len(resp.Value) < 8 {
			return nil
		}

		lastOffsetTsNs = int64(util.BytesToUint64(resp.Value))

		return nil
	})

	return

}

func writeSyncOffset(grpcDialOption grpc.DialOption, filer string, filerSignature int32, offsetTsNs int64) error {
	return pb.WithFilerClient(filer, grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {

		syncKey := []byte(SyncKeyPrefix + "____")
		util.Uint32toBytes(syncKey[len(SyncKeyPrefix):len(SyncKeyPrefix)+4], uint32(filerSignature))

		valueBuf := make([]byte, 8)
		util.Uint64toBytes(valueBuf, uint64(offsetTsNs))

		resp, err := client.KvPut(context.Background(), &filer_pb.KvPutRequest{
			Key:   syncKey,
			Value: valueBuf,
		})
		if err != nil {
			return err
		}

		if len(resp.Error) != 0 {
			return errors.New(resp.Error)
		}

		return nil

	})

}
