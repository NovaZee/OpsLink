package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denovo/permission/pkg/etcd"
	"github.com/denovo/permission/pkg/service/role"
)

type RoleClient struct {
	client SClient
}

type RoleClientInterface interface {
	Create(ctx context.Context, v any) error
	Update(ctx context.Context, v any, a any) (*role.Role, error)
	Delete(ctx context.Context, v any) (int64, error)
	Get(ctx context.Context, k string) ([]*role.Role, error)
	List(ctx context.Context, v any, a any) error
	Watch(ctx context.Context, v any, a any) error
}

func (r RoleClient) Create(ctx context.Context, v any) error {
	key := convertKey(v)
	result, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err1 := r.client.Backend.Create(ctx, key, string(result))
	if err1 != nil {
		return err1
	}
	return nil
}

func (r RoleClient) Update(ctx context.Context, old any, new any) (*role.Role, error) {
	//TODO implement me
	getKey := convertKey(old)
	getResp, err := r.Get(ctx, getKey)
	if err != nil {
		return nil, err
	}
	// 假设键存在
	if len(getResp) > 0 {
		// 获取当前值
		//currentValue := getResp.Kvs[0].Value
		// 2. 修改值
		putKey := convertKey(new)
		newValue, err2 := json.Marshal(new)
		if err2 != nil {
			return nil, err2
		}
		// 3. 更新键的值
		err3 := r.List(ctx, putKey, string(newValue))
		if err3 != nil {
			return nil, err3
		}
		r := new.(*role.Role)
		return r, nil
	} else {
		return nil, errors.New("键不存在")
	}
}

func (r RoleClient) Delete(ctx context.Context, v any) (int64, error) {
	key := convertKey(v)
	i, err := r.client.Backend.Delete(ctx, key)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

func (r RoleClient) Get(ctx context.Context, k string) ([]*role.Role, error) {
	k1 := convertKey(k)
	get, err := r.client.Backend.Get(ctx, k1)
	// 处理获取的结果
	var roles []*role.Role
	for _, kv := range get {
		var r *role.Role
		// fmt.Printf("键：%s，值：%s\n", kv.Key, kv.Value)
		if err2 := json.Unmarshal(kv.Value, r); err != nil {
			return nil, err2
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func (r RoleClient) List(ctx context.Context, v any, a any) error {
	//TODO implement me
	panic("implement me")
}

func (r RoleClient) Watch(ctx context.Context, v any, a any) error {
	//TODO implement me
	panic("implement me")
}

func convertKey(v any) (k string) {
	switch t := v.(type) {
	case *role.Role:
		k = etcd.RoleKey + t.FrontRole.Name
		return
	default:
		k = fmt.Sprintf("Unhandled type: %T", v)
		return
	}
}
