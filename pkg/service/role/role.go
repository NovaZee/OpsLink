package role

import (
	"math/rand"
	"sync"
	"time"
)

type FrontRole struct {
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"password"`
}

type Role struct {
	Id        int64 `json:"id" yaml:"id"`
	FrontRole `json:"front_role" yaml:"front_role"`
	mu        *sync.Mutex
}

func NewRole(frontRole FrontRole) (role *Role) {
	role = &Role{
		FrontRole: frontRole,
	}
	rand.Seed(time.Now().UnixNano())
	role.Id = rand.Int63()
	return
}

// 实现 Message 接口的 Reset 方法
func (r *Role) Reset() {
	// 在此方法中重置 Role 结构的字段
}

// 实现 Message 接口的 String 方法
func (r *Role) String() string {
	// 返回 Role 结构的字符串表示形式
	// 通常用于调试目的
	return "Role as a string"
}

// 实现 Message 接口的 ProtoMessage 方法
func (r *Role) ProtoMessage() {
	// 在此方法中返回 Role 结构的 proto.Message
	// 通常是消息自身，因此返回 r 即可
}
