package internal

import (
	"fmt"
	config "github.com/denovo/permission/configration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB

var (
	CasbinSetting *config.CasbinModelPath
)

func InitDBConnection(cfg *config.Config) error {
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
