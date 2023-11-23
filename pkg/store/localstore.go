package store

import (
	"context"
	"github.com/denovo/permission/config"
	"github.com/denovo/permission/protoc"
	"github.com/denovo/permission/protoc/model"
	"github.com/golang/protobuf/proto"
	"github.com/oppslink/protocol/logger"
	"os"
	"sync"
)

type LocalStore struct {
	DistPath string

	//todo:结构缩减，重置protoc文件
	LocalRoles *model.RolesSlice

	lock            sync.RWMutex
	globalLock      sync.Mutex
	dataSyncCounter int
	dataSync        chan int

	protoc.Signal
}

func NewLocalStore() (*LocalStore, error) {
	localStore := &LocalStore{
		DistPath:        config.LocalStorePath,
		LocalRoles:      &model.RolesSlice{Roles: []*model.Role{}},
		lock:            sync.RWMutex{},
		dataSyncCounter: 0,
		dataSync:        make(chan int, 10),
	}

	err := localStore.loadDistFile()
	if err != nil {
		return nil, err
	}

	go localStore.dataSyncHandler()

	return localStore, nil
}

func (ls *LocalStore) Stop() {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	if ls.dataSyncCounter != 0 {
		ls.WriteData()
	}
	ls.dataSyncCounter = 0
}

func (ls *LocalStore) Create(_ context.Context, v *model.Role) error {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.LocalRoles.Roles = append(ls.LocalRoles.Roles, v)
	ls.dataSync <- 1
	return nil
}

func (ls *LocalStore) Update(_ context.Context, _ *model.Role, new *model.Role) (*model.Role, error) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	var uname = new.Name
	for i, r := range ls.LocalRoles.Roles {
		if r.Name == uname {
			ls.LocalRoles.Roles[i] = new
			break
		}
	}
	ls.dataSync <- 1
	return new, nil
}

func (ls *LocalStore) Delete(_ context.Context, v any) (int64, error) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	roles := v.(*model.Role)
	var uname = roles.Name
	var result int64
	for i, r := range ls.LocalRoles.Roles {
		if r.Name == uname {
			last := len(ls.LocalRoles.Roles) - 1
			//target moved to end
			ls.LocalRoles.Roles[i], ls.LocalRoles.Roles[last] = ls.LocalRoles.Roles[last], ls.LocalRoles.Roles[i]
			ls.LocalRoles.Roles = ls.LocalRoles.Roles[:len(ls.LocalRoles.Roles)-1]
			result += 1
			break
		}
	}
	if result != 0 {
		ls.dataSync <- 1
		return result, nil
	}
	return result, nil
}

func (ls *LocalStore) Get(_ context.Context, name string) (*model.Role, error) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	var roles = ls.LocalRoles.Roles
	for i := range roles {
		if roles[i].Name == name {
			return roles[i], nil
		}
	}
	return nil, nil
}

func (ls *LocalStore) List(_ context.Context, key string) ([]*model.Role, error) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	return ls.LocalRoles.Roles, nil
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
	rs := &model.RolesSlice{}
	err = proto.Unmarshal(serializedData, rs)
	if err != nil {
		return err
	}
	ls.ConvertRoles(rs)
	return nil
}

func (ls *LocalStore) WriteData() error {
	// 序列化RoleMap消息为二进制数据
	serializedRoleMap, err := proto.Marshal(ls.LocalRoles)
	if err != nil {
		return err
	}
	err = os.WriteFile(ls.DistPath, serializedRoleMap, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LocalStore) dataSyncHandler() {
	for {
		select {
		case <-ls.dataSync:
			ls.dataSyncCounter++
			if ls.dataSyncCounter == 10 {
				ls.lock.Lock()
				err := ls.WriteData()
				if err != nil {
					return
				}
				ls.dataSyncCounter = 0
				ls.lock.Unlock()
			}
		}
	}
}

// ConvertRoles pb struct convert to runtime role struct
func (ls *LocalStore) ConvertRoles(pbRoles *model.RolesSlice) *LocalStore {
	for _, r := range pbRoles.GetRoles() {
		ls.LocalRoles.Roles = append(ls.LocalRoles.Roles, r)
	}
	return ls
}
