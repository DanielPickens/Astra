package testingutil

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	astralabels "github\.com/danielpickens/astra/pkg/labels"
)

// CreateFakePod creates a fake pod with the given pod name and component name
func CreateFakePod(componentName, podName, containerName string) *corev1.Pod {
	fakePod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: astralabels.GetLabels(componentName, "app", "", astralabels.ComponentDevMode, false),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: containerName,
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
	return fakePod
}
