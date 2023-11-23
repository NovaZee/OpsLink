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

const CasbinCsvPath = "./config/file/casbin_policy.csv"
const LocalStorePath = "./config/file/role.bin"
const CasbinRuleKey = "casbin_policy"
const RoleKey = "role_key/"

type OpsLinkConfig struct {
	EtcdConfig EtcdConfig      `yaml:"etcd"`
	Server     ServerConfig    `yaml:"server"`
	Logging    LoggingConfig   `yaml:"logging,omitempty"`
	CMPath     CasbinModelPath `yaml:"casbin_path,omitempty"`
	Kubernetes Kubernetes      `yaml:"kubernetes,omitempty"`
}

// EtcdConfig Etcd配置
type EtcdConfig struct {
	Endpoint          []string `yaml:"endpoint"`
	DialTimeout       int      `yaml:"dial_timeout"`
	DialKeepAliveTime int      `yaml:"dial_keep_alive_time,omitempty"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	RunMode      string `yaml:"run_mode,omitempty"`
	HttpPort     string `yaml:"http_port,omitempty"`
	ReadTimeout  int64  `yaml:"read_timeout,omitempty"`
	WriteTimeout int64  `yaml:"write_timeout,omitempty"`
}

// Kubernetes a K8s specific struct to hold config
type Kubernetes struct {
	K8sAPIRoot string `json:"k8s_api_root"`
	Kubeconfig string `json:"kubeconfig"`
	NodeName   string `json:"node_name"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	logger.Config `yaml:",inline"`
}

var DefaultConfig = OpsLinkConfig{
	Server: ServerConfig{
		RunMode:      "dev",
		HttpPort:     "8080",
		ReadTimeout:  60,
		WriteTimeout: 60,
	},
	Kubernetes: Kubernetes{
		K8sAPIRoot: "",
		Kubeconfig: "",
		NodeName:   "",
	},
	CMPath: CasbinModelPath{
		ModelPath: "./config/file/rbac_model.conf",
	},
	Logging: LoggingConfig{},
}

type CasbinModelPath struct {
	ModelPath string `yaml:"model_path,omitempty"`
}

func NewConfig(confString string, strictMode bool) (*OpsLinkConfig, error) {
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
