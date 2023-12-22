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

// domain 域
const (
	KubeSystem   = "kube-system"
	PolicyModule = "policy"
	Default      = "default"
	All          = "*"
)

// 角色身份
// owner 所有权  write
// user 用户   read
// super_admin 超级管理员
const (
	Admin       = "admin"
	User        = "user"
	SuperAdmin  = "superAdmin"
	InitAdmin   = "initAdmin"
	PolicyAdmin = "policyAdmin"
)

// 内置用户
// 初始化admin：opslink 初始化默认资源的超级权限(kube-system&default)-----initAdmin
// 超级root：admin k8s所有资源的超级权限-----superAdmin
const (
	AdminRoot = "admin"
	OpsLink   = "opslink"
)

// 资源,行为
const ()

// 权限 C->B->A  A包括B,C  B包括C
const (
	Read  = "read"
	Write = "write"
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
	Domain   string `json:"domain" form:"v1" description:"域/角色"`
	Source   string `json:"source" form:"v2" description:"资源"`
	Behavior string `json:"behavior" form:"v3" description:"行为"`
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
	//p, initAdmin, default, *, write
	//p, initAdmin, kube-system, *, write

	//g, admin, superAdmin
	//g, opslink, initAdmin

	//g, admin, superAdmin
	adminExists := c.Enforcer.HasGroupingPolicy(AdminRoot, SuperAdmin)
	if !adminExists {
		c.Enforcer.AddGroupingPolicy(AdminRoot, SuperAdmin)
		logger.Infow("InitPermission", AdminRoot, "admin role init success!")
	}
	//g, opslink, policyAdmin
	policy := c.Enforcer.HasGroupingPolicy(OpsLink, PolicyAdmin)
	if !policy {
		c.Enforcer.AddGroupingPolicy(OpsLink, PolicyAdmin)
	}
	//g, opslink, initAdmin
	opsKube := c.Enforcer.HasGroupingPolicy(OpsLink, InitAdmin)
	if !opsKube {
		c.Enforcer.AddGroupingPolicy(OpsLink, InitAdmin)
	}
	//p, policyAdmin, policy, *, write
	policyModule := c.Enforcer.HasPolicy(PolicyAdmin, PolicyModule, All, Write)
	if !policyModule {
		c.Enforcer.AddPolicy(PolicyAdmin, PolicyModule, All, Write)
	}
	//p, superAdmin, *, *, write
	superAdmin := c.Enforcer.HasPolicy(SuperAdmin, All, All, Write)
	if !superAdmin {
		c.Enforcer.AddPolicy(SuperAdmin, All, All, Write)
	}
	//p, initAdmin, default, *, write
	kubeRead := c.Enforcer.HasPolicy(InitAdmin, Default, All, Write)
	if !kubeRead {
		c.Enforcer.AddPolicy(InitAdmin, Default, All, Write)
	}
	//p, initAdmin, kube-system, *, write
	deafultWrite := c.Enforcer.HasPolicy(InitAdmin, KubeSystem, All, Write)
	if !deafultWrite {
		c.Enforcer.AddPolicy(InitAdmin, KubeSystem, All, Write)
	}

	err := c.Enforcer.SavePolicy()
	if err != nil {
		return
	}

}
func NewCasbinModel(s0, s1, s2, s3, s4 string) *CasbinModel {
	return &CasbinModel{
		PType:    s0,
		Role:     s1,
		Domain:   s2,
		Source:   s3,
		Behavior: s4,
	}
}

// Casbin Casbin: usage for policy upate
func (c *CasbinAdapter) Casbin() (*casbin.Enforcer, error) {
	// Init etcd adapter
	adapter := etcdadapter.NewAdapter(c.etcdEndpoint, c.key)
	enforcer := casbin.NewEnforcer(c.modelConf, adapter)
	enforcer.AddFunction("isSuper", isSuper)
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
		//ABAC
		if casbinModel.PType == "p" {
			addExist := c.Enforcer.HasPolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
			if !addExist {
				result := c.Enforcer.AddPolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
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
				result := c.Enforcer.AddGroupingPolicy(casbinModel.Role, casbinModel.Domain)
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
	s := c.Enforcer.AddRoleForUser(role, group)
	if s {
		logger.Infow("InitPermission", role+":"+group, "权限初始化成功")
		return s
	}
	return false
}
func (c *Casbin) Update(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		//todo:判断当前身份是否已存在，若从只读更新为可写，新增write，删除只读。因为可写包括只读，冗余了。
		if casbinModel.PType == "p" {
			//过滤下标从0,从属性开始，
			policy := c.Enforcer.GetFilteredNamedPolicy("p", 0, casbinModel.Role)
			for _, strings := range policy {
				println(strings)
			}
		}
		if casbinModel.PType == "g" {
			policy := c.Enforcer.GetFilteredGroupingPolicy(0, casbinModel.Role)
			for _, strings := range policy {
				println(strings)
			}
		}
	}
	return false
}
func (c *Casbin) Delete(a any) bool {
	if casbinModel, ok := a.(*CasbinModel); ok {
		//ABAC
		if casbinModel.PType == "p" {
			result := c.Enforcer.RemovePolicy(casbinModel.Role, casbinModel.Domain, casbinModel.Source, casbinModel.Behavior)
			if result {
				err := c.Enforcer.SavePolicy()
				if err != nil {
					return false
				}
			}
			return result
		}
		//role
		if casbinModel.PType == "g" {
			result := c.Enforcer.RemoveGroupingPolicy(casbinModel.Role, casbinModel.Domain)
			if result {
				err := c.Enforcer.SavePolicy()
				if err != nil {
					return false
				}
				return true
			}
		}
		return false
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
	enforcer.AddFunction("isSuper", isSuper)
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
	enforcer.AddFunction("isSuper", isSuper)
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

func isSuperAdminMatch(userName string) bool {
	return userName == AdminRoot
}
func isSuper(args ...interface{}) (interface{}, error) {
	userName := args[0].(string)
	return (bool)(isSuperAdminMatch(userName)), nil
}
