package config

import (
	"errors"
	"github.com/oppslink/protocol/logger"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

var (
	ErrCfgFail = errors.New("config parse error")
)

type Config struct {
	EtcdConfig EtcdConfig      `yaml:"etcd"`
	Server     ServerConfig    `yaml:"server"`
	Logging    LoggingConfig   `yaml:"logging,omitempty"`
	CMPath     CasbinModelPath `yaml:"casbin_path,omitempty"`
}

// EtcdConfig Etcd配置
type EtcdConfig struct {
	Endpoint          []string      `yaml:"endpoint"`
	DialTimeout       int           `yaml:"dial_timeout"`
	DialKeepAliveTime time.Duration `yaml:"dial_keep_alive_time,omitempty"`
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
	EtcdConfig: EtcdConfig{
		Endpoint:    []string{"127.0.0.1:2379"},
		DialTimeout: 5,
	},
	Server: ServerConfig{
		RunMode:      "dev",
		HttpPort:     "8080",
		ReadTimeout:  60,
		WriteTimeout: 60,
	},
	CMPath: CasbinModelPath{
		ModelPath: "/media/denovo/data1/go/OpsLink/OpsLink/configration/cfg/rbac_model.conf",
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
