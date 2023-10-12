package componment

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/service/role"
	"time"
)

var etcdClient *clientv3.Client

const HTTP = "http://"

const CasbinRuleKey = "casbin_policy"
const RoleKey = "role_key"

type KvClient struct {
	clientv3.KV
}

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

func (c3 *KvClient) PutKv(v any) error {
	var key string
	if r, ok := v.(*role.Role); ok {
		key = RoleKey + r.Name
	}
	result, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, e := c3.Put(context.TODO(), key, string(result))
	if e != nil {
		return e
	}
	return nil
}

func (c3 *KvClient) GetkV(key string) error {
	_, e := c3.Get(context.TODO(), key)
	if e != nil {
		return e
	}
	return nil
}

func (c3 *KvClient) DelkV(v any) error {
	return nil
}
