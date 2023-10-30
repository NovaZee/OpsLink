package service

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/service/role"
	"sync"
)

type LocalStore struct {
	Role   *role.Role
	Policy []*casbin.CasbinModel

	lock       sync.RWMutex
	globalLock sync.Mutex
}

func NewLocalStore() *LocalStore {
	return nil
}
