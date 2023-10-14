package casbin

import (
	"github.com/casbin/casbin"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/etcd"
	"github.com/oppslink/protocol/logger"
	"github.com/sebastianliu/etcd-adapter"
)

var (
	CasbinSetting *config.CasbinModelPath
)

// 权限
const (
	Read  = "read"
	Write = "write"
	Admin = "owner"
)

// 权限组
const (
	GroupRead    = "role_read"
	GroupWrite   = "role_write"
	GroupManager = "role_manager"
)

// 资源
const (
	// http 资源
	HttpV1      = "/v1"
	HttpManager = "/manager"

	//todo：k8s资源
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
	PType    string `json:"p_type" form:"p_type" description:"策略"`
	Role     string `json:"role" form:"v0" description:"角色/用户"`
	Source   string `json:"source" form:"v1" description:"资源"`
	Behavior string `json:"behavior" form:"v2" description:"行为"`
}

func InitCasbin(conf *config.Config) (*Casbin, error) {
	casbinAdapter := &CasbinAdapter{
		etcdEndpoint: conf.EtcdConfig.Endpoint,
		key:          etcd.CasbinRuleKey,
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
	roleRead := c.DefaultPolicy.e.HasPolicy(GroupRead, HttpV1, Read)
	if !roleRead {
		c.DefaultPolicy.e.AddPolicy(GroupRead, HttpV1, Read)
		logger.Infow("InitPermission", GroupRead, "权限初始化成功")
	}
	roleWrite := c.DefaultPolicy.e.HasPolicy(GroupWrite, HttpV1, Write)
	if !roleWrite {
		c.DefaultPolicy.e.AddPolicy(GroupWrite, HttpV1, Write)
		logger.Infow("InitPermission", GroupWrite, "权限初始化成功")
	}
	roleManager := c.DefaultPolicy.e.HasPolicy(GroupManager, HttpManager, Admin)
	if !roleManager {
		c.DefaultPolicy.e.AddPolicy(GroupManager, HttpManager, Admin)
		logger.Infow("InitPermission", GroupManager, "权限初始化成功")
	}

	// 角色初始化
	_ = c.DefaultPolicy.AddGroupingPolicy("admin", GroupManager)

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
	AddGroupingPolicy(role string, group string) bool
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
