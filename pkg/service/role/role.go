package role

import (
	"github.com/google/uuid"
	"sync"
)

type FrontRole struct {
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"password"`
}

type Role struct {
	Id        uuid.UUID `json:"id" yaml:"id"`
	FrontRole `json:"front_role" yaml:"front_role"`
	mu        *sync.Mutex
}

func NewRole(frontRole FrontRole) (role *Role) {
	role = &Role{
		FrontRole: frontRole,
	}
	role.Id = uuid.New()
	return
}
