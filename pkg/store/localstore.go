package store

import (
	"context"
	"github.com/denovo/permission/pkg/service/role"
)

type LocalStore struct {
}

func (r *LocalStore) Create(ctx context.Context, v *role.Role) error {
	return nil
}

func (r *LocalStore) Update(ctx context.Context, old *role.Role, new *role.Role) (*role.Role, error) {
	return nil, nil
}

func (r *LocalStore) Delete(ctx context.Context, v any) (int64, error) {
	return 0, nil
}

func (r *LocalStore) Get(ctx context.Context, k string) ([]*role.Role, error) {
	return nil, nil
}

func (r *LocalStore) List(ctx context.Context, key string) ([]*role.Role, error) {
	return nil, nil
}
