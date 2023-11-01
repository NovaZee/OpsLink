package role

import (
	"math/rand"
	"sync"
	"time"
)

type RolesSlice struct {
	Roles []*Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles"`
}

type Role struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id" yaml:"id"`

	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name" yaml:"name"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password" yaml:"password"`

	mu *sync.Mutex
}

func (r RolesSlice) Reset() {
	//TODO implement me
	panic("implement me")
}

func (r RolesSlice) String() string {
	//TODO implement me
	panic("implement me")
}

func (r RolesSlice) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

func NewRole(name, password string) (role *Role) {
	role = &Role{
		Name:     name,
		Password: password,
	}
	rand.Seed(time.Now().UnixNano())
	role.Id = rand.Int63()
	return
}

func NewSlice(roles ...*Role) *RolesSlice {
	slices := make([]*Role, 0, 100)
	for _, role := range roles {
		slices = append(slices, role)
	}
	return &RolesSlice{
		Roles: slices,
	}
}
