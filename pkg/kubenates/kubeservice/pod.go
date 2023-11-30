package kubeservice

import (
	"github.com/denovo/permission/pkg/kubenates/informer"
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
