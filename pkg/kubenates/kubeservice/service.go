package kubeservice

import (
	"context"
	"fmt"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	v3yaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/kubernetes"
)

type ServiceService struct {
	Si     *informer.ServiceInformer
	Client kubernetes.Interface
	helper *helper
}

func NewServiceService(client kubernetes.Interface) *ServiceService {
	return &ServiceService{Si: &informer.ServiceInformer{}, Client: client, helper: &helper{}}
}

func (s *ServiceService) ListServiceByNamespace(namespace string) (res []*kube.Service, err error) {
	services, err := s.Si.ListAll(namespace)
	if err != nil {
		return
	}
	for _, service := range services {
		res = append(res, &kube.Service{
			Name:       service.Name,
			Namespace:  service.Namespace,
			Type:       string(service.Spec.Type),
			ClusterIp:  service.Spec.ClusterIP,
			ClusterIps: service.Spec.ClusterIPs,
			Ports:      ports(service.Spec.Ports),
			Selector:   s.helper.ToArray(service.Spec.Selector),
		})
	}
	return
}

func (s *ServiceService) Get(ns, name string) (*corev1.Service, error) {
	services, err := s.Si.ListAll(ns)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if service.Name == name {
			//apiVersion: v1
			//kind: Service
			return service, nil
		}
	}
	return nil, nil
}

func (s *ServiceService) DownToYaml(ns, name string) ([]byte, error) {
	ss, err := s.Get(ns, name)
	if err != nil {
		return nil, err
	}
	sDate, converterErr := runtime.DefaultUnstructuredConverter.ToUnstructured(ss)
	if converterErr != nil {
		return nil, converterErr
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(sDate, &corev1.Pod{})
	if err != nil {
		return nil, err
	}
	sDate["apiVersion"] = "v1"
	sDate["kind"] = "Service"
	sByte, err := v3yaml.Marshal(sDate)
	if err != nil {
		return nil, err
	}
	return sByte, nil

	return nil, nil
}

func (s *ServiceService) ApplyByYaml(ctx context.Context, ns string, in []byte, isUpdate bool) error {
	// create unstructured object
	decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	service := &corev1.Service{}
	_, _, err := decode.Decode(in, nil, service)
	if err != nil {
		return err
	}
	if isUpdate {
		_, err = s.Client.CoreV1().Services(ns).Update(ctx, service, metav1.UpdateOptions{})
	} else {
		_, err = s.Client.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}
	return nil
}

func ports(servicePort []corev1.ServicePort) []string {
	res := make([]string, 0)
	for _, s := range servicePort {
		if s.NodePort == 0 {
			continue
		}
		res = append(res, fmt.Sprintf("%d/%s", s.NodePort, s.Protocol))
	}
	return res
}
