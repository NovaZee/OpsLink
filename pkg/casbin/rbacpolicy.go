package casbin

import (
	"github.com/casbin/casbin/v2"
)

type RbacPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *RbacPolicy) Add(any2 any) bool {
	return false
}
func (c *RbacPolicy) Update() bool {
	return false
}
func (c *RbacPolicy) Delete() bool {
	return false
}

func NewRbacPolicy(cba CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.Casbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
