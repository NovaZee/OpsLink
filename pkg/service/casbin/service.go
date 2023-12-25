package casbin

import (
	"errors"
	"github.com/casbin/casbin/v2"
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
	Add(a any) bool
	AddGroupingPolicy(role string, group string) bool
	Update(ur *UpdateRequest) (bool, error)
	Delete(a any) (bool, error)
	ListMyPolicy(uname string) ([]CasbinModel, error)
}

type CasbinModel struct {
	PType    string `json:"p_type" form:"p_type" description:"策略"`
	Role     string `json:"role" form:"v0" description:"角色/用户"`
	Domain   string `json:"domain" form:"v1" description:"域/角色"`
	Source   string `json:"source" form:"v2" description:"资源"`
	Behavior string `json:"behavior" form:"v3" description:"行为"`
}

// UpdateRequest 是请求的结构体，包含了旧实体和新实体的信息
type UpdateRequest struct {
	OldData *CasbinModel `json:"old_policy"`
	NewData *CasbinModel `json:"new_policy"`
}

func (c *Casbin) Add(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		//ABAC
		if casbinModel.PType == "p" {
			addExist := c.Enforcer.HasPolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
			if !addExist {
				result, _ := c.Enforcer.AddPolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
				if result {
					err := c.Enforcer.SavePolicy()
					if err != nil {
						return false
					}
				}
				return true
			}
		}
		//role
		if casbinModel.PType == "g" {
			addExist := c.Enforcer.HasGroupingPolicy(casbinModel.Role, casbinModel.Domain)
			if !addExist {
				result, _ := c.Enforcer.AddGroupingPolicy(casbinModel.Role, casbinModel.Domain)
				if result {
					err := c.Enforcer.SavePolicy()
					if err != nil {
						return false
					}
				}
				return true
			}
		}
		return false
	}
	return false
}
func (c *Casbin) AddGroupingPolicy(role string, group string) bool {
	s, _ := c.Enforcer.AddRoleForUser(role, group)
	if s {
		logger.Infow("InitPermission", role+":"+group, "权限初始化成功")
		return s
	}
	return false
}
func (c *Casbin) Update(ur *UpdateRequest) (bool, error) {

	if ur.OldData.PType == "p" && ur.NewData.PType == "p" {
		updateExist := c.Enforcer.HasPolicy(ur.NewData.Role, ur.NewData.Domain, ur.NewData.Source, ur.NewData.Behavior)
		if updateExist {
			return false, errors.New("已存在！")
		}
		//old,new
		policy, err := c.Enforcer.UpdatePolicy([]string{ur.OldData.Role, ur.OldData.Domain, ur.OldData.Source, ur.OldData.Behavior}, []string{ur.NewData.Role, ur.NewData.Domain, ur.NewData.Source, ur.NewData.Behavior})
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

func (c *Casbin) ListMyPolicy(uname string) (res []CasbinModel, err error) {
	policy := c.Enforcer.GetFilteredNamedPolicy("p", 0, uname)

	if len(policy) != 0 {
		res = c.ConvertToCasbinModel(policy, "p")
	}

	group := c.Enforcer.GetFilteredGroupingPolicy(0, uname)
	groupList := c.ConvertToCasbinModel(group, "g")
	res = append(res, groupList...)
	return res, nil
}
func (c *Casbin) Delete(a any) (bool, error) {
	if casbinModel, ok := a.(*CasbinModel); ok {
		//ABAC
		if casbinModel.PType == "p" {

			policy := c.Enforcer.GetFilteredGroupingPolicy(1, casbinModel.Role)
			if len(policy) != 0 {
				return false, errors.New("该角色下存在绑定用户！")
			}
			result, _ := c.Enforcer.RemovePolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
			if result {
				err := c.Enforcer.SavePolicy()
				if err != nil {
					return false, errors.New("内部错误！")
				}
			}
			return result, nil
		}
		//role
		if casbinModel.PType == "g" {
			result, err := c.Enforcer.RemoveGroupingPolicy(casbinModel.Role, casbinModel.Domain)
			if err != nil {
				logger.Warnw("RemoveGroupingPolicy", err, "")
				return false, err
			}
			if result {
				err = c.Enforcer.SavePolicy()
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
		return false, errors.New("类型不存在！")
	}
	return false, errors.New("实体错误！")
}

// 将 [][]string 转换为 []CasbinModel
func (c *Casbin) ConvertToCasbinModel(data [][]string, pType string) (res []CasbinModel) {
	if pType == "p" {
		for _, row := range data {
			if len(row) >= 4 {
				model := CasbinModel{
					Role:     row[0],
					Domain:   row[1],
					Source:   row[2],
					Behavior: row[3],
					PType:    pType,
				}
				res = append(res, model)
			}
		}
	}
	if pType == "g" {
		for _, row := range data {
			if len(row) >= 2 {
				model := CasbinModel{
					Role:   row[0],
					Domain: row[1],
					PType:  pType,
				}
				res = append(res, model)
				policy := c.Enforcer.GetFilteredNamedPolicy("p", 0, model.Domain)
				if len(policy) != 0 {
					res = append(res, c.ConvertToCasbinModel(policy, "p")...)
				}
			}
		}
	}
	return res
}
