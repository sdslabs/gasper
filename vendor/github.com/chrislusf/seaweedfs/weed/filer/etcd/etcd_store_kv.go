package etcd

import (
	"context"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/filer"
)

func (store *EtcdStore) KvPut(ctx context.Context, key []byte, value []byte) (err error) {

	_, err = store.client.Put(ctx, string(key), string(value))

	if err != nil {
		return fmt.Errorf("kv put: %v", err)
	}

	return nil
}

func (store *EtcdStore) KvGet(ctx context.Context, key []byte) (value []byte, err error) {

	resp, err := store.client.Get(ctx, string(key))

	if err != nil {
		return nil, fmt.Errorf("kv get: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, filer.ErrKvNotFound
	}

	return resp.Kvs[0].Value, nil
}

func (store *EtcdStore) KvDelete(ctx context.Context, key []byte) (err error) {

	_, err = store.client.Delete(ctx, string(key))

	if err != nil {
		return fmt.Errorf("kv delete: %v", err)
	}

	return nil
}
