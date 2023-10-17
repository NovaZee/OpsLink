package etcdv3

import (
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	_ "github.com/oppslink/protocol/logger"
	"time"
)

// client implements the client.Interface.
type SClient struct {
	//resources resourceInterface
	Backend Client
	config  config.OpsLinkConfig
}

func New(config *config.OpsLinkConfig) (Interface, error) {
	be, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return SClient{
		Backend: be,
		//resources: &resources{backend: be},
	}, nil
}

func NewClient(config *config.OpsLinkConfig) (Client, error) {
	// 创建一个 etcd 客户端连接
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdConfig.Endpoint, // etcd节点地址
		DialTimeout: time.Duration(config.EtcdConfig.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &etcdV3Client{etcdClient: etcd}, nil
}

// RolesCfg returns an interface for managing the Roles configuration resources.
func (c SClient) RolesCfg() RoleClientInterface {
	return RoleClient{client: c}
}
