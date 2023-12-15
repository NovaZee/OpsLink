package kubeservice

import "context"

type IHub interface {
	ApplyByYaml(ctx context.Context, ns string, yaml []byte) error
}
