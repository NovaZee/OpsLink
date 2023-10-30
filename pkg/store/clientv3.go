package store

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

var (
	clientTimeout    = 10 * time.Second
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second
)

type Client interface {
	Create(ctx context.Context, key string, result string) error
	Update(ctx context.Context, object any) error
	Delete(ctx context.Context, key string) (int64, error)
	DeleteKVP(ctx context.Context, object any) error
	Get(ctx context.Context, key string) ([]*mvccpb.KeyValue, error)
	List(ctx context.Context, key string) ([]*mvccpb.KeyValue, error)
	EnsureInitialized() error

	//Apply(ctx context.Context, object any) (*model.KVPair, error)
	//Watch(ctx context.Context, list any, revision string) (WatchInterface, error)
	//Clean() error

}

type etcdV3Client struct {
	etcdClient *clientv3.Client
}

func (e etcdV3Client) Create(ctx context.Context, key string, value string) error {
	_, err := e.etcdClient.KV.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (e etcdV3Client) Delete(ctx context.Context, key string) (int64, error) {
	response, err := e.etcdClient.Delete(ctx, key)
	if err != nil {
		return 0, err
	}
	return response.Deleted, nil
}

func (e etcdV3Client) Get(ctx context.Context, k string) ([]*mvccpb.KeyValue, error) {
	get, err := e.etcdClient.KV.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	return get.Kvs, nil
}

func (e etcdV3Client) List(ctx context.Context, key string) ([]*mvccpb.KeyValue, error) {
	//find all value of prefix key
	get, err := e.etcdClient.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("list for prefix error:", err)
		return nil, err
	}
	return get.Kvs, nil
}

func (e etcdV3Client) Update(ctx context.Context, object any) error {
	panic("implement me")
}

func (e etcdV3Client) DeleteKVP(ctx context.Context, object any) error {
	panic("implement me")
}

func (e etcdV3Client) EnsureInitialized() error {
	panic("implement me")
}
