package weed_server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/chrislusf/seaweedfs/weed/util/grace"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/stats"
	"github.com/chrislusf/seaweedfs/weed/util"

	"github.com/chrislusf/seaweedfs/weed/filer"
	_ "github.com/chrislusf/seaweedfs/weed/filer/cassandra"
	_ "github.com/chrislusf/seaweedfs/weed/filer/elastic/v7"
	_ "github.com/chrislusf/seaweedfs/weed/filer/etcd"
	_ "github.com/chrislusf/seaweedfs/weed/filer/leveldb"
	_ "github.com/chrislusf/seaweedfs/weed/filer/leveldb2"
	_ "github.com/chrislusf/seaweedfs/weed/filer/mongodb"
	_ "github.com/chrislusf/seaweedfs/weed/filer/mysql"
	_ "github.com/chrislusf/seaweedfs/weed/filer/postgres"
	_ "github.com/chrislusf/seaweedfs/weed/filer/redis"
	_ "github.com/chrislusf/seaweedfs/weed/filer/redis2"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/notification"
	_ "github.com/chrislusf/seaweedfs/weed/notification/aws_sqs"
	_ "github.com/chrislusf/seaweedfs/weed/notification/gocdk_pub_sub"
	_ "github.com/chrislusf/seaweedfs/weed/notification/google_pub_sub"
	_ "github.com/chrislusf/seaweedfs/weed/notification/kafka"
	_ "github.com/chrislusf/seaweedfs/weed/notification/log"
	"github.com/chrislusf/seaweedfs/weed/security"
)

type FilerOption struct {
	Masters            []string
	Collection         string
	DefaultReplication string
	DisableDirListing  bool
	MaxMB              int
	DirListingLimit    int
	DataCenter         string
	DefaultLevelDbDir  string
	DisableHttp        bool
	Host               string
	Port               uint32
	recursiveDelete    bool
	Cipher             bool
	Filers             []string
}

type FilerServer struct {
	option         *FilerOption
	secret         security.SigningKey
	filer          *filer.Filer
	grpcDialOption grpc.DialOption

	// metrics read from the master
	metricsAddress     string
	metricsIntervalSec int

	// notifying clients
	listenersLock sync.Mutex
	listenersCond *sync.Cond

	brokers     map[string]map[string]bool
	brokersLock sync.Mutex
}

func NewFilerServer(defaultMux, readonlyMux *http.ServeMux, option *FilerOption) (fs *FilerServer, err error) {

	fs = &FilerServer{
		option:         option,
		grpcDialOption: security.LoadClientTLS(util.GetViper(), "grpc.filer"),
		brokers:        make(map[string]map[string]bool),
	}
	fs.listenersCond = sync.NewCond(&fs.listenersLock)

	if len(option.Masters) == 0 {
		glog.Fatal("master list is required!")
	}

	fs.filer = filer.NewFiler(option.Masters, fs.grpcDialOption, option.Host, option.Port, option.Collection, option.DefaultReplication, func() {
		fs.listenersCond.Broadcast()
	})
	fs.filer.Cipher = option.Cipher

	fs.maybeStartMetrics()

	go fs.filer.KeepConnectedToMaster()

	v := util.GetViper()
	if !util.LoadConfiguration("filer", false) {
		v.Set("leveldb2.enabled", true)
		v.Set("leveldb2.dir", option.DefaultLevelDbDir)
		_, err := os.Stat(option.DefaultLevelDbDir)
		if os.IsNotExist(err) {
			os.MkdirAll(option.DefaultLevelDbDir, 0755)
		}
		glog.V(0).Infof("default to create filer store dir in %s", option.DefaultLevelDbDir)
	}
	util.LoadConfiguration("notification", false)

	fs.option.recursiveDelete = v.GetBool("filer.options.recursive_delete")
	v.SetDefault("filer.options.buckets_folder", "/buckets")
	fs.filer.DirBucketsPath = v.GetString("filer.options.buckets_folder")
	fs.filer.FsyncBuckets = v.GetStringSlice("filer.options.buckets_fsync")
	fs.filer.LoadConfiguration(v)

	notification.LoadConfiguration(v, "notification.")

	handleStaticResources(defaultMux)
	if !option.DisableHttp {
		defaultMux.HandleFunc("/", fs.filerHandler)
	}
	if defaultMux != readonlyMux {
		readonlyMux.HandleFunc("/", fs.readonlyFilerHandler)
	}

	fs.filer.AggregateFromPeers(fmt.Sprintf("%s:%d", option.Host, option.Port), option.Filers)

	fs.filer.LoadBuckets()

	grace.OnInterrupt(func() {
		fs.filer.Shutdown()
	})

	return fs, nil
}

func (fs *FilerServer) maybeStartMetrics() {

	for _, master := range fs.option.Masters {
		_, err := pb.ParseFilerGrpcAddress(master)
		if err != nil {
			glog.Fatalf("invalid master address %s: %v", master, err)
		}
	}

	isConnected := false
	var readErr error
	for !isConnected {
		for _, master := range fs.option.Masters {
			fs.metricsAddress, fs.metricsIntervalSec, readErr = readFilerConfiguration(fs.grpcDialOption, master)
			if readErr == nil {
				isConnected = true
			} else {
				time.Sleep(7 * time.Second)
			}
		}
	}

	go stats.LoopPushingMetric("filer", stats.SourceName(fs.option.Port), fs.metricsAddress, fs.metricsIntervalSec)
}

func readFilerConfiguration(grpcDialOption grpc.DialOption, masterAddress string) (metricsAddress string, metricsIntervalSec int, err error) {
	err = operation.WithMasterServerClient(masterAddress, grpcDialOption, func(masterClient master_pb.SeaweedClient) error {
		resp, err := masterClient.GetMasterConfiguration(context.Background(), &master_pb.GetMasterConfigurationRequest{})
		if err != nil {
			return fmt.Errorf("get master %s configuration: %v", masterAddress, err)
		}
		metricsAddress, metricsIntervalSec = resp.MetricsAddress, int(resp.MetricsIntervalSeconds)
		return nil
	})
	return
}
