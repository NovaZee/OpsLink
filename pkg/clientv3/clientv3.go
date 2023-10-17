package etcdv3

import (
	"context"
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
	List(ctx context.Context, list any, revision string) error
	EnsureInitialized() error

	//Apply(ctx context.Context, object any) (*model.KVPair, error)
	//Watch(ctx context.Context, list any, revision string) (WatchInterface, error)
	//Clean() error

}

type etcdV3Client struct {
	etcdClient *clientv3.Client
}

func (e etcdV3Client) Create(ctx context.Context, key string, result string) error {
	//TODO 处理返回值
	_, err := e.etcdClient.KV.Put(ctx, key, result)
	if err != nil {
		return err
	}
	return nil
}

func (e etcdV3Client) Update(ctx context.Context, object any) error {
	//TODO implement me
	panic("implement me")
}

func (e etcdV3Client) Delete(ctx context.Context, key string) (int64, error) {
	response, err := e.etcdClient.Delete(ctx, key)
	if err != nil {
		return 0, err
	}
	return response.Deleted, nil
}

func (e etcdV3Client) DeleteKVP(ctx context.Context, object any) error {
	//TODO implement me
	panic("implement me")
}

func (e etcdV3Client) Get(ctx context.Context, k string) ([]*mvccpb.KeyValue, error) {
	//TODO implement me
	get, err := e.etcdClient.KV.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	return get.Kvs, nil
}

func (e etcdV3Client) List(ctx context.Context, list any, revision string) error {
	//TODO 处理返回值
	panic("implement me")
}

func (e etcdV3Client) EnsureInitialized() error {
	//TODO implement me
	panic("implement me")
}

//type resourceInterface interface {
//	Create(ctx context.Context, v any, kind string, in any) error
//	Update(ctx context.Context, v any, kind string, in any) error
//	Delete(ctx context.Context, v any, kind, ns, name string) error
//	Get(ctx context.Context, v any, kind, ns, name string) error
//	List(ctx context.Context, v any, kind, listkind string) error
//	Watch(ctx context.Context, v any, kind string) (watch.Interface, error)
//}
//
//// resources implements resourceInterface.
//type resources struct {
//	backend Client
//}
//
//func (r resources) Create(ctx context.Context, v any, kind string, in any) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r resources) Update(ctx context.Context, v any, kind string, in any) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r resources) Delete(ctx context.Context, v any, kind, ns, name string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r resources) Get(ctx context.Context, v any, kind, ns, name string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r resources) List(ctx context.Context, v any, kind, listkind string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r resources) Watch(ctx context.Context, v any, kind string) (watch.Interface, error) {
//	//TODO implement me
//	panic("implement me")
//}
