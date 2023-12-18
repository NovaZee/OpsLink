
# OpsLink

Role-based authentication Kubernetes container management platform


## Features

- Role-based subscription for Kubernetes (k8s) resource delivery
- Basic CRUD Operations for k8s Resources
- Identity Authentication and Authorization based on JWT and Casbin

## Running

Clone

```bash
  git clone https://github.com/NovaZee/OpsLink.git
```

Install

```bash
  go mod tidy
```

Run

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

