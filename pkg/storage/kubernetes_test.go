package storage

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"

	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/testingutil"
	"github\.com/danielpickens/astra/pkg/util"
)

func Test_kubernetesClient_List(t *testing.T) {
	type fields struct {
		generic generic
	}
	tests := []struct {
		name                string
		fields              fields
		returnedDeployments *appsv1.DeploymentList
		returnedPVCs        *corev1.PersistentVolumeClaimList
		want                StorageList
		wantErr             bool
	}{
		{
			name: "case 1: should error out for multiple pods returned",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeployment("nodejs", true),
					*testingutil.CreateFakeDeployment("nodejs", true),
				},
			},
			wantErr: true,
		},
		{
			name: "case 2: pod not found",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{},
			},
			want:    StorageList{},
			wantErr: false,
		},
		{
			name: "case 3: no volume mounts on pod",
			fields: fields{
				generic: generic{
					componentName: "nodejs",
					appName:       "app",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeployment("nodejs", true),
				},
			},
			want:    StorageList{},
			wantErr: false,
		},
		{
			name: "case 4: two volumes mounted on a single container",
			fields: fields{
				generic: generic{
					componentName: "nodejs",
					appName:       "app",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-vol", MountPath: "/data"},
							{Name: "volume-1-vol", MountPath: "/path"},
						}),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
					*testingutil.FakePVC("volume-1", "10Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-1").Labels()),
				},
			},
			want: StorageList{
				Items: []Storage{
					generateStorage(NewStorage("volume-0", "5Gi", "/data", nil), "", "container-0"),
					generateStorage(NewStorage("volume-1", "10Gi", "/path", nil), "", "container-0"),
				},
			},
			wantErr: false,
		},
		{
			name: "case 5: one volume is mounted on a single container and another on both",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-vol", MountPath: "/data"},
							{Name: "volume-1-vol", MountPath: "/path"},
						}),
						testingutil.CreateFakeContainerWithVolumeMounts("container-1", []corev1.VolumeMount{
							{Name: "volume-1-vol", MountPath: "/path"},
						}),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
					*testingutil.FakePVC("volume-1", "10Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-1").Labels()),
				},
			},
			want: StorageList{
				Items: []Storage{
					generateStorage(NewStorage("volume-0", "5Gi", "/data", nil), "", "container-0"),
					generateStorage(NewStorage("volume-1", "10Gi", "/path", nil), "", "container-0"),
					generateStorage(NewStorage("volume-1", "10Gi", "/path", nil), "", "container-1"),
				},
			},
			wantErr: false,
		},
		{
			name: "case 6: pvc for volumeMount not found",
			fields: fields{
				generic: generic{
					componentName: "nodejs",
					appName:       "app",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-vol", MountPath: "/data"},
						}),
						testingutil.CreateFakeContainer("container-1"),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
					*testingutil.FakePVC("volume-1", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-1").Labels()),
				},
			},
			wantErr: true,
		},
		{
			name: "case 7: the storage label should be used as the name of the storage",
			fields: fields{
				generic: generic{
					componentName: "nodejs",
					appName:       "app",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-nodejs-vol", MountPath: "/data"},
						}),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0-nodejs", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
				},
			},
			want: StorageList{
				Items: []Storage{
					generateStorage(NewStorage("volume-0", "5Gi", "/data", nil), "", "container-0"),
				},
			},
			wantErr: false,
		},
		{
			name: "case 8: no pvc found for mount path",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-vol", MountPath: "/data"},
							{Name: "volume-1-vol", MountPath: "/path"},
						}),
						testingutil.CreateFakeContainerWithVolumeMounts("container-1", []corev1.VolumeMount{
							{Name: "volume-1-vol", MountPath: "/path"},
							{Name: "volume-vol", MountPath: "/path1"},
						}),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
					*testingutil.FakePVC("volume-1", "10Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-1").Labels()),
				},
			},
			want:    StorageList{},
			wantErr: true,
		},
		{
			name: "case 9: avoid the source volume's mount",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					*testingutil.CreateFakeDeploymentsWithContainers("nodejs", []corev1.Container{
						testingutil.CreateFakeContainerWithVolumeMounts("container-0", []corev1.VolumeMount{
							{Name: "volume-0-vol", MountPath: "/data"},
							{Name: "volume-1-vol", MountPath: "/path"},
						}),
						testingutil.CreateFakeContainerWithVolumeMounts("container-1", []corev1.VolumeMount{
							{Name: "volume-1-vol", MountPath: "/path"},
							{Name: astraSourceVolume, MountPath: "/path1"},
						}),
					}, []corev1.Container{}, true),
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{
					*testingutil.FakePVC("volume-0", "5Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-0").Labels()),
					*testingutil.FakePVC("volume-1", "10Gi", astralabels.Builder().WithComponent("nodejs").WithDevfileStorageName("volume-1").Labels()),
				},
			},
			want: StorageList{
				Items: []Storage{
					generateStorage(NewStorage("volume-0", "5Gi", "/data", nil), "", "container-0"),
					generateStorage(NewStorage("volume-1", "10Gi", "/path", nil), "", "container-0"),
					generateStorage(NewStorage("volume-1", "10Gi", "/path", nil), "", "container-1"),
				},
			},
			wantErr: false,
		},
		{
			name: "case 10: avoid the mandatory volume mounts used by astra",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			returnedDeployments: &appsv1.DeploymentList{
				Items: []appsv1.Deployment{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "nodejs",
							Labels: map[string]string{
								"component": "nodejs",
							},
						},
						Spec: appsv1.DeploymentSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									InitContainers: []corev1.Container{
										{
											Name: "my-container-with-shared-project",
											VolumeMounts: []corev1.VolumeMount{
												{
													Name:      "astra-shared-project",
													MountPath: "/opt/",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			returnedPVCs: &corev1.PersistentVolumeClaimList{
				Items: []corev1.PersistentVolumeClaim{},
			},
			want:    StorageList{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient, fakeClientSet := kclient.FakeNew()

			fakeClientSet.Kubernetes.PrependReactor("list", "persistentvolumeclaims", func(action ktesting.Action) (bool, runtime.Object, error) {
				return true, tt.returnedPVCs, nil
			})

			fakeClientSet.Kubernetes.PrependReactor("list", "deployments", func(action ktesting.Action) (bool, runtime.Object, error) {
				return true, tt.returnedDeployments, nil
			})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			k := kubernetesClient{
				generic: tt.fields.generic,
				client:  fakeClient,
			}
			got, err := k.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("kubernetesClient.List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_kubernetesClient_Create(t *testing.T) {
	type fields struct {
		generic generic
	}
	type args struct {
		storage Storage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case 1: valid storage",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			args: args{
				storage: NewStorageWithContainer("storage-0", "5Gi", "/data", "runtime", util.GetBool(false)),
			},
		},
		{
			name: "case 2: invalid storage size",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			args: args{
				storage: NewStorageWithContainer("storage-0", "example", "/data", "runtime", util.GetBool(false)),
			},
			wantErr: true,
		},
		{
			name: "case 3: valid astra project related storage",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			args: args{
				storage: NewStorageWithContainer("astra-projects-vol", "5Gi", "/data", "runtime", util.GetBool(false)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fkclient, fkclientset := kclient.FakeNew()

			k := kubernetesClient{
				generic: tt.fields.generic,
				client:  fkclient,
			}
			if err := k.Create(tt.args.storage); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == true {
				return
			}

			// Check for validating actions performed
			if (len(fkclientset.Kubernetes.Actions()) != 1) && (tt.wantErr != true) {
				t.Errorf("expected 1 action in CreatePVC got: %v", fkclientset.Kubernetes.Actions())
				return
			}

			createdPVC := fkclientset.Kubernetes.Actions()[0].(ktesting.CreateAction).GetObject().(*corev1.PersistentVolumeClaim)
			quantity, err := resource.ParseQuantity(tt.args.storage.Spec.Size)
			if err != nil {
				t.Errorf("failed to create quantity by calling resource.ParseQuantity(%v)", tt.args.storage.Spec.Size)
			}

			wantLabels := astralabels.GetLabels(tt.fields.generic.componentName, tt.fields.generic.appName, "", astralabels.ComponentDevMode, false)
			astralabels.AddStorageInfo(wantLabels, tt.args.storage.Name, strings.Contains(tt.args.storage.Name, astraSourceVolume))

			// created PVC should be labeled with labels passed to CreatePVC
			if diff := cmp.Diff(wantLabels, createdPVC.Labels); diff != "" {
				t.Errorf("kubernetesClient.Create() wantLabels mismatch (-want +got):\n%s", diff)
			}
			// name, size of createdPVC should be matching to size, name passed to CreatePVC
			if diff := cmp.Diff(quantity, createdPVC.Spec.Resources.Requests["storage"]); diff != "" {
				t.Errorf("kubernetesClient.Create() quantity mismatch (-want +got):\n%s", diff)
			}

			wantedPVCName, err := generatePVCName(tt.args.storage.Name, tt.fields.generic.componentName, tt.fields.generic.appName)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(wantedPVCName, createdPVC.Name); diff != "" {
				t.Errorf("kubernetesClient.Create() wantedPVCName mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_kubernetesClient_Delete(t *testing.T) {
	pvcName := "pvc-0"
	returnedPVCs := corev1.PersistentVolumeClaimList{
		Items: []corev1.PersistentVolumeClaim{
			*testingutil.FakePVC(pvcName, "5Gi", getStorageLabels("storage-0", "nodejs", "app")),
		},
	}

	type fields struct {
		generic generic
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case 1: delete successful",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			args: args{
				"storage-0",
			},
			wantErr: false,
		},
		{
			name: "case 2: pvc not found",
			fields: fields{
				generic: generic{
					appName:       "app",
					componentName: "nodejs",
				},
			},
			args: args{
				"storage-example",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fkclient, fkclientset := kclient.FakeNew()

			fkclientset.Kubernetes.PrependReactor("list", "persistentvolumeclaims", func(action ktesting.Action) (bool, runtime.Object, error) {
				return true, &returnedPVCs, nil
			})

			fkclientset.Kubernetes.PrependReactor("delete", "persistentvolumeclaims", func(action ktesting.Action) (bool, runtime.Object, error) {
				if action.(ktesting.DeleteAction).GetName() != pvcName {
					t.Errorf("delete called with = %v, want %v", action.(ktesting.DeleteAction).GetName(), pvcName)
				}
				return true, nil, nil
			})

			k := kubernetesClient{
				generic: tt.fields.generic,
				client:  fkclient,
			}
			if err := k.Delete(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(fkclientset.Kubernetes.Actions()) != 2 {
				t.Errorf("expected 2 action, got %v", len(fkclientset.Kubernetes.Actions()))
			}
		})
	}
}
