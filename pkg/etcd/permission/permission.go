package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/etcd"
	"github.com/denovo/permission/pkg/service/role"
	"github.com/oppslink/protocol/logger"
	"sync"
	"time"
)

type PermissionEtcdClient struct {
	etcd.PermissionClient

	closed bool
	mu     sync.Mutex

	client *clientv3.Client
	logger logger.Logger
}

func InitEtcd(cfg *config.Config) (*PermissionEtcdClient, error) {
	// 创建一个 etcd 客户端连接
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.EtcdConfig.Endpoint, // etcd节点地址
		DialTimeout: time.Duration(cfg.EtcdConfig.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return newEtcdEtcdPermissonClient(etcd, logger.GetLogger()), nil
}

// NewEtcdEtcdPermissonClient  create a register based on etcd
func newEtcdEtcdPermissonClient(c *clientv3.Client, opsLog logger.Logger) *PermissionEtcdClient {
	return &PermissionEtcdClient{
		client: c,
		logger: opsLog,
	}
}

func (pc *PermissionEtcdClient) SetPermissionPolicy(v any) error {
	defer pc.client.Close()
	key := key(v)
	result, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, e := pc.client.KV.Put(context.TODO(), key, string(result))
	if e != nil {
		return e
	}
	return nil
}

func (pc *PermissionEtcdClient) GetPermissionPolicy(v any) ([]*role.Role, error) {
	defer pc.client.Close()
	key := key(v)
	get, err := pc.client.KV.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	// 处理获取的结果
	var roles []*role.Role
	for _, kv := range get.Kvs {
		var r *role.Role
		// fmt.Printf("键：%s，值：%s\n", kv.Key, kv.Value)
		if err2 := json.Unmarshal(kv.Value, r); err != nil {
			return nil, err2
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func (pc *PermissionEtcdClient) DeletePermissionPolicy(v any) (int64, error) {
	defer pc.client.Close()
	key := key(v)
	response, err := pc.client.KV.Delete(context.TODO(), key)
	if err != nil {
		return 0, err
	}
	return response.Deleted, nil
}

func (pc *PermissionEtcdClient) UpdatePermissionPolicy(old any, new any) (*role.Role, error) {
	ctx := context.TODO()
	getKey := key(old)
	getResp, err := pc.client.KV.Get(ctx, getKey)
	if err != nil {
		return nil, err
	}
	// 假设键存在
	if len(getResp.Kvs) > 0 {
		// 获取当前值
		//currentValue := getResp.Kvs[0].Value
		// 2. 修改值
		putKey := key(new)
		newValue, err2 := json.Marshal(new)
		if err2 != nil {
			return nil, err2
		}
		// 3. 更新键的值
		_, err3 := pc.client.KV.Put(ctx, putKey, string(newValue))
		if err3 != nil {
			return nil, err3
		}
		r := new.(*role.Role)
		return r, nil
	} else {
		return nil, errors.New("键不存在")
	}
}

func key(v any) (k string) {
	switch t := v.(type) {
	case *role.Role:
		k = etcd.RoleKey + t.FrontRole.Name
		return
	default:
		k = fmt.Sprintf("Unhandled type: %T", v)
		return
	}
}

func (pc *PermissionEtcdClient) close() {
	// 使用互斥锁确保关闭操作只执行一次
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.closed {
		// 执行关闭逻辑，释放资源等
		err := pc.client.Close()
		if err != nil {
			return
		}
		pc.closed = true
	}
}
