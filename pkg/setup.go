package pkg

import (
	"fmt"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/oppslink/protocol/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB

var (
	CasbinSetting *config.CasbinModelPath
)

func InitializeServer(cfg *config.Config) (*OpsLinkServer, error) {
	err := initDBConnection(cfg)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewCasbin(DBEngine, cfg.CMPath.ModelPath)
	if err != nil {
		return nil, err
	}
	success, _ := e.Enforce("1", "2", "#")
	logger.Infow("success:", "bool", success)
	os, err := NewOpsLinkServer(cfg)
	return os, err
}

func initDBConnection(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		cfg.Dbs.Username,
		cfg.Dbs.Password,
		cfg.Dbs.Host,
		cfg.Dbs.DBName,
		cfg.Dbs.Charset,
		cfg.Dbs.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DBEngine = db
	if err != nil {
		return err
	}
	return nil
}
