package kubeclient

import (
	"context"
	"fmt"
	config "github.com/denovo/permission/configration"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientInterface interface {
	Get()
	List(ctx context.Context, namespace string) *v1.PodList
	Apply()
}

type KubeClientType string

const (
	K8sClientTypeKubernetes KubeClientType = "kubernetes"
	K8sClientTypeDynamic    KubeClientType = "dynamic"
)

// NewClientInterface Kubernetes客户端的接口实例
func NewClientInterface(conf *config.OpsLinkConfig, clientType KubeClientType) (ClientInterface, error) {
	var err error
	kubeconfig := conf.Kubernetes.Kubeconfig
	configOverrides := &clientcmd.ConfigOverrides{}
	var kubecfg *rest.Config

	switch clientType {
	case K8sClientTypeKubernetes:
		kubecfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			configOverrides).ClientConfig()
		if err != nil {
			return nil, err
		}
		k8sClient, err := kubernetes.NewForConfig(kubecfg)
		if err != nil {
			return nil, err
		}
		return &KubernetesClient{kc: k8sClient}, nil

	case K8sClientTypeDynamic:
		kubecfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			configOverrides).ClientConfig()
		if err != nil {
			return nil, err
		}
		dynClient, err := dynamic.NewForConfig(kubecfg)
		if err != nil {
			return nil, err
		}
		return &DynamicClient{dc: dynClient}, nil

	default:
		return nil, fmt.Errorf("Unsupported client type: %s", clientType)
	}
}

// GetClientSet 断言转化接口实例为ClientSet，便于调用底层包中方法
func GetClientSet(kube ClientInterface) *KubernetesClient {
	if clientSet, ok := kube.(*KubernetesClient); ok {
		return clientSet
	}
	return nil
}

// GetDynamicClient 断言转化接口实例为DynamicClient，便于调用底层包中方法
func GetDynamicClient(kube ClientInterface) *DynamicClient {
	if dynClient, ok := kube.(*DynamicClient); ok {
		return dynClient
	}
	return nil
}
