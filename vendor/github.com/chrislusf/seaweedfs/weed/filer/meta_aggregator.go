package filer

import (
	"context"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/util"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/util/log_buffer"
)

type MetaAggregator struct {
	filers         []string
	grpcDialOption grpc.DialOption
	MetaLogBuffer  *log_buffer.LogBuffer
	// notifying clients
	ListenersLock sync.Mutex
	ListenersCond *sync.Cond
}

// MetaAggregator only aggregates data "on the fly". The logs are not re-persisted to disk.
// The old data comes from what each LocalMetadata persisted on disk.
func NewMetaAggregator(filers []string, grpcDialOption grpc.DialOption) *MetaAggregator {
	t := &MetaAggregator{
		filers:         filers,
		grpcDialOption: grpcDialOption,
	}
	t.ListenersCond = sync.NewCond(&t.ListenersLock)
	t.MetaLogBuffer = log_buffer.NewLogBuffer(LogFlushInterval, nil, func() {
		t.ListenersCond.Broadcast()
	})
	return t
}

func (ma *MetaAggregator) StartLoopSubscribe(f *Filer, self string) {
	for _, filer := range ma.filers {
		go ma.subscribeToOneFiler(f, self, filer)
	}
}

func (ma *MetaAggregator) subscribeToOneFiler(f *Filer, self string, peer string) {

	/*
		Each filer reads the "filer.store.id", which is the store's signature when filer starts.

		When reading from other filers' local meta changes:
		* if the received change does not contain signature from self, apply the change to current filer store.

		Upon connecting to other filers, need to remember their signature and their offsets.

	*/

	var maybeReplicateMetadataChange func(*filer_pb.SubscribeMetadataResponse)
	lastPersistTime := time.Now()
	lastTsNs := time.Now().Add(-LogFlushInterval).UnixNano()

	peerSignature, err := ma.readFilerStoreSignature(peer)
	for err != nil {
		glog.V(0).Infof("connecting to peer filer %s: %v", peer, err)
		time.Sleep(1357 * time.Millisecond)
		peerSignature, err = ma.readFilerStoreSignature(peer)
	}

	if peerSignature != f.Signature {
		if prevTsNs, err := ma.readOffset(f, peer, peerSignature); err == nil {
			lastTsNs = prevTsNs
		}

		glog.V(0).Infof("follow peer: %v, last %v (%d)", peer, time.Unix(0, lastTsNs), lastTsNs)
		var counter int64
		var synced bool
		maybeReplicateMetadataChange = func(event *filer_pb.SubscribeMetadataResponse) {
			if err := Replay(f.Store, event); err != nil {
				glog.Errorf("failed to reply metadata change from %v: %v", peer, err)
				return
			}
			counter++
			if lastPersistTime.Add(time.Minute).Before(time.Now()) {
				if err := ma.updateOffset(f, peer, peerSignature, event.TsNs); err == nil {
					if event.TsNs < time.Now().Add(-2*time.Minute).UnixNano() {
						glog.V(0).Infof("sync with %s progressed to: %v %0.2f/sec", peer, time.Unix(0, event.TsNs), float64(counter)/60.0)
					} else if !synced {
						synced = true
						glog.V(0).Infof("synced with %s", peer)
					}
					lastPersistTime = time.Now()
					counter = 0
				} else {
					glog.V(0).Infof("failed to update offset for %v: %v", peer, err)
				}
			}
		}
	}

	processEventFn := func(event *filer_pb.SubscribeMetadataResponse) error {
		data, err := proto.Marshal(event)
		if err != nil {
			glog.Errorf("failed to marshal subscribed filer_pb.SubscribeMetadataResponse %+v: %v", event, err)
			return err
		}
		dir := event.Directory
		// println("received meta change", dir, "size", len(data))
		ma.MetaLogBuffer.AddToBuffer([]byte(dir), data, 0)
		if maybeReplicateMetadataChange != nil {
			maybeReplicateMetadataChange(event)
		}
		return nil
	}

	for {
		err := pb.WithFilerClient(peer, ma.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream, err := client.SubscribeLocalMetadata(ctx, &filer_pb.SubscribeMetadataRequest{
				ClientName: "filer:" + self,
				PathPrefix: "/",
				SinceNs:    lastTsNs,
			})
			if err != nil {
				return fmt.Errorf("subscribe: %v", err)
			}

			for {
				resp, listenErr := stream.Recv()
				if listenErr == io.EOF {
					return nil
				}
				if listenErr != nil {
					return listenErr
				}

				if err := processEventFn(resp); err != nil {
					return fmt.Errorf("process %v: %v", resp, err)
				}
				lastTsNs = resp.TsNs
			}
		})
		if err != nil {
			glog.V(0).Infof("subscribing remote %s meta change: %v", peer, err)
			time.Sleep(1733 * time.Millisecond)
		}
	}
}

func (ma *MetaAggregator) readFilerStoreSignature(peer string) (sig int32, err error) {
	err = pb.WithFilerClient(peer, ma.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
		if err != nil {
			return err
		}
		sig = resp.Signature
		return nil
	})
	return
}

const (
	MetaOffsetPrefix = "Meta"
)

func (ma *MetaAggregator) readOffset(f *Filer, peer string, peerSignature int32) (lastTsNs int64, err error) {

	key := []byte(MetaOffsetPrefix + "xxxx")
	util.Uint32toBytes(key[len(MetaOffsetPrefix):], uint32(peerSignature))

	value, err := f.Store.KvGet(context.Background(), key)

	if err == ErrKvNotFound {
		glog.Warningf("readOffset %s not found", peer)
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("readOffset %s : %v", peer, err)
	}

	lastTsNs = int64(util.BytesToUint64(value))

	glog.V(0).Infof("readOffset %s : %d", peer, lastTsNs)

	return
}

func (ma *MetaAggregator) updateOffset(f *Filer, peer string, peerSignature int32, lastTsNs int64) (err error) {

	key := []byte(MetaOffsetPrefix + "xxxx")
	util.Uint32toBytes(key[len(MetaOffsetPrefix):], uint32(peerSignature))

	value := make([]byte, 8)
	util.Uint64toBytes(value, uint64(lastTsNs))

	err = f.Store.KvPut(context.Background(), key, value)

	if err != nil {
		return fmt.Errorf("updateOffset %s : %v", peer, err)
	}

	glog.V(4).Infof("updateOffset %s : %d", peer, lastTsNs)

	return
}
