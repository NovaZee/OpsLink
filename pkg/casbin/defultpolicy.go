package casbin

import (
	"github.com/casbin/casbin/v2"
)

type DefaultPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *DefaultPolicy) Add(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		// a 是 CasbinModel 类型
		result, _ := c.e.AddPolicy(casbinModel)
		return result
	}
	return false
}
func (c *DefaultPolicy) Update() bool {
	return false
}
func (c *DefaultPolicy) Delete() bool {
	return false
}

func NewDefaultPolicy(cba *CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.Casbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
