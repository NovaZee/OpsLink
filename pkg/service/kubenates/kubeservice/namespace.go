package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/service/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	Nsi    *informer.NamespaceInformer
	Client kubernetes.Interface
}

func NewNamespaceService(client kubernetes.Interface) *NamespaceService {
	return &NamespaceService{Nsi: &informer.NamespaceInformer{}, Client: client}
}

func (nss *NamespaceService) List() (res []*kube.Namespace) {
	all := nss.Nsi.ListAll(0)
	for _, ns := range all {
		res = append(res, &kube.Namespace{
			Namespace:  ns.Name,
			Status:     string(ns.Status.Phase),
			CreateTime: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

func (nss *NamespaceService) Create(ctx context.Context, ns string) (*kube.Namespace, error) {
	// 新建一个 Namespace 对象
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
	n, err := nss.Client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return &kube.Namespace{Namespace: n.Name, Status: string(n.Status.Phase), CreateTime: n.CreationTimestamp.Format("2006-01-02 15:04:05")}, nil
}

func (nss *NamespaceService) Remove(ctx context.Context, ns string) error {
	err := nss.Client.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
