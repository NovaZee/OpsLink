package componment

import (
	"fmt"
	config "github.com/denovo/permission/configration"
	"github.com/oppslink/protocol/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB

// initDBConnection 初始化db引擎
func InitDBConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		cfg.Dbs.Username,
		cfg.Dbs.Password,
		cfg.Dbs.Host,
		cfg.Dbs.DBName,
		cfg.Dbs.Charset,
		cfg.Dbs.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	logger.Infow(" init Db-connection", "InitDBConnection", cfg.Dbs.Host)
	return db, nil
}
