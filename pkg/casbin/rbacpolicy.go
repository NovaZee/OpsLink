package casbin

import (
	"github.com/casbin/casbin"
)

type RbacPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *RbacPolicy) AddGroupingPolicy(role string, group string) bool {
	return false
}
func (c *RbacPolicy) Update(a any) bool {
	return false
}
func (c *RbacPolicy) Delete(a any) bool {
	return false
}

func NewRbacPolicy(cba *CasbinAdapter) (*RbacPolicy, error) {
	enforcer, _ := cba.Casbin()
	return &RbacPolicy{
		e: enforcer,
	}, nil
}
