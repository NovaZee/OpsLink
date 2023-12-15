package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapService struct {
	Cmi    *informer.ConfigMapInformer
	Client kubernetes.Interface
	helper *helper
}

func NewConfigMapService(client kubernetes.Interface) *ConfigMapService {
	return &ConfigMapService{Cmi: &informer.ConfigMapInformer{}, Client: client, helper: &helper{}}
}
func (cms *ConfigMapService) Apply(ctx context.Context, cf *kube.ConfigMap) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cf.Name,
			Namespace: cf.Namespace,
		},
		Data: cms.helper.ToMap(cf.Data),
	}
	var err error
	if cf.Update {
		_, err = cms.Client.CoreV1().ConfigMaps(cf.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	} else {
		_, err = cms.Client.CoreV1().ConfigMaps(cf.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}
	return nil
}

func (cms *ConfigMapService) Delete(ctx context.Context, ns, name string) error {
	err := cms.Client.CoreV1().ConfigMaps(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (cms *ConfigMapService) GetConfigMap(ns, name string) (res *kube.ConfigMap, err error) {
	configmap, err := cms.Cmi.Get(ns, name)
	if err != nil {
		return
	}
	var maps []*kube.Map
	for k, v := range configmap.Data {
		newMap := &kube.Map{
			Key:   k,
			Value: v,
		}
		maps = append(maps, newMap)
	}
	res = &kube.ConfigMap{
		Name:       configmap.Name,
		Namespace:  configmap.Namespace,
		CreateTime: configmap.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Data:       maps,
	}
	return
}
func (cms *ConfigMapService) ListByNamespace(namespace string) (res []*kube.ConfigMap, err error) {
	cfgms, err := cms.Cmi.List(namespace)
	if err != nil {
		return
	}
	for _, configmap := range cfgms {
		res = append(res, &kube.ConfigMap{
			Name:       configmap.Name,
			Namespace:  configmap.Namespace,
			CreateTime: configmap.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Data:       cms.helper.ToArray(configmap.Data),
		})
	}
	return
}

func (cms *ConfigMapService) ApplyByYaml(ctx context.Context, ns string, in []byte) error {
	// create unstructured object
	decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	configMap := &corev1.ConfigMap{}
	_, _, err := decode.Decode(in, nil, configMap)
	if err != nil {
		return err
	}
	_, err = cms.Client.CoreV1().ConfigMaps(ns).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
