package store

import (
	"context"
	"encoding/json"
	"errors"
	opsconfig "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/role"
)

type V3Store struct {
	Backend Client

	config opsconfig.EtcdConfig
}

func (r *V3Store) Stop() {
	//TODO implement me
	panic("implement me")
}

func (r *V3Store) Start() {
	//TODO implement me
	panic("implement me")
}

func (r *V3Store) Create(ctx context.Context, v *role.Role) error {
	key := ConvertKey(v)
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
	getKey := ConvertKey(old)
	getResp, err := r.Get(ctx, getKey)
	if err != nil {
		return nil, err
	}
	// 假设键存在
	if getResp != nil {
		oldValue := getResp
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
	key := ConvertKey(v)
	i, err := r.Backend.Delete(ctx, key)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

func (r *V3Store) Get(ctx context.Context, k string) (*role.Role, error) {
	k1 := ConvertKey(k)
	get, err := r.Backend.Get(ctx, k1)
	// 处理获取的结果
	var roles *role.Role
	if err2 := json.Unmarshal(get[0].Value, &roles); err != nil {
		return nil, err2
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
