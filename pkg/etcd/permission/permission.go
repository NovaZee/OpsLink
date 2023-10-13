package etcd

import (
	"github.com/denovo/permission/pkg/etcd"
	"github.com/denovo/permission/pkg/service/role"
	"github.com/oppslink/protocol/logger"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func NewPermissionClient(etcdClient *clientv3.Client) {

}

const HTTP = "http://"

const CasbinRuleKey = "casbin_policy"
const RoleKey = "role_key"

type PermissionEtcdClient struct {
	etcd.PermissionClient

	client *clientv3.Client
	logger logger.Logger
}

func InitEtcd(endpoint []string) (*PermissionEtcdClient, error) {
	// 创建一个 etcd 客户端连接
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint, // etcd节点地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return NewEtcdEtcdPermissonClient(etcd, logger.GetLogger()), nil
}

// NewEtcdEtcdPermissonClient  create a register based on etcd
func NewEtcdEtcdPermissonClient(c *clientv3.Client, opsLog logger.Logger) *PermissionEtcdClient {
	return &PermissionEtcdClient{
		client: c,
		logger: opsLog,
	}
}

func (pc *PermissionEtcdClient) SetPermissionPolicy(v any) error {
	var key string
	if r, ok := v.(*role.Role); ok {
		key = RoleKey + r.FrontRole.Name
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

func (c3 *EtcdRegister) GetkV(key string) error {
	_, e := c3.Get(context.TODO(), key)
	if e != nil {
		return e
	}
	return nil
}

func (c3 *KvClient) DelkV(v any) error {
	return nil
}
