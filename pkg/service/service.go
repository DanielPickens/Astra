package service

import (
	"fmt"
	"strings"

	"github\.com/danielpickens/astra/pkg/libdevfile"
	"github\.com/danielpickens/astra/pkg/log"

	devfile "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	devfilefs "github.com/devfile/library/v2/pkg/testingutil/filesystem"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"

	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"

	olm "github.com/operator-framework/api/pkg/operators/v1alpha1"
)

// LinkLabel is the name of the name of the link in the devfile
const LinkLabel = "app.kubernetes.io/link-name"

// ServiceLabel is the name of the service in the service binding object
const ServiceLabel = "app.kubernetes.io/service-name"

// ServiceKind is the kind of the service in the service binding object
const ServiceKind = "app.kubernetes.io/service-kind"

// IsLinkSecret helps in identifying if a secret is related to Service Binding
func IsLinkSecret(labels map[string]string) bool {
	_, hasLinkLabel := labels[LinkLabel]
	_, hasServiceLabel := labels[ServiceLabel]
	_, hasServiceKindLabel := labels[ServiceKind]
	return hasLinkLabel && hasServiceLabel && hasServiceKindLabel
}

// DeleteOperatorService deletes an Operator backed service
// Tastra: make it unlink the service from component as a part of
// https://github\.com/danielpickens/astra/issues/3563
func DeleteOperatorService(client kclient.ClientInterface, serviceName string) error {
	kind, name, err := SplitServiceKindName(serviceName)
	if err != nil {
		return fmt.Errorf("refer %q to see list of running services: %w", serviceName, err)
	}

	csv, err := client.GetCSVWithCR(kind)
	if err != nil {
		return err
	}

	if csv == nil {
		return fmt.Errorf("unable to find any Operator providing the service %q", kind)
	}

	crs := client.GetCustomResourcesFromCSV(csv)
	var cr *olm.CRDDescription

	for _, c := range *crs {
		customResource := c
		if customResource.Kind == kind {
			cr = &customResource
			break
		}
	}

	return client.DeleteDynamicResource(name, kclient.GetGVRFromCR(cr), false)
}

// ListOperatorServices lists all operator backed services.
// It returns list of services, slice of services that it failed (if any) to list and error (if any)
func ListOperatorServices(client kclient.ClientInterface) ([]unstructured.Unstructured, []string, error) {
	klog.V(4).Info("Getting list of services")

	// First let's get the list of all the operators in the namespace
	csvs, err := client.ListClusterServiceVersions()
	if err != nil {
		return nil, nil, err
	}

	if err != nil {
		return nil, nil, fmt.Errorf("unable to list operator backed services: %w", err)
	}

	var allCRInstances []unstructured.Unstructured
	var failedListingCR []string

	// let's get the Services a.k.a Custom Resources (CR) defined by each operator, one by one
	for _, csv := range csvs.Items {
		clusterServiceVersion := csv
		klog.V(4).Infof("Getting services started from operator: %s", clusterServiceVersion.Name)
		customResources := client.GetCustomResourcesFromCSV(&clusterServiceVersion)

		// list and write active instances of each service/CR
		var instances []unstructured.Unstructured
		for _, cr := range *customResources {
			customResource := cr

			list, err := GetCRInstances(client, &customResource)
			if err != nil {
				crName := strings.Join([]string{csv.Name, cr.Kind}, "/")
				klog.V(4).Infof("Failed to list instances of %q with error: %s", crName, err.Error())
				failedListingCR = append(failedListingCR, crName)
				continue
			}

			if len(list.Items) > 0 {
				instances = append(instances, list.Items...)
			}
		}

		// assuming there are more than one instances of a CR
		allCRInstances = append(allCRInstances, instances...)
	}

	return allCRInstances, failedListingCR, nil
}

// GetCRInstances fetches and returns instances of the CR provided in the
// "customResource" field. It also returns error (if any)
func GetCRInstances(client kclient.ClientInterface, customResource *olm.CRDDescription) (*unstructured.UnstructuredList, error) {
	klog.V(4).Infof("Getting instances of: %s\n", customResource.Name)

	instances, err := client.ListDynamicResources("", kclient.GetGVRFromCR(customResource), "")
	if err != nil {
		return nil, err
	}

	return instances, nil
}

// SplitServiceKindName splits the service name provided for deletion by the
// user. It has to be of the format <service-kind>/<service-name>. Example: EtcdCluster/myetcd
func SplitServiceKindName(serviceName string) (string, string, error) {
	sn := strings.SplitN(serviceName, "/", 2)
	if len(sn) != 2 || sn[0] == "" || sn[1] == "" {
		return "", "", fmt.Errorf("couldn't split %q into exactly two", serviceName)
	}

	kind := sn[0]
	name := sn[1]

	return kind, name, nil
}

