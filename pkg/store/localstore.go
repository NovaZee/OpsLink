package store

import (
	"context"
	"github.com/denovo/permission/config"
	opslink "github.com/denovo/permission/pkg/protoc/opslink"
	"github.com/denovo/permission/pkg/service/role"
	"github.com/golang/protobuf/proto"
	"github.com/oppslink/protocol/logger"
	"os"
	"sync"
)

type LocalStore struct {
	DistPath string

	Role *role.RolesSlice

	lock       sync.RWMutex
	globalLock sync.Mutex

	dataSync chan struct{}
}

func NewLocalStore() (*LocalStore, error) {
	localStore := &LocalStore{
		DistPath: config.LocalStorePath,
		Role:     &role.RolesSlice{Roles: []*role.Role{}},
		lock:     sync.RWMutex{},
		dataSync: make(chan struct{}),
	}

	err := localStore.loadDistFile()
	if err != nil {
		return nil, err
	}

	return localStore, nil
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

func (ls *LocalStore) loadDistFile() error {
	ls.lock.Lock()
	defer ls.lock.Unlock()

	if _, err := os.Stat(ls.DistPath); os.IsNotExist(err) {
		// create empty file
		emptyFile, createErr := os.Create(ls.DistPath)
		if createErr != nil {
			logger.Errorw("Load Roles File Error!", createErr)
			return createErr
		}
		//ls.WriteData()
		defer emptyFile.Close()
		logger.Infow("Load Roles File Success!")
	} else {
		err := ls.ReadData()
		if err != nil {
			return err
		}
		logger.Infow("Load Roles File Success!", "path", ls.DistPath)
	}
	return nil
}
func (ls *LocalStore) ReadData() error {
	serializedData, err := os.ReadFile(ls.DistPath)
	if err != nil {
		return err
	}
	// 反序列化二进制数据
	rs := &opslink.RolesSlice{}
	err = proto.Unmarshal(serializedData, rs)
	if err != nil {
		return err
	}
	ls.ConvertRoles(rs)
	return nil
}

func (ls *LocalStore) WriteData(rs *role.RolesSlice) error {
	// 序列化RoleMap消息为二进制数据
	serializedRoleMap, err := proto.Marshal(rs)
	if err != nil {
		return err
	}
	err = os.WriteFile(ls.DistPath, serializedRoleMap, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LocalStore) dealSyncData() {
	go func() {

	}()
}

func (ls *LocalStore) Stop() {
	go func() {

	}()
}

// ConvertRoles pb struct convert to runtime role struct
func (ls *LocalStore) ConvertRoles(pbRoles *opslink.RolesSlice) *LocalStore {
	for _, r := range pbRoles.GetRoles() {
		ls.Role.Roles = append(ls.Role.Roles, &role.Role{
			Name:     r.GetName(),
			Password: r.GetPassword(),
			Id:       r.GetId(),
		})
	}
	return ls
}
