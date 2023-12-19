package kubeservice

import (
	"github.com/denovo/permission/pkg/kubenates/informer"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

type RBACService struct {
	Ri     *informer.RBACInformer
	Sai    *informer.SAInformer
	Client kubernetes.Interface
}

func NewRBACService(client kubernetes.Interface) *RBACService {
	return &RBACService{Ri: &informer.RBACInformer{}, Sai: &informer.SAInformer{}, Client: client}
}

// ListRoles 获取roles列表
func (rs *RBACService) ListRoles(ns string) []*rbacv1.Role {
	list := rs.Ri.ListAllByNs(ns)
	return list
}

func (rs *RBACService) ListSa(ns string) []*corev1.ServiceAccount {
	return rs.Sai.ListTargetAll(ns)
}
