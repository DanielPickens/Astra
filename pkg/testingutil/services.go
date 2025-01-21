package testingutil

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	astralabels "github\.com/danielpickens/astra/pkg/labels"
)

func FakeKubeService(componentName, serviceName string) corev1.Service {
	labels := astralabels.GetLabels(componentName, "app", "", astralabels.ComponentDevMode, false)
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: labels,
		},
	}
}

func FakeKubeServices(componentName string) []corev1.Service {
	return []corev1.Service{
		FakeKubeService(componentName, "service-1"),
		FakeKubeService(componentName, "service-2"),
		FakeKubeService(componentName, "service-3"),
	}
}
