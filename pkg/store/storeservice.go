package store

import (
	"context"
	opsconfig "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/role"
)

type StoreService interface {
	Create(ctx context.Context, v *role.Role) error
	Update(ctx context.Context, v *role.Role, a *role.Role) (*role.Role, error)
	Delete(ctx context.Context, v any) (int64, error)
	Get(ctx context.Context, k string) ([]*role.Role, error)
	List(ctx context.Context, k string) ([]*role.Role, error)
}

// NewStoreService 创建 StoreService 实例，根据配置选择合适的存储中间件
func NewStoreService(config *opsconfig.OpsLinkConfig) (StoreService, error) {
	if len(config.EtcdConfig.Endpoint) != 0 {
		return &LocalStore{}, nil
	} else {
		return &V3Store{}, nil
	}
}
