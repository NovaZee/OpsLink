package config

import (
	"errors"
	"github.com/oppslink/protocol/logger"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	ErrCfgFail = errors.New("config parse error")
)

type Config struct {
	Dbs     DBConfig        `yaml:"dbs"`
	Server  ServerConfig    `yaml:"server"`
	Logging LoggingConfig   `yaml:"logging,omitempty"`
	CMPath  CasbinModelPath `yaml:"CMPath"`
}

// DBConfig 数据库配置
type DBConfig struct {
	DBType      string `yaml:"db_type"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	DBName      string `yaml:"db_name"`
	TablePrefix string `yaml:"table_prefix"`
	Charset     string `yaml:"charset,omitempty"`
	ParseTime   bool   `yaml:"parse_time,omitempty"`
	MaxIdleTime int    `yaml:"max_idle_time,omitempty"`
	MaxOpenConn int    `yaml:"max_open_conn,omitempty"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	RunMode      string `yaml:"run_mode,omitempty"`
	HttpPort     string `yaml:"http_port,omitempty"`
	ReadTimeout  int64  `yaml:"read_timeout,omitempty"`
	WriteTimeout int64  `yaml:"write_timeout,omitempty"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	logger.Config `yaml:",inline"`
}

var DefaultConfig = Config{
	Dbs: DBConfig{
		DBType:   "mysql",
		Username: "root",
		Password: "root",
		Host:     "127.0.01",
		DBName:   "OpsLink",
	},
	Server: ServerConfig{
		RunMode:      "dev",
		HttpPort:     "8080",
		ReadTimeout:  60,
		WriteTimeout: 60,
	},
	Logging: LoggingConfig{},
}

type CasbinModelPath struct {
	ModelPath string `yaml:"model_path,omitempty"`
}

func NewConfig(confString string, strictMode bool) (*Config, error) {
	cfg := DefaultConfig
	if confString != "" {
		decoder := yaml.NewDecoder(strings.NewReader(confString))
		decoder.KnownFields(strictMode)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func InitLoggerFromConfig(config LoggingConfig) {
	logger.InitLogConfig(config.Config, "OpsLink")
}
