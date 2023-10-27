package kubeclient

import (
	"context"
	config "github.com/denovo/permission/configration"
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// clientset 调用结构体
type KubernetesClient struct {
	kc *kubernetes.Clientset
}

func NewKubernetesClient(conf *config.OpsLinkConfig) *KubernetesClient {
	clientInterface, _ := newClientInterface(conf, K8sClientTypeKubernetes)
	clientInterface
	return newClientInterface(conf, K8sClientTypeKubernetes)
}

func (kc *KubernetesClient) Get() {
	//TODO implement me
	panic("implement me")
}

func (kc *KubernetesClient) Apply() {
	//TODO implement me
	panic("implement me")
}

func (kc *KubernetesClient) List(ctx context.Context, namespace string) *v1.PodList {
	get, _ := kc.kc.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	for _, item := range get.Items {
		logger.Infow("kubeclient-system", "Namespace", item.Namespace, "Name", item.GetName())
	}
	return get
}
