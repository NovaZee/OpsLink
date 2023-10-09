package model

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/denovo/permission/internal"
	"log"
	"strings"
)

type CasinoModel struct {
	PType  string `gorm:"column:p_type" json:"p_type" form:"p_type" description:"策略类型"`
	RoleId string `gorm:"column:v0" json:"role_id" form:"v0" description:"角色id"`
	Path   string `gorm:"column:v1" json:"path" form:"v1" description:"api路径"`
	Method string `gorm:"column:v2" json:"method" form:"v2" description:"方法"`
}

func (c *CasinoModel) TableName() string {
	return "casbin_rule"
}
func (c *CasinoModel) AddPolicy() error {
	return nil
}

func Casbin() *casbin.Enforcer {
	// 使用MySQL数据库初始化一个orm适配器
	adapter, err := gormadapter.NewAdapterByDB(internal.DBEngine)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	enforcer, err := casbin.NewEnforcer(internal.CasbinSetting, adapter)
	enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	enforcer.LoadPolicy()
	return enforcer
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
