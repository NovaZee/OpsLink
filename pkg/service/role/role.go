package role

import "sync"

type Role struct {
	Id       int64  `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"password"`

	Token string `json:"token" yaml:"token"`
	mu    *sync.Mutex
}

func NewRole(id int64, name string, pwd string) *Role {
	return &Role{
		Id:       id,
		Name:     name,
		Password: pwd,
	}
}

func (r *Role) LoadToken(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Token = token
}
