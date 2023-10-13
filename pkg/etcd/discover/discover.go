package etcd

import (
	"github.com/denovo/permission/pkg/etcd"
)

type DiscoveryModule struct {
	etcdClient etcd.DiscoveryClient
}
