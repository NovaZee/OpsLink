package casbin

import (
	"github.com/casbin/casbin"
)

type DefaultPolicy struct {
	Policy
	e *casbin.Enforcer
}

func (c *DefaultPolicy) Add(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		result := c.e.AddPolicy(casbinModel.RoleId, casbinModel.Path, casbinModel.Method)
		return result
	}
	return false
}
func (c *DefaultPolicy) Update(a any) bool {
	if _, ok := a.([]*CasbinModel); ok {
		// 遍历集合中的每个 CasbinModel 并添加策略
		return true
	}
	return false
}
func (c *DefaultPolicy) Delete(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		result := c.e.RemovePolicy(casbinModel.RoleId, casbinModel.Path, casbinModel.Method)
		return result
	}
	return false
}

func NewDefaultPolicy(cba *CasbinAdapter) (*DefaultPolicy, error) {
	enforcer, _ := cba.Casbin()
	return &DefaultPolicy{
		e: enforcer,
	}, nil
}
