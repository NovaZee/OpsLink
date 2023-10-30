package store

import (
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/config"
	_ "github.com/oppslink/protocol/logger"
	"time"
)

//
//func New(config *config.OpsLinkConfig) (Interface, error) {
//	//etcd未配置 跳过
//	if len(config.EtcdConfig.Endpoint) == 0 {
//		return nil, nil
//	}
//	be, err := NewClient(config)
//	//if err != nil {
//		return nil, err
//	}
//	return &SClient{
//		Backend: be,
//	}, nil
//}

func NewClient(config *config.OpsLinkConfig) (Client, error) {
	// create etcd3 connection
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdConfig.Endpoint, // etcd节点地址
		DialTimeout: time.Duration(config.EtcdConfig.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &etcdV3Client{etcdClient: etcd}, nil
}
