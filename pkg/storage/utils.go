package storage

import (
	"fmt"

	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/util"
)

// getPVCNameFromStorageName returns the PVC associated with the given storage
func getPVCNameFromStorageName(client kclient.ClientInterface, storageName string) (string, error) {
	var selector = astralabels.Builder().WithStorageName(storageName).Selector()
	pvcs, err := client.ListPVCNames(selector)
	if err != nil {
		return "", fmt.Errorf("unable to get PVC names for selector %v: %w", selector, err)
	}
	numPVCs := len(pvcs)
	if numPVCs != 1 {
		return "", fmt.Errorf("expected exactly one PVC attached to storage %v, but got %v, %v", storageName, numPVCs, pvcs)
	}
	return pvcs[0], nil
}

// generatePVCName generates a PVC name from the Devfile volume name, component name and app name
func generatePVCName(volName, componentName, appName string) (string, error) {

	pvcName, err := util.NamespaceKubernetesObject(volName, componentName)
	if err != nil {
		return "", fmt.Errorf("unable to create namespaced name: %w", err)
	}

	pvcName, err = util.NamespaceKubernetesObject(pvcName, appName)
	if err != nil {
		return "", fmt.Errorf("unable to create namespaced name: %w", err)
	}

	return pvcName, nil
}

// ConvertListLocalToMachine converts storage config list to StorageList type
func ConvertListLocalToMachine(storageListConfig []LocalStorage) StorageList {

	var storageListLocal []Storage

	for _, storeLocal := range storageListConfig {
		s := NewStorage(storeLocal.Name, storeLocal.Size, storeLocal.Path, storeLocal.Ephemeral)
		s.Spec.ContainerName = storeLocal.Container
		storageListLocal = append(storageListLocal, s)
	}

	return NewStorageList(storageListLocal)
}
