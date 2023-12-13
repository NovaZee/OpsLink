package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeService struct {
	Ni     *informer.NodeInformer
	Metric *versioned.Clientset
}

func NewNodeService(client *versioned.Clientset) *NodeService {
	return &NodeService{Ni: &informer.NodeInformer{}, Metric: client}
}

func (ns *NodeService) List(ctx context.Context) (res []*kube.Node) {
	all := ns.Ni.ListAll()
	for _, node := range all {
		usage := GetUsage(ns.Metric, node, ctx)
		res = append(res, &kube.Node{
			Name:     node.Name,
			Ip:       node.Status.Addresses[0].Address,
			HostName: node.Status.Addresses[1].Address,
			//Labels: helpers.FilterLabels(node.Labels),
			//Taints: helpers.FilterTaints(node.Spec.Taints),
			Capacity: &kube.NodeCapacity{
				Cpu:    node.Status.Capacity.Cpu().Value(),
				Memory: node.Status.Capacity.Memory().Value(),
				Pods:   node.Status.Capacity.Pods().Value(),
			},
			Usage: &kube.NodeUsage{
				Pods:   int32(ns.Ni.GetPodsNum(node.Name)),
				Cpu:    usage[0],
				Memory: usage[1],
			},
			CreateTime: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

func GetUsage(c *versioned.Clientset, node *v1.Node, ctx context.Context) []float64 {
	nodeMetric, _ := c.MetricsV1beta1().
		NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
	cpu := float64(nodeMetric.Usage.Cpu().MilliValue()) / float64(node.Status.Capacity.Cpu().MilliValue())
	memory := float64(nodeMetric.Usage.Memory().MilliValue()) / float64(node.Status.Capacity.Memory().MilliValue())
	return []float64{cpu, memory}
}
