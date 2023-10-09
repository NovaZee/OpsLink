package casbin

import (
	"github.com/casbin/casbin/v2"
)

type DefaultPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *DefaultPolicy) Add() bool {
	filteredPolicy := c.e.GetFilteredPolicy(0, "alice")
}
func (c *DefaultPolicy) Update() bool {
	return false
}
func (c *DefaultPolicy) Delete() bool {
	return false
}

func NewDefaultPolicy(cba CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.NewCasbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
