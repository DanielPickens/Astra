package storage

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devfile/library/v2/pkg/devfile/generator"
	devfileParser "github.com/devfile/library/v2/pkg/devfile/parser"
	parsercommon "github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	dfutil "github.com/devfile/library/v2/pkg/util"

	"github\.com/danielpickens/astra/pkg/configAutomount"
	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/storage"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VolumeInfo is a struct to hold the pvc name and the volume name to create a volume.
// To be moved to devfile/library.
type VolumeInfo struct {
	PVCName    string
	VolumeName string
}

// GetVolumeInfos returns the PVC name attached to the `astra-projects` directory and a map of other PVCs
func GetVolumeInfos(pvcs []corev1.PersistentVolumeClaim) (astraSourcePVCName string, infos map[string]VolumeInfo, _ error) {
	infos = make(map[string]VolumeInfo)
	for _, pvc := range pvcs {
		// check if the pvc is in the terminating state
		if pvc.DeletionTimestamp != nil {
			continue
		}

		generatedVolumeName, e := generateVolumeNameFromPVC(pvc.Name)
		if e != nil {
			return "", nil, fmt.Errorf("unable to generate volume name from pvc name: %w", e)
		}

		storageName := astralabels.GetStorageName(pvc.Labels)
		if storageName == storage.astraSourceVolume {
			astraSourcePVCName = pvc.Name
			continue
		}

		infos[storageName] = VolumeInfo{
			PVCName:    pvc.Name,
			VolumeName: generatedVolumeName,
		}
	}
	return astraSourcePVCName, infos, nil
}

// GetPersistentVolumesAndVolumeMounts gets the PVC volumes and updates the containers with the volume mounts.
// volumeNameToVolInfo is a map of the devfile volume name to the volume info containing the pvc name and the volume name.
// To be moved to devfile/library.
func GetPersistentVolumesAndVolumeMounts(devfileObj devfileParser.DevfileObj, containers []corev1.Container, initContainers []corev1.Container, volumeNameToVolInfo map[string]VolumeInfo, options parsercommon.DevfileOptions) ([]corev1.Volume, error) {

	containerComponents, err := devfileObj.Data.GetDevfileContainerComponents(options)
	if err != nil {
		return nil, err
	}

	var pvcVols []corev1.Volume

	// We need to sort volumes to create Deployment in a deterministic way
	keys := make([]string, 0, len(volumeNameToVolInfo))
	for k := range volumeNameToVolInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, volName := range keys {
		volInfo := volumeNameToVolInfo[volName]
		pvcVols = append(pvcVols, getPVC(volInfo.VolumeName, volInfo.PVCName))

		// containerNameToMountPaths is a map of the Devfile container name to their Devfile Volume Mount Paths for a given Volume Name
		containerNameToMountPaths := make(map[string][]string)
		for _, containerComp := range containerComponents {
			for _, volumeMount := range containerComp.Container.VolumeMounts {
				if volName == volumeMount.Name {
					containerNameToMountPaths[containerComp.Name] = append(containerNameToMountPaths[containerComp.Name], generator.GetVolumeMountPath(volumeMount))
				}
			}
		}

		addVolumeMountToContainers(containers, initContainers, volInfo.VolumeName, containerNameToMountPaths)
	}
	return pvcVols, nil
}

func GetEphemeralVolumesAndVolumeMounts(devfileObj devfileParser.DevfileObj, containers []corev1.Container, initContainers []corev1.Container, ephemerals map[string]storage.Storage, options parsercommon.DevfileOptions) ([]corev1.Volume, error) {
	containerComponents, err := devfileObj.Data.GetDevfileContainerComponents(options)
	if err != nil {
		return nil, err
	}
	var emptydirVols []corev1.Volume

	// We need to sort volumes to create Deployment in a deterministic way
	keys := make([]string, 0, len(ephemerals))
	for k := range ephemerals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, volName := range keys {
		volInfo := ephemerals[volName]
		emptyDir, err := getEmptyDir(volInfo.Name, volInfo.Spec.Size)
		if err != nil {
			return nil, err
		}
		emptydirVols = append(emptydirVols, emptyDir)

		// containerNameToMountPaths is a map of the Devfile container name to their Devfile Volume Mount Paths for a given Volume Name
		containerNameToMountPaths := make(map[string][]string)
		for _, containerComp := range containerComponents {
			for _, volumeMount := range containerComp.Container.VolumeMounts {
				if volName == volumeMount.Name {
					containerNameToMountPaths[containerComp.Name] = append(containerNameToMountPaths[containerComp.Name], generator.GetVolumeMountPath(volumeMount))
				}
			}
		}

		addVolumeMountToContainers(containers, initContainers, volInfo.Name, containerNameToMountPaths)
	}
	return emptydirVols, nil
}

