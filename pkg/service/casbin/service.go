package casbin

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/denovo/permission/protoc/model"
	"github.com/oppslink/protocol/logger"
)

type Casbin struct {
	Enforcer *casbin.Enforcer
}

// EnforcerProvider 接口定义
type EnforcerProvider interface {
	GetEnforcer(modelConf string) (*casbin.Enforcer, error)
}

type Policy interface {
	BindingRoles(gModel *model.GModel) (bool, error)
	UnBindingRoles(gModel *model.GModel) (bool, error)
	AddPolicy(pModel *model.PModel) (bool, error)
	Delete(pModel *model.PModel) (bool, error)
	Update(ur *model.UpdateRequest) (bool, error)

	AddGroupingPolicy(role string, group string) bool
	ListMyPolicy(uname string) *model.BackResp
	ListRoles(role string) []*model.GModel
}

func (c *Casbin) UnBindingRoles(gModel *model.GModel) (bool, error) {
	if gModel.PType == "g" {
		addExist := c.Enforcer.HasGroupingPolicy(gModel.User, gModel.Role)
		if addExist {
			policy := c.Enforcer.GetFilteredGroupingPolicy(1, gModel.Role)
			if len(policy) <= 1 {
				filteredPolicy, err := c.Enforcer.RemoveFilteredPolicy(0, gModel.Role)
				if err != nil {
					return false, err
				}
				if !filteredPolicy {
					return false, errors.New("内部错误")
				}
			}
			result, _ := c.Enforcer.RemoveGroupingPolicy(gModel.User, gModel.Role)
			if result {
				err := c.Enforcer.SavePolicy()
				if err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, errors.New("权限不存在")
}

func (c *Casbin) BindingRoles(gModel *model.GModel) (bool, error) {
	if gModel.PType == "g" {
		policy := c.Enforcer.GetFilteredPolicy(0, gModel.Role)
		if len(policy) == 0 {
			return false, errors.New("角色未进行权限初始化,请先初始化角色权限！")
		}
		addExist := c.Enforcer.HasGroupingPolicy(gModel.User, gModel.Role)
		if !addExist {
			result, _ := c.Enforcer.AddGroupingPolicy(gModel.User, gModel.Role)
			if result {
				err := c.Enforcer.SavePolicy()
				if err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, errors.New("内部错误！")
}

func (c *Casbin) AddPolicy(pModel *model.PModel) (bool, error) {
	if pModel.PType == "p" {
		addExist := c.Enforcer.HasPolicy(pModel.UserRole, pModel.Namespace, pModel.Source, pModel.Action)
		if !addExist {
			result, err := c.Enforcer.AddPolicy(pModel.UserRole, pModel.Namespace, pModel.Source, pModel.Action)
			if err != nil {
				return false, err
			}
			if result {
				err = c.Enforcer.SavePolicy()
				if err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, errors.New("内部错误！")
}
func (c *Casbin) AddGroupingPolicy(role string, group string) bool {
	s, _ := c.Enforcer.AddRoleForUser(role, group)
	if s {
		logger.Infow("InitPermission", role+":"+group, "权限初始化成功")
		return s
	}
	return false
}
func (c *Casbin) Update(ur *model.UpdateRequest) (bool, error) {

	if ur.OldPolicy.PType == "p" && ur.NewPolicy.PType == "p" {
		updateExist := c.Enforcer.HasPolicy(ur.NewPolicy.UserRole, ur.NewPolicy.Namespace, ur.NewPolicy.Source, ur.NewPolicy.Action)
		if updateExist {
			return false, errors.New("已存在！")
		}
		oldExist := c.Enforcer.HasPolicy(ur.OldPolicy.UserRole, ur.OldPolicy.Namespace, ur.OldPolicy.Source, ur.OldPolicy.Action)
		if !oldExist {
			return false, errors.New("权限不存在！")
		}
		//old,new
		policy, err := c.Enforcer.UpdatePolicy([]string{ur.OldPolicy.UserRole, ur.OldPolicy.Namespace, ur.OldPolicy.Source, ur.OldPolicy.Action}, []string{ur.NewPolicy.UserRole, ur.NewPolicy.Namespace, ur.NewPolicy.Source, ur.NewPolicy.Action})
		if err != nil {
			return policy, err
		}
		if policy {
			err = c.Enforcer.SavePolicy()
			if err != nil {
				return false, err
			}
		}
		return policy, err

	}
	return false, errors.New("实体有误！")
}

func (c *Casbin) ListMyPolicy(uname string) (res *model.BackResp) {
	res = &model.BackResp{}
	group := c.Enforcer.GetFilteredGroupingPolicy(0, uname)
	var roles []string
	var policies [][]string
	if len(group) != 0 {
		res.GPolicy = filterGModel(group)
		roles = filterRoles(group)
		roles = append(roles, uname)
	}
	if len(roles) != 0 {
		for i := range roles {
			policies = append(policies, c.Enforcer.GetFilteredPolicy(0, roles[i])...)
		}
	}
	if len(policies) != 0 {
		res.PPolicy = filterPModel(policies)
	}
	return res
}
func (c *Casbin) Delete(pModel *model.PModel) (bool, error) {
	if pModel.PType == "p" {
		hasPolicy := c.Enforcer.HasPolicy(pModel.UserRole, pModel.Namespace, pModel.Source, pModel.Action)
		if !hasPolicy {
			return false, errors.New("权限不存在！")
		}
		filteredPolicy := c.Enforcer.GetFilteredPolicy(0, pModel.UserRole)
		// 若删除的是角色,查询角色下是否有绑定用户
		if len(filteredPolicy) <= 1 {
			policy := c.Enforcer.GetFilteredGroupingPolicy(1, pModel.UserRole)
			if len(policy) != 0 {
				return false, errors.New("该角色下存在绑定用户！")
			}
		}
		result, _ := c.Enforcer.RemovePolicy(pModel.UserRole, pModel.Namespace, pModel.Source, pModel.Action)
		if result {
			err := c.Enforcer.SavePolicy()
			if err != nil {
				return false, errors.New("内部错误！")
			}
		}
		return result, nil
	}
	return false, errors.New("实体错误！")
}

func (c *Casbin) ListRoles(role string) (gModel []*model.GModel) {
	policy := c.Enforcer.GetFilteredGroupingPolicy(1, role)
	if len(policy) != 0 {
		gModel = filterGModel(policy)
		return
	}
	return
}
