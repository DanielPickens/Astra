package testingutil

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	astralabels "github\.com/danielpickens/astra/pkg/labels"
)

// CreateFakeDeployment creates a fake deployment with the given pod name and labels
// isPartOfComponent bool decides if the deployment is supposed to be a part of the core resources deployed by `astra dev`
func CreateFakeDeployment(podName string, isPartOfComponent bool) *appsv1.Deployment {
	fakeUID := types.UID("12345")
	labels := astralabels.Builder().
		WithApp("app").
		WithAppName("app").
		WithComponentName(podName).
		WithManager("astra").
		WithMode(astralabels.ComponentDevMode)
	if isPartOfComponent {
		labels = labels.WithComponent(podName)
	}
	deployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        podName,
			UID:         fakeUID,
			Labels:      labels.Labels(),
			Annotations: astralabels.Builder().WithProjectType(podName).Labels(),
		},
	}
	return &deployment
}

// CreateFakeDeploymentsWithContainers creates a fake pod with the given pod name, container name and containers
// isPartOfComponent bool decides if the deployment is supposed to be a part of the core resources deployed by `astra dev`
func CreateFakeDeploymentsWithContainers(podName string, containers []corev1.Container, initContainers []corev1.Container, isPartOfComponent bool) *appsv1.Deployment {
	fakeDeployment := CreateFakeDeployment(podName, isPartOfComponent)
	fakeDeployment.Spec.Template.Spec.Containers = containers
	fakeDeployment.Spec.Template.Spec.InitContainers = initContainers
	return fakeDeployment
}
