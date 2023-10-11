package componment

import (
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"time"
)

var etcdClient *clientv3.Client

const HTTP = "http://"

const CasbinRuleKey = "casbin_policy"

func InitEtcd(config *config.Config) (*clientv3.Client, error) {
	// 创建一个 etcd 客户端连接
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{HTTP + config.EtcdConfig.Host + ":" + config.EtcdConfig.Port}, // etcd节点地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// 将客户端连接赋值给全局变量
	etcdClient = etcd
	return etcdClient, nil
}
