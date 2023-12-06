package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	Pi     *informer.PodInformer
	Client kubernetes.Interface

	helper *helper
}

func NewPodService(client kubernetes.Interface) *PodService {
	return &PodService{Pi: &informer.PodInformer{}, Client: client, helper: &helper{}}
}

func (ps *PodService) Get(ctx context.Context, ns, name string) (*v1.Pod, error) {
	get, err := ps.Client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}
