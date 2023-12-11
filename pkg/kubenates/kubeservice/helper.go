package kubeservice

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type helper struct{}

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