// getPVC gets a pvc type volume with the given volume name and pvc name.
func getPVC(volumeName, pvcName string) corev1.Volume {

	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
			},
		},
	}
}

// getEmptyDir gets an emptyDir type volume with the given volume name and size.
// size should be parseable as a Kubernetes `Quantity` or an error will be returned
func getEmptyDir(volumeName string, size string) (corev1.Volume, error) {

	emptyDir := &corev1.EmptyDirVolumeSource{}
	qty, err := resource.ParseQuantity(size)
	if err != nil {
		return corev1.Volume{}, err
	}
	emptyDir.SizeLimit = &qty
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: emptyDir,
		},
	}, nil
}

// addVolumeMountToContainers adds the Volume Mounts in containerNameToMountPaths to the containers for a given pvc and volumeName
// containerNameToMountPaths is a map of a container name to an array of its Mount Paths.
// To be moved to devfile/library.
func addVolumeMountToContainers(containers []corev1.Container, initContainers []corev1.Container, volumeName string, containerNameToMountPaths map[string][]string) {

	for containerName, mountPaths := range containerNameToMountPaths {
		for i := range containers {
			if containers[i].Name == containerName {
				for _, mountPath := range mountPaths {
					containers[i].VolumeMounts = append(containers[i].VolumeMounts, corev1.VolumeMount{
						Name:      volumeName,
						MountPath: mountPath,
						SubPath:   "",
					},
					)
				}
			}
		}
		for i := range initContainers {
			if strings.HasPrefix(initContainers[i].Name, containerName) {
				for _, mountPath := range mountPaths {
					initContainers[i].VolumeMounts = append(initContainers[i].VolumeMounts, corev1.VolumeMount{
						Name:      volumeName,
						MountPath: mountPath,
						SubPath:   "",
					},
					)
				}
			}
		}
	}
}

// generateVolumeNameFromPVC generates a volume name based on the pvc name
func generateVolumeNameFromPVC(pvc string) (volumeName string, err error) {
	volumeName, err = dfutil.NamespaceOpenShiftObject(pvc, "vol")
	if err != nil {
		return "", err
	}
	return
}

// HandleastraSourceStorage creates or deletes the volume containing project sources, based on the preference setting
// - if Ephemeral preference is true, any PVC with labels "component=..." and "astra-source-pvc=astra-projects" is removed
// - if Ephemeral preference is false and no PVC with matching labels exists, it is created
func HandleastraSourceStorage(client kclient.ClientInterface, storageClient storage.Client, componentName string, isEphemeral bool) error {
	selector := astralabels.Builder().WithComponentName(componentName).WithSourcePVC(storage.astraSourceVolume).Selector()
	pvcs, err := client.ListPVCs(selector)
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	if !isEphemeral {
		if len(pvcs) == 0 {
			err := storageClient.Create(storage.Storage{
				ObjectMeta: metav1.ObjectMeta{
					Name: storage.astraSourceVolume,
				},
				Spec: storage.StorageSpec{
					Size: storage.astraSourceVolumeSize,
				},
			})

			if err != nil {
				return err
			}
		} else if len(pvcs) > 1 {
			return fmt.Errorf("number of source volumes shouldn't be greater than 1")
		}
	} else {
		if len(pvcs) > 0 {
			for _, pvc := range pvcs {
				err := client.DeletePVC(pvc.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GetAutomountVolumes(configAutomountClient configAutomount.Client, containers, initContainers []corev1.Container) ([]corev1.Volume, error) {
	volumesInfos, err := configAutomountClient.GetAutomountingVolumes()
	if err != nil {
		return nil, err
	}

	var volumes []corev1.Volume
	for _, volumeInfo := range volumesInfos {
		switch volumeInfo.VolumeType {
		case configAutomount.VolumeTypePVC:
			volumes = mountPVC(volumeInfo, containers, initContainers, volumes)
		case configAutomount.VolumeTypeSecret:
			volumes = mountSecret(volumeInfo, containers, initContainers, volumes)
		case configAutomount.VolumeTypeConfigmap:
			volumes = mountConfigMap(volumeInfo, containers, initContainers, volumes)
		}
	}
	return volumes, nil
}

func mountPVC(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	volumeName := "auto-pvc-" + volumeInfo.VolumeName

	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		addVolumeMountToContainer(container, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: volumeInfo.MountPath,
			ReadOnly:  volumeInfo.ReadOnly,
		})
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: volumeInfo.VolumeName,
			},
		},
	})
	return volumes
}

