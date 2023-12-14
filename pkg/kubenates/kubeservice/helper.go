package kubeservice

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type helper struct{}

const pattern = "[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\\.?"

func (h *helper) GetImages(deployment v1.Deployment) string {
	return h.GetImagesByPod(deployment.Spec.Template.Spec.Containers)
}

func (h *helper) GetImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	if imagesLen := len(containers); imagesLen > 1 {
		images += fmt.Sprintf("其他%d个镜像", imagesLen-1)
	}
	return images
}

func (h *helper) PodIsReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != "Running" {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != "True" {
			return false
		}
	}
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return true
}

func (h *helper) LabelsFilter(labels map[string]string) (ls []string) {
	for k, v := range labels {
		ls = append(ls, fmt.Sprintf("%s=%s", k, v))
	}
	return
}

func (h *helper) TaintsFilter(taints []corev1.Taint) (ret []string) {
	for _, taint := range taints {
		ret = append(ret, fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect))
	}
	return
}