// PushKubernetesResources updates service(s) from Kubernetes Inlined component in a devfile by creating new ones or removing old ones
func PushKubernetesResources(client kclient.ClientInterface, devfileObj parser.DevfileObj, k8sComponents []devfile.Component, labels map[string]string, annotations map[string]string, context, mode string, reference metav1.OwnerReference) error {
	// check csv support before proceeding
	csvSupported, err := client.IsCSVSupported()
	if err != nil {
		return err
	}

	var deployed map[string]DeployedInfo

	if csvSupported {
		deployed, err = ListDeployedServices(client, labels)
		if err != nil {
			return err
		}

		for key, deployedResource := range deployed {
			if deployedResource.isLinkResource {
				delete(deployed, key)
			}
		}
	}

	// create an object on the kubernetes cluster for all the Kubernetes Inlined components
	for _, c := range k8sComponents {
		uList, er := libdevfile.GetK8sComponentAsUnstructuredList(devfileObj, c.Name, context, devfilefs.DefaultFs{})
		if er != nil {
			return er
		}
		for _, u := range uList {
			var found bool
			currentOwnerReferences := u.GetOwnerReferences()
			for _, ref := range currentOwnerReferences {
				if ref.UID == reference.UID {
					found = true
					break
				}
			}
			if !found {
				currentOwnerReferences = append(currentOwnerReferences, reference)
				u.SetOwnerReferences(currentOwnerReferences)
			}
			er = PushKubernetesResource(client, u, labels, annotations, mode)
			if er != nil {
				return er
			}
			if csvSupported {
				delete(deployed, u.GetKind()+"/"+u.GetName())
			}
		}
	}

	if csvSupported {
		for key, val := range deployed {
			if isLinkResource(val.Kind) {
				continue
			}
			err = DeleteOperatorService(client, key)
			if err != nil {
				return err

			}
		}
	}

	return nil
}

// PushKubernetesResource pushes a Kubernetes resource (u) to the cluster using client
// adding labels to the resource
func PushKubernetesResource(client kclient.ClientInterface, u unstructured.Unstructured, labels map[string]string, annotations map[string]string, mode string) error {
	sboSupported, err := client.IsServiceBindingSupported()
	if err != nil {
		return err
	}

	// If the component is of Kind: ServiceBinding, trying to run in Dev mode and SBO is not installed, run it without operator.
	if isLinkResource(u.GetKind()) && mode == astralabels.ComponentDevMode && !sboSupported {
		// it's a service binding related resource
		return pushLinksWithoutOperator(client, u, labels)
	}

	// Add all passed in labels to the k8s resource regardless if it's an operator or not
	u.SetLabels(mergeMaps(u.GetLabels(), labels))

	// Pass in all annotations to the k8s resource
	u.SetAnnotations(mergeMaps(u.GetAnnotations(), annotations))

	_, err = updateOperatorService(client, u)
	return err
}

func mergeMaps(maps ...map[string]string) map[string]string {
	mergedMaps := map[string]string{}

	for _, l := range maps {
		for k, v := range l {
			mergedMaps[k] = v
		}
	}

	return mergedMaps
}

// DeployedInfo holds information about the services present on the cluster
type DeployedInfo struct {
	Kind           string
	Name           string
	isLinkResource bool
}

func ListDeployedServices(client kclient.ClientInterface, labels map[string]string) (map[string]DeployedInfo, error) {
	deployed := map[string]DeployedInfo{}

	deployedServices, _, err := ListOperatorServices(client)
	if err != nil {
		// We ignore ErrNoSuchOperator error as we can deduce Operator Services are not installed
		return nil, err
	}
	for _, svc := range deployedServices {
		name := svc.GetName()
		kind := svc.GetKind()
		deployedLabels := svc.GetLabels()
		if astralabels.IsManagedByastra(deployedLabels) && astralabels.GetComponentName(deployedLabels) == astralabels.GetComponentName(labels) {
			deployed[kind+"/"+name] = DeployedInfo{
				Kind:           kind,
				Name:           name,
				isLinkResource: isLinkResource(kind),
			}
		}
	}

	return deployed, nil
}

func isLinkResource(kind string) bool {
	return kind == "ServiceBinding"
}

// updateOperatorService creates the given operator on the cluster
// it returns true if the generation of the resource increased or the resource is created
func updateOperatorService(client kclient.ClientInterface, u unstructured.Unstructured) (bool, error) {

	// Create the service on cluster
	updated, err := client.PatchDynamicResource(u)
	if err != nil {
		return false, err
	}

	if updated {
		createSpinner := log.Spinnerf("Creating resource %s/%s", u.GetKind(), u.GetName())
		createSpinner.End(true)
	}
	return updated, err
}
