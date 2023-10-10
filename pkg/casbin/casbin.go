package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	config "github.com/denovo/permission/configration"
	"gorm.io/gorm"
	"log"
	"strings"
)

var (
	CasbinSetting *config.CasbinModelPath
)

type Casbin struct {
	cba           CasbinAdapter
	defaultpolicy DefaultPolicy
	rbacPolicy    RbacPolicy
}

type CasbinAdapter struct {
	engine *gorm.DB
	conf   string
}

type CasbinModel struct {
	PType  string `gorm:"column:p_type" json:"p_type" form:"p_type" description:"策略类型"`
	RoleId string `gorm:"column:v0" json:"role_id" form:"v0" description:"角色id"`
	Path   string `gorm:"column:v1" json:"path" form:"v1" description:"api路径"`
	Method string `gorm:"column:v2" json:"method" form:"v2" description:"方法"`
}

func NewCasbin(cba CasbinAdapter, rbac RbacPolicy, policy DefaultPolicy) *Casbin {
	return &Casbin{
		cba:           cba,
		defaultpolicy: policy,
		rbacPolicy:    rbac,
	}
}
func NewCasbinModel(s2 string, s3 string, s4 string) *CasbinModel {
	return &CasbinModel{
		RoleId: s2,
		Path:   s3,
		Method: s4,
	}
}

func NewCasbinAdapter(engine *gorm.DB, conf *config.Config) *CasbinAdapter {
	return &CasbinAdapter{
		engine: engine,
		conf:   conf.CMPath.ModelPath,
	}
}

// NewCasbin: usage for policy upate
func (c *CasbinAdapter) Casbin() (*casbin.Enforcer, error) {
	// 使用MySQL数据库初始化一个orm适配器
	adapter, err := gormadapter.NewAdapterByDB(c.engine)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	enforcer, err := casbin.NewEnforcer(c.conf, adapter)
	if err != nil {
		return nil, err
	}
	//enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = enforcer.LoadPolicy()
	return enforcer, nil
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	return util.KeyMatch2(key1, key2)
}

// 注册func到casbin
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	return ParamsMatch(name1, name2), nil
}

type Policy interface {
	Add(a any) bool
	Update() bool
	Delete() bool
}
