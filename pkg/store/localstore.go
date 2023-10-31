package store

import (
	"context"
	"github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/role"
	"github.com/oppslink/protocol/logger"
	"os"
	"sync"
)

type roleName string

type LocalStore struct {
	DistPath string

	Role map[roleName]*role.Role

	lock       sync.RWMutex
	globalLock sync.Mutex

	dataSync chan struct{}
}

func NewLocalStore() *LocalStore {
	localStore := &LocalStore{
		DistPath: config.LocalStorePath,
		Role:     make(map[roleName]*role.Role),
		lock:     sync.RWMutex{},
	}

	localStore.loadDistFile()

	return localStore
}

func (ls *LocalStore) Create(ctx context.Context, v *role.Role) error {
	return nil
}

func (ls *LocalStore) Update(ctx context.Context, old *role.Role, new *role.Role) (*role.Role, error) {
	return nil, nil
}

func (ls *LocalStore) Delete(ctx context.Context, v any) (int64, error) {
	return 0, nil
}

func (ls *LocalStore) Get(ctx context.Context, k string) ([]*role.Role, error) {
	return nil, nil
}

func (ls *LocalStore) List(ctx context.Context, key string) ([]*role.Role, error) {
	return nil, nil
}

func (ls *LocalStore) loadDistFile() {
	ls.lock.Lock()
	defer ls.lock.Unlock()

	if _, err := os.Stat(ls.DistPath); os.IsNotExist(err) {
		// create empty file
		emptyFile, createErr := os.Create(ls.DistPath)
		if createErr != nil {
			logger.Errorw("Load Roles File Error!", createErr)
			return
		}
		defer emptyFile.Close()
		logger.Infow("Load Roles File Success!")
	} else {
		logger.Infow("Load Roles File Success!", "path", config.CasbinCsvPath)
	}
}
