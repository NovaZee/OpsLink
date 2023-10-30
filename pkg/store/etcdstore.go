package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/role"
)

type V3Store struct {
	Backend Client
	config  config.OpsLinkConfig
}

func (r *V3Store) Create(ctx context.Context, v *role.Role) error {
	key := convertKey(v)
	result, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err1 := r.Backend.Create(ctx, key, string(result))
	if err1 != nil {
		return err1
	}
	return nil
}

func (r *V3Store) Update(ctx context.Context, old *role.Role, new *role.Role) (*role.Role, error) {
	//TODO implement me
	getKey := convertKey(old)
	getResp, err := r.Get(ctx, getKey)
	if err != nil {
		return nil, err
	}
	// 假设键存在
	if len(getResp) > 0 {
		oldValue := getResp[0]
		if oldValue.Password == new.Password && oldValue.Name == new.Password {
			return oldValue, nil
		}
		err3 := r.Create(ctx, new)
		if err3 != nil {
			return nil, err3
		}
		return new, nil
	} else {
		return nil, errors.New("键不存在")
	}
}

func (r *V3Store) Delete(ctx context.Context, v any) (int64, error) {
	key := convertKey(v)
	i, err := r.Backend.Delete(ctx, key)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

func (r *V3Store) Get(ctx context.Context, k string) ([]*role.Role, error) {
	k1 := convertKey(k)
	get, err := r.Backend.Get(ctx, k1)
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

func (r *V3Store) List(ctx context.Context, key string) ([]*role.Role, error) {
	list, err := r.Backend.List(ctx, key)
	if err != nil {
		return nil, err
	}
	var roles []*role.Role
	for _, kv := range list {
		var r *role.Role
		// fmt.Printf("键：%s，值：%s\n", kv.Key, kv.Value)
		if err2 := json.Unmarshal(kv.Value, r); err != nil {
			return nil, err2
		}
		roles = append(roles, r)
	}
	return roles, nil
}

// Watch simple crud don't need watch
func (r *V3Store) Watch(ctx context.Context, v any, a any) error {
	//TODO implement me
	panic("implement me")
}

func convertKey(v any) (k string) {
	switch t := v.(type) {
	case *role.Role:
		k = config.RoleKey + t.FrontRole.Name
		return
	default:
		k = fmt.Sprintf("Unhandled type: %T", v)
		return
	}
}
