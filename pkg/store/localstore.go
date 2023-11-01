package store

import (
	"context"
	"github.com/denovo/permission/config"
	pb "github.com/denovo/permission/pkg/protoc"
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

func NewLocalStore() *LocalStore {
	localStore := &LocalStore{
		DistPath: config.LocalStorePath,
		//Role:     &role.RoleMap{Roles: make(map[string]*role.RoleEntry)},
		lock:     sync.RWMutex{},
		dataSync: make(chan struct{}),
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
		ls.WriteData()
		defer emptyFile.Close()
		logger.Infow("Load Roles File Success!")
	} else {
		//todo：Marshal 二进制
		ls.ReadData()
		logger.Infow("Load Roles File Success!", "path", ls.DistPath)
	}
}
func (ls *LocalStore) ReadData() {
	serializedData, err := os.ReadFile(ls.DistPath)
	if err != nil {
		// 处理错误
		return
	}
	// 反序列化二进制数据
	rs := &pb.RolesSlice{}
	err = proto.Unmarshal(serializedData, rs)
	if err != nil {
		// 处理错误
		return
	}
	println(rs.String())
}

func (ls *LocalStore) WriteData() {
	newRole1 := role.NewRole("1", "2")
	newRole2 := role.NewRole("3", "4")
	slice := role.NewSlice(newRole1, newRole2)
	// 序列化RoleMap消息为二进制数据
	serializedRoleMap, err := proto.Marshal(slice)
	if err != nil {
		// 处理错误
	}

	err = os.WriteFile(ls.DistPath, serializedRoleMap, 0644)
	if err != nil {
		// 处理错误
		return
	}
}

func (ls *LocalStore) dealSyncData() {
	go func() {

	}()
}

func (ls *LocalStore) Stop() {
	go func() {

	}()
}
