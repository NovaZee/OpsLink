package casbin

import (
	"github.com/casbin/casbin"
	config "github.com/denovo/permission/config"
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
	Enforcer *casbin.Enforcer
}

func NewCasbin(Enforcer *casbin.Enforcer) *Casbin {
	return &Casbin{
		Enforcer: Enforcer,
	}
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

func InitCasbin(conf *config.OpsLinkConfig) (*Casbin, error) {
	provider, err := NewEnforcerProvider(conf)
	if err != nil {
		return nil, err
	}
	enforcer, err := provider.GetEnforcer(conf.CMPath.ModelPath)
	if err != nil {
		return nil, err
	}
	newCasbin := NewCasbin(enforcer)
	// 初始化权限  读，写，管理
	newCasbin.InitPermission()

	return newCasbin, nil
}

func (c *Casbin) InitPermission() {

	// p, role_read, /v1, read
	// p, role_write, /v1, write
	// p, role_manager, /v1/manager, owner

	// 用户初始化
	roleRead := c.Enforcer.HasPolicy(GroupRead, HttpV1, Read)
	if !roleRead {
		c.Enforcer.AddPolicy(GroupRead, HttpV1, Read)
		logger.Infow("InitPermission", GroupRead, "权限初始化成功")
	}
	roleWrite := c.Enforcer.HasPolicy(GroupWrite, HttpV1, Write)
	if !roleWrite {
		c.Enforcer.AddPolicy(GroupWrite, HttpV1, Write)
		logger.Infow("InitPermission", GroupWrite, "权限初始化成功")
	}
	roleManager := c.Enforcer.HasPolicy(GroupManager, HttpManager, Admin)
	if !roleManager {
		c.Enforcer.AddPolicy(GroupManager, HttpManager, Admin)
		logger.Infow("InitPermission", GroupManager, "权限初始化成功")
	}

	// 角色初始化
	_ = c.Enforcer.AddGroupingPolicy("admin", GroupManager)
	err := c.Enforcer.SavePolicy()
	if err != nil {
		return
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
	AddGroupingPolicy(role string, group string) bool
	Update(a any) bool
	Delete(a any) bool
}

func (c *Casbin) Add(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		result := c.Enforcer.AddPolicy(casbinModel.Role, casbinModel.Source, casbinModel.Behavior)
		if result {
			err := c.Enforcer.SavePolicy()
			if err != nil {
				return false
			}
		}
		return result
	}
	return false
}
func (c *Casbin) AddGroupingPolicy(role string, group string) bool {
	s := c.Enforcer.AddRoleForUser(role, group)
	if s {
		logger.Infow("InitPermission", role+":"+group, "权限初始化成功")
		return s
	}
	return false
}
func (c *Casbin) Update(a any) bool {
	if _, ok := a.([]*CasbinModel); ok {
		// 遍历集合中的每个 CasbinModel 并添加策略
		return true
	}
	return false
}
func (c *Casbin) Delete(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		result := c.Enforcer.RemovePolicy(casbinModel.Role, casbinModel.Source, casbinModel.Behavior)
		if result {
			err := c.Enforcer.SavePolicy()
			if err != nil {
				return false
			}
		}
		return result
	}
	return false
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

// EnforcerProvider 接口定义
type EnforcerProvider interface {
	GetEnforcer(modelConf string) (*casbin.Enforcer, error)
}

// EtcdAdapterProvider 结构体实现 EnforcerProvider 接口
type EtcdAdapterProvider struct {
	etcdEndpoint []string
	key          string
}

func (eap *EtcdAdapterProvider) GetEnforcer(modelConf string) (*casbin.Enforcer, error) {
	adapter := etcdadapter.NewAdapter(eap.etcdEndpoint, eap.key)
	enforcer := casbin.NewEnforcer(modelConf, adapter)
	_ = enforcer.LoadPolicy()
	enforcer.EnableAutoSave(true)
	return enforcer, nil
}

// CsvAdapterProvider 结构体实现 EnforcerProvider 接口
type CsvAdapterProvider struct {
	csvFilePath string
}

func (cap *CsvAdapterProvider) GetEnforcer(modelConf string) (*casbin.Enforcer, error) {
	enforcer := casbin.NewEnforcer(modelConf, cap.csvFilePath)
	_ = enforcer.LoadPolicy()
	// 启用自动保存选项。
	enforcer.EnableAutoSave(true)
	return enforcer, nil
}

func NewEnforcerProvider(conf *config.OpsLinkConfig) (EnforcerProvider, error) {
	if len(conf.EtcdConfig.Endpoint) != 0 {
		//etcd配置地址不为空 权限策略存入etcd中
		return &EtcdAdapterProvider{conf.EtcdConfig.Endpoint, config.CasbinRuleKey}, nil
	} else {
		//etcd为空，权限策略走磁盘
		return &CsvAdapterProvider{config.CasbinCsvPath}, nil
	}
}
