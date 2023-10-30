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
