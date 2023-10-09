package casbin

import (
	"github.com/casbin/casbin/v2"
)

type DefaultPolicy struct {
	CasbinI
	e *casbin.Enforcer
}

func (c *DefaultPolicy) Add() bool {
	filteredPolicy := c.e.GetFilteredPolicy(0, "alice")
}
func (c *DefaultPolicy) Update() bool {
	return nil
}
func (c *DefaultPolicy) Delete() bool {
	return nil
}

func NewDefaultPolicy(cba CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.NewCasbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
