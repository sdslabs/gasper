package exclusive_locks

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/wdclient"
)

const (
	RenewInteval     = 4 * time.Second
	SafeRenewInteval = 3 * time.Second
	InitLockInteval  = 1 * time.Second
	AdminLockName    = "admin"
)

type ExclusiveLocker struct {
	masterClient *wdclient.MasterClient
	token        int64
	lockTsNs     int64
	isLocking    bool
}

func NewExclusiveLocker(masterClient *wdclient.MasterClient) *ExclusiveLocker {
	return &ExclusiveLocker{
		masterClient: masterClient,
	}
}
func (l *ExclusiveLocker) IsLocking() bool {
	return l.isLocking
}

func (l *ExclusiveLocker) GetToken() (token int64, lockTsNs int64) {
	for time.Unix(0, atomic.LoadInt64(&l.lockTsNs)).Add(SafeRenewInteval).Before(time.Now()) {
		// wait until now is within the safe lock period, no immediate renewal to change the token
		time.Sleep(100 * time.Millisecond)
	}
	return atomic.LoadInt64(&l.token), atomic.LoadInt64(&l.lockTsNs)
}

func (l *ExclusiveLocker) RequestLock() {
	if l.isLocking {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// retry to get the lease
	for {
		if err := l.masterClient.WithClient(func(client master_pb.SeaweedClient) error {
			resp, err := client.LeaseAdminToken(ctx, &master_pb.LeaseAdminTokenRequest{
				PreviousToken:    atomic.LoadInt64(&l.token),
				PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
				LockName:         AdminLockName,
			})
			if err == nil {
				atomic.StoreInt64(&l.token, resp.Token)
				atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs)
			}
			return err
		}); err != nil {
			// println("leasing problem", err.Error())
			time.Sleep(InitLockInteval)
		} else {
			break
		}
	}

	l.isLocking = true

	// start a goroutine to renew the lease
	go func() {
		ctx2, cancel2 := context.WithCancel(context.Background())
		defer cancel2()

		for l.isLocking {
			if err := l.masterClient.WithClient(func(client master_pb.SeaweedClient) error {
				resp, err := client.LeaseAdminToken(ctx2, &master_pb.LeaseAdminTokenRequest{
					PreviousToken:    atomic.LoadInt64(&l.token),
					PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
					LockName:         AdminLockName,
				})
				if err == nil {
					atomic.StoreInt64(&l.token, resp.Token)
					atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs)
					// println("ts", l.lockTsNs, "token", l.token)
				}
				return err
			}); err != nil {
				glog.Errorf("failed to renew lock: %v", err)
				return
			} else {
				time.Sleep(RenewInteval)
			}

		}
	}()

}

func (l *ExclusiveLocker) ReleaseLock() {
	l.isLocking = false

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l.masterClient.WithClient(func(client master_pb.SeaweedClient) error {
		client.ReleaseAdminToken(ctx, &master_pb.ReleaseAdminTokenRequest{
			PreviousToken:    atomic.LoadInt64(&l.token),
			PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
			LockName:         AdminLockName,
		})
		return nil
	})
	atomic.StoreInt64(&l.token, 0)
	atomic.StoreInt64(&l.lockTsNs, 0)
}
