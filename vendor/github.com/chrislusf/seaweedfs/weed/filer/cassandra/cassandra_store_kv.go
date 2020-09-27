package cassandra

import (
	"context"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/filer"
	"github.com/gocql/gocql"
)

func (store *CassandraStore) KvPut(ctx context.Context, key []byte, value []byte) (err error) {
	dir, name := genDirAndName(key)

	if err := store.session.Query(
		"INSERT INTO filemeta (directory,name,meta) VALUES(?,?,?) USING TTL ? ",
		dir, name, value, 0).Exec(); err != nil {
		return fmt.Errorf("kv insert: %s", err)
	}

	return nil
}

func (store *CassandraStore) KvGet(ctx context.Context, key []byte) (data []byte, err error) {
	dir, name := genDirAndName(key)

	if err := store.session.Query(
		"SELECT meta FROM filemeta WHERE directory=? AND name=?",
		dir, name).Consistency(gocql.One).Scan(&data); err != nil {
		if err != gocql.ErrNotFound {
			return nil, filer.ErrKvNotFound
		}
	}

	if len(data) == 0 {
		return nil, filer.ErrKvNotFound
	}

	return data, nil
}

func (store *CassandraStore) KvDelete(ctx context.Context, key []byte) (err error) {
	dir, name := genDirAndName(key)

	if err := store.session.Query(
		"DELETE FROM filemeta WHERE directory=? AND name=?",
		dir, name).Exec(); err != nil {
		return fmt.Errorf("kv delete: %v", err)
	}

	return nil
}

func genDirAndName(key []byte) (dir string, name string) {
	for len(key) < 8 {
		key = append(key, 0)
	}

	dir = string(key[:8])
	name = string(key[8:])

	return
}
