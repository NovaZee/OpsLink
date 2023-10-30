package kubeclient

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
)

type DynamicClient struct {
	dc *dynamic.DynamicClient
}

func (d DynamicClient) List(ctx context.Context, namespace string) *v1.PodList {
	return nil
}

func (d DynamicClient) Get() {
	//TODO implement me
	panic("implement me")
}

func (d DynamicClient) Apply() {
	//TODO implement me
	panic("implement me")
}