func mountSecret(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	switch volumeInfo.MountAs {
	case configAutomount.MountAsFile:
		return mountSecretAsFile(volumeInfo, containers, initContainers, volumes)
	case configAutomount.MountAsEnv:
		return mountSecretAsEnv(volumeInfo, containers, initContainers, volumes)
	case configAutomount.MountAsSubpath:
		return mountSecretAsSubpath(volumeInfo, containers, initContainers, volumes)
	}
	return volumes
}

func mountSecretAsFile(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	volumeName := "auto-secret-" + volumeInfo.VolumeName

	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		addVolumeMountToContainer(container, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: volumeInfo.MountPath,
			ReadOnly:  volumeInfo.ReadOnly,
		})
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  volumeInfo.VolumeName,
				DefaultMode: volumeInfo.MountAccessMode,
			},
		},
	})
	return volumes
}

func mountSecretAsEnv(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		addEnvFromToContainer(container, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: volumeInfo.VolumeName,
				},
			},
		})
	})
	return volumes
}

func mountSecretAsSubpath(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	volumeName := "auto-secret-" + volumeInfo.VolumeName

	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		for _, key := range volumeInfo.Keys {
			addVolumeMountToContainer(container, corev1.VolumeMount{
				Name:      volumeName,
				MountPath: filepath.ToSlash(filepath.Join(volumeInfo.MountPath, key)),
				SubPath:   key,
				ReadOnly:  volumeInfo.ReadOnly,
			})
		}
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  volumeInfo.VolumeName,
				DefaultMode: volumeInfo.MountAccessMode,
			},
		},
	})
	return volumes
}

func mountConfigMap(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	switch volumeInfo.MountAs {
	case configAutomount.MountAsFile:
		return mountConfigMapAsFile(volumeInfo, containers, initContainers, volumes)
	case configAutomount.MountAsEnv:
		return mountConfigMapAsEnv(volumeInfo, containers, initContainers, volumes)
	case configAutomount.MountAsSubpath:
		return mountConfigMapAsSubpath(volumeInfo, containers, initContainers, volumes)
	}
	return volumes
}

func mountConfigMapAsFile(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	volumeName := "auto-cm-" + volumeInfo.VolumeName

	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		addVolumeMountToContainer(container, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: volumeInfo.MountPath,
			ReadOnly:  volumeInfo.ReadOnly,
		})
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				DefaultMode: volumeInfo.MountAccessMode,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: volumeInfo.VolumeName,
				},
			},
		},
	})
	return volumes
}

func mountConfigMapAsEnv(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		addEnvFromToContainer(container, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: volumeInfo.VolumeName,
				},
			},
		})
	})
	return volumes
}

func mountConfigMapAsSubpath(volumeInfo configAutomount.AutomountInfo, containers, initContainers []corev1.Container, volumes []corev1.Volume) []corev1.Volume {
	volumeName := "auto-cm-" + volumeInfo.VolumeName

	inAllContainers(containers, initContainers, func(container *corev1.Container) {
		for _, key := range volumeInfo.Keys {
			addVolumeMountToContainer(container, corev1.VolumeMount{
				Name:      volumeName,
				MountPath: filepath.ToSlash(filepath.Join(volumeInfo.MountPath, key)),
				SubPath:   key,
				ReadOnly:  volumeInfo.ReadOnly,
			})
		}
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: volumeInfo.VolumeName,
				},
				DefaultMode: volumeInfo.MountAccessMode,
			},
		},
	})
	return volumes
}

func inAllContainers(containers, initContainers []corev1.Container, f func(container *corev1.Container)) {
	for i := range containers {
		f(&containers[i])
	}
	for i := range initContainers {
		f(&initContainers[i])
	}
}

func addVolumeMountToContainer(container *corev1.Container, volumeMount corev1.VolumeMount) {
	container.VolumeMounts = append(container.VolumeMounts, volumeMount)
}

func addEnvFromToContainer(container *corev1.Container, envFrom corev1.EnvFromSource) {
	container.EnvFrom = append(container.EnvFrom, envFrom)
}
