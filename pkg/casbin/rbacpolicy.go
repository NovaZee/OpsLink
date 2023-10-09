package casbin

import (
	"github.com/casbin/casbin/v2"
)

type RbacPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *RbacPolicy) Add() bool {
	return nil
}
func (c *RbacPolicy) Update() bool {
	return nil
}
func (c *RbacPolicy) Delete() bool {
	return nil
}

func NewRbacPolicy(cba CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.NewCasbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
