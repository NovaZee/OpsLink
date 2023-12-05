package casbin

import (
	"github.com/casbin/casbin"
	config "github.com/denovo/permission/config"
	"github.com/oppslink/protocol/logger"
	"github.com/sebastianliu/etcd-adapter"
	"os"
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
	// init permission  read,write,manager,admin
	newCasbin.InitPermission()

	return newCasbin, nil
}

func (c *Casbin) InitPermission() {

	// p, role_read, /v1, read
	// p, role_write, /v1, write
	// p, role_manager, /v1/manager, owner

	// roles init
	roleRead := c.Enforcer.HasPolicy(GroupRead, HttpV1, Read)
	if !roleRead {
		c.Enforcer.AddPolicy(GroupRead, HttpV1, Read)
		logger.Infow("InitPermission", GroupRead, "read policy init success!")
	}
	roleWrite := c.Enforcer.HasPolicy(GroupWrite, HttpV1, Write)
	if !roleWrite {
		c.Enforcer.AddPolicy(GroupWrite, HttpV1, Write)
		logger.Infow("InitPermission", GroupWrite, "write policy ini success!")
	}
	roleManager := c.Enforcer.HasPolicy(GroupManager, HttpManager, Admin)
	if !roleManager {
		c.Enforcer.AddPolicy(GroupManager, HttpManager, Admin)
		logger.Infow("InitPermission", GroupManager, "admin policy init success!")
	}

	// admin role init
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
	// Init etcd adapter
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
	// enable auto save todo：is valid?
	enforcer.EnableAutoSave(true)
	return enforcer, nil
}

func NewEnforcerProvider(conf *config.OpsLinkConfig) (EnforcerProvider, error) {
	if len(conf.EtcdConfig.Endpoint) != 0 {
		//The etcd endpoint is not null
		//load policy from etcd
		return &EtcdAdapterProvider{conf.EtcdConfig.Endpoint, config.CasbinRuleKey}, nil
	} else {
		//The etcd endpoint is null Init casbin policy csv
		//Loading policy from dist
		loadCsv()
		return &CsvAdapterProvider{config.CasbinCsvPath}, nil
	}
}

func loadCsv() {
	// check exit
	if _, err := os.Stat(config.CasbinCsvPath); os.IsNotExist(err) {
		// create empty file
		emptyFile, createErr := os.Create(config.CasbinCsvPath)
		if createErr != nil {
			logger.Errorw("Init CSV File Error!", createErr)
			return
		}
		defer emptyFile.Close()
		logger.Infow("Init CSV File Success!")
	} else {
		logger.Infow("CSV File Already exits!", "path", config.CasbinCsvPath)
	}
}
