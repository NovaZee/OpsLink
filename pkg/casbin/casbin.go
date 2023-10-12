package casbin

import (
	"github.com/casbin/casbin"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/componment"
	"github.com/oppslink/protocol/logger"
	"github.com/sebastianliu/etcd-adapter"
)

var (
	CasbinSetting *config.CasbinModelPath
)

const (
	Read  = "read"
	Write = "write"
	Admin = "admin"
)

type Casbin struct {
	Adapter       *CasbinAdapter
	DefaultPolicy *DefaultPolicy
	RbacPolicy    *RbacPolicy
}

type CasbinAdapter struct {
	etcdEndpoint []string
	key          string
	modelConf    string
}

type CasbinModel struct {
	PType    string `gorm:"column:p_type" json:"p_type" form:"p_type" description:"策略类型"`
	Role     string `gorm:"column:v0" json:"role" form:"v0" description:"角色id"`
	Source   string `gorm:"column:v1" json:"source" form:"v1" description:"资源"`
	Behavior string `gorm:"column:v2" json:"behavior" form:"v2" description:"权限"`
}

func InitCasbin(conf *config.Config) (*Casbin, error) {
	casbinAdapter := &CasbinAdapter{
		etcdEndpoint: []string{componment.HTTP + conf.EtcdConfig.Host + ":" + conf.EtcdConfig.Port},
		key:          componment.CasbinRuleKey,
		modelConf:    conf.CMPath.ModelPath,
	}
	defaultPolicy, _ := NewDefaultPolicy(casbinAdapter)
	rbacPolicy, _ := NewRbacPolicy(casbinAdapter)
	c := &Casbin{
		Adapter:       casbinAdapter,
		DefaultPolicy: defaultPolicy,
		RbacPolicy:    rbacPolicy,
	}

	// 初始化权限  读，写，管理
	c.InitPermission()

	return c, nil
}

func (c *Casbin) InitPermission() {

	// p, role_read, /v1, read
	// p, role_write, /v1, write
	// p, role_manager, /v1/manager, owner

	// 用户初始化
	roleRead := c.DefaultPolicy.e.HasPolicy("role_read", "/v1", "read")
	if !roleRead {
		c.DefaultPolicy.e.AddPolicy("role_read", "/v1", "read")
		logger.Infow("InitPermission", "role_read", "权限初始化成功")
	}
	roleWrite := c.DefaultPolicy.e.HasPolicy("role_write", "/v1", "write")
	if !roleWrite {
		c.DefaultPolicy.e.AddPolicy("role_write", "/v1", "write")
		logger.Infow("InitPermission", "role_write", "权限初始化成功")
	}
	roleManager := c.DefaultPolicy.e.HasPolicy("role_manager", "/v1/manager", "owner")
	if !roleManager {
		c.DefaultPolicy.e.AddPolicy("role_manager", "/v1/manager", "owner")
		logger.Infow("InitPermission", "role_manager", "权限初始化成功")
	}
	// 角色初始化
	user := c.RbacPolicy.e.AddRoleForUser("admin", "owner")
	if user {
		logger.Infow("InitPermission", "admin", "权限初始化成功")
	}

}
func NewCasbinModel(s2 string, s3 string, s4 string) *CasbinModel {
	return &CasbinModel{
		Role:     s2,
		Source:   s3,
		Behavior: s4,
	}
}

// Casbin Casbin: usage for policy upate
func (c *CasbinAdapter) Casbin() (*casbin.Enforcer, error) {
	// 初始化etcd适配器
	adapter := etcdadapter.NewAdapter(c.etcdEndpoint, c.key)
	enforcer := casbin.NewEnforcer(c.modelConf, adapter)
	_ = enforcer.LoadPolicy()
	return enforcer, nil
}

type Policy interface {
	Add(a any) bool
	Update(a any) bool
	Delete(a any) bool
}

//func ParamsMatch(fullNameKey1 string, key2 string) bool {
//	key1 := strings.Split(fullNameKey1, "?")[0]
//	return util.KeyMatch2(key1, key2)
//}
//
//// 注册func到casbin
//func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
//	name1 := args[0].(string)
//	name2 := args[1].(string)
//	return ParamsMatch(name1, name2), nil
//}
