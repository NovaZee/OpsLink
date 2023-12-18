
# OpsLink

基于角色订阅的K8S管理平台

[中文](./config/file/README-zh-cn.md) | [English](README.md)
## 特性

- 基于角色的 Kubernetes （k8s） 资源订阅模式资源推送
- Kubernetes （k8s） 资源的基本CRUD操作
- 基于 JWT 和 Casbin 的身份认证和授权

## 运行

克隆

```bash
  git clone https://github.com/NovaZee/OpsLink.git
```

安装依赖

```bash
  go mod tidy
```

运行

```bash
  ./OpsLink --config /config.yaml
```
# Config

```yaml
# Server config
server:
  # Run mode ： dev/debug/product
  run_mode: debug
  # Http port
  http_port: 8082
  # Read time out /S
  read_timeout: 60
  # Write time out /S
  write_timeout: 60
# By default, no intermediate storage is used. Role identities are stored locally in binary files, and extensibility allows configuration with etcd.
etcd:
  endpoint:
    - "http://127.0.0.1:2379"
  dial_timeout: 5
# The default is to read the local kubeconfig file, but you can also manually specify or upload files that are not system-loaded (note permissions, as they might prevent communication with the API server).
kubernetes:
  kubeconfig: './test.yaml'
logging:
  level: debug
  json: false
# The program loads the default rbac.conf file by default.
casbin_path:
  model_path: './config/file/rbac_model.conf'
```

