package store

import (
	"context"
	"fmt"
	opsconfig "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/role"
)

type StoreService interface {
	Create(ctx context.Context, v *role.Role) error
	Update(ctx context.Context, v *role.Role, a *role.Role) (*role.Role, error)
	Delete(ctx context.Context, v any) (int64, error)
	Get(ctx context.Context, k string) (*role.Role, error)
	List(ctx context.Context, k string) ([]*role.Role, error)
}

// NewStoreService 创建 StoreService 实例，根据配置选择合适的存储中间件
func NewStoreService(config *opsconfig.OpsLinkConfig) (StoreService, error) {
	if len(config.EtcdConfig.Endpoint) == 0 {
		store, err := NewLocalStore()
		if err != nil {
			return nil, err
		}
		return store, nil
	} else {
		client, err := NewClient(config)
		if err != nil {
			return nil, err
		}
		return &V3Store{Backend: client, config: config.EtcdConfig}, nil
	}
}

func ConvertKey(v any) (k string) {
	switch t := v.(type) {
	case *role.Role:
		k = opsconfig.RoleKey + t.Name
		return
	case string:
		k = opsconfig.RoleKey + t
		return
	default:
		k = fmt.Sprintf("Unhandled type: %T", v)
		return
	}
}
