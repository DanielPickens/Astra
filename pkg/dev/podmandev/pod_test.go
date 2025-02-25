package podmandev

import (
	"context"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/devfile/library/v2/pkg/devfile/parser/data"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/libdevfile/generator"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/podman"
	"github\.com/danielpickens/astra/pkg/version"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	devfileName = "mycmp"
	appName     = "app"

	command = generator.GetExecCommand(generator.ExecCommandParams{
		Id:          "run",
		Component:   "mycomponent",
		CommandLine: "./run",
		IsDefault:   pointer.Bool(true),
		Kind:        v1alpha2.RunCommandGroupKind,
	})

	baseComponent = generator.GetContainerComponent(generator.ContainerComponentParams{
		Name: "mycomponent",
		Container: v1alpha2.Container{
			Image:   "myimage",
			Args:    []string{"-f", "/dev/null"},
			Command: []string{"tail"},
		},
	})

	volume = generator.GetVolumeComponent(generator.VolumeComponentParams{
		Name: "myvolume",
	})
)

func buildBasePod(withPfContainer bool) *corev1.Pod {
	basePod := corev1.Pod{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: "mycmp-app",
			Labels: map[string]string{
				"app":                                  appName,
				"app.kubernetes.io/instance":           devfileName,
				"app.kubernetes.io/managed-by":         "astra",
				"app.kubernetes.io/managed-by-version": version.VERSION,
				"app.kubernetes.io/part-of":            appName,
				"component":                            devfileName,
				"astra.dev/mode":                         labels.ComponentDevMode,
				"astra.dev/project-type":                 "Not available",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Args:    []string{"-f", "/dev/null"},
					Command: []string{"tail"},
					Env: []corev1.EnvVar{
						{
							Name:  "PROJECTS_ROOT",
							Value: "/projects",
						},
						{
							Name:  "PROJECT_SOURCE",
							Value: "/projects",
						},
					},
					Image:           "myimage",
					ImagePullPolicy: "Always",
					Name:            "mycomponent",
					VolumeMounts: []corev1.VolumeMount{
						{
							MountPath: "/projects",
							Name:      "astra-projects",
						},
						{
							MountPath: "/opt/astra/",
							Name:      "astra-shared-data",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "astra-projects",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "astra-projects-mycmp-app",
						},
					},
				},
				{
					Name: "astra-shared-data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "astra-shared-data-mycmp-app",
						},
					},
				},
			},
		},
	}
	if withPfContainer {
		basePod.Spec.Containers = append(basePod.Spec.Containers, corev1.Container{
			Args:    []string{"-f", "/dev/null"},
			Command: []string{"tail"},
			Image:   portForwardingHelperImage,
			Name:    portForwardingHelperContainerName,
		})
	}
	return &basePod
}

func Test_createPodFromComponent(t *testing.T) {

	type args struct {
		devfileObj           func() parser.DevfileObj
		componentName        string
		appName              string
		debug                bool
		buildCommand         string
		runCommand           string
		debugCommand         string
		forwardLocalhost     bool
		customForwardedPorts []api.ForwardedPort
		customAddress        string
	}
	tests := []struct {
		name         string
		capabilities podman.Capabilities
		args         args
		wantPod      func(basePod *corev1.Pod) *corev1.Pod
		wantFwPorts  []api.ForwardedPort
		wantErr      bool
	}{
		{
			name: "basic component without command / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					_ = data.AddComponents([]v1alpha2.Component{baseComponent})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				return basePod.DeepCopy()
			},
		},
		{
			name: "basic component without command / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					_ = data.AddComponents([]v1alpha2.Component{baseComponent})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				return pod
			},
		},
		{
			name: "basic component with command / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Command = []string{"./cmd"}
					cmp.Container.Args = []string{"arg1", "arg2"}
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Command = []string{"./cmd"}
				pod.Spec.Containers[0].Args = []string{"arg1", "arg2"}
				return pod
			},
		},
		{
			name: "basic component with command / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Command = []string{"./cmd"}
					cmp.Container.Args = []string{"arg1", "arg2"}
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Command = []string{"./cmd"}
				pod.Spec.Containers[0].Args = []string{"arg1", "arg2"}
				return pod
			},
		},
		{
			name: "basic component + memory limit / forwardLocalhost=false, cgroup=v2",
			capabilities: podman.Capabilities{
				Cgroupv2: true,
			},
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.MemoryLimit = "1Gi"
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
					"memory": resource.MustParse("1Gi"),
				}
				return pod
			},
		},
		{
			name: "basic component + memory limit / forwardLocalhost=false, cgroup=v1",
			capabilities: podman.Capabilities{
				Cgroupv2: false,
			},
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.MemoryLimit = "1Gi"
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				return pod
			},
		},
		{
			name: "basic component + memory limit / forwardLocalhost=true, cgroup=v2",
			capabilities: podman.Capabilities{
				Cgroupv2: true,
			},
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.MemoryLimit = "1Gi"
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
					"memory": resource.MustParse("1Gi"),
				}
				return pod
			},
		},
		{
			name: "basic component + application endpoint / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 8080,
					HostPort:      20001,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20001,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - without debug / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 8080,
					HostPort:      20001,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - without debug / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20001,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - with debug / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 8080,
					HostPort:      20001,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "debug",
					ContainerPort: 5858,
					HostPort:      20002,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 5858,
					IsDebug:       true,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - with debug / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 8080,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				debug:            true,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20001,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "debug",
					ContainerPort: 20002,
					HostPort:      20002,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 5858,
					IsDebug:       true,
				},
			},
		},
		{
			name: "basic component with volume mount / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					_ = data.AddComponents([]v1alpha2.Component{baseComponent, volume})
					_ = data.AddVolumeMounts(baseComponent.Name, []v1alpha2.VolumeMount{
						{
							Name: volume.Name,
							Path: "/path/to/mount",
						},
					})

					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
					Name: volume.Name,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: volume.Name + "-" + devfileName + "-" + appName,
						},
					},
				})
				pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      volume.Name,
					MountPath: "/path/to/mount",
				})
				return pod
			},
		},
		{
			name: "basic component with volume mount / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					_ = data.AddComponents([]v1alpha2.Component{baseComponent, volume})
					_ = data.AddVolumeMounts(baseComponent.Name, []v1alpha2.VolumeMount{
						{
							Name: volume.Name,
							Path: "/path/to/mount",
						},
					})

					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
					Name: volume.Name,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: volume.Name + "-" + devfileName + "-" + appName,
						},
					},
				})
				pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      volume.Name,
					MountPath: "/path/to/mount",
				})
				return pod
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint + container ports known - with debug / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 20001,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 20002,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug-1",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20003,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "debug",
					ContainerPort: 20002,
					HostPort:      20004,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "debug-1",
					ContainerPort: 5858,
					HostPort:      20005,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20003,
					ContainerPort: 20001,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20004,
					ContainerPort: 20002,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug-1",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20005,
					ContainerPort: 5858,
					IsDebug:       true,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint + container ports known - with debug / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 20001,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 20002,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug-1",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				debug:            true,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20003,
					HostPort:      20003,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "debug",
					ContainerPort: 20004,
					HostPort:      20004,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "debug-1",
					ContainerPort: 20005,
					HostPort:      20005,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20003,
					ContainerPort: 20001,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20004,
					ContainerPort: 20002,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug-1",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20005,
					ContainerPort: 5858,
					IsDebug:       true,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint + container ports known - without debug / forwardLocalhost=false",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 20001,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 20002,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug-1",
						TargetPort: 5858,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http-1",
						TargetPort: 8080,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         false,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20002,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http-1",
					ContainerPort: 8080,
					HostPort:      20003,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 20001,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http-1",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20003,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint + container ports known - without debug / forwardLocalhost=true",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 20001,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 20002,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug-1",
						TargetPort: 5858,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http-1",
						TargetPort: 8080,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName:    devfileName,
				appName:          appName,
				debug:            false,
				forwardLocalhost: true,
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20002,
					HostPort:      20002,
					HostIP:        "127.0.0.1",
				})
				pod.Spec.Containers[1].Ports = append(pod.Spec.Containers[1].Ports, corev1.ContainerPort{
					Name:          "http-1",
					ContainerPort: 20003,
					HostPort:      20003,
					HostIP:        "127.0.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 20001,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http-1",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20003,
					ContainerPort: 8080,
					IsDebug:       false,
				},
			},
		},

		{
			name: "basic component + application endpoint + debug endpoint - with debug / custom mapping for port forwarding with container name (customForwardedPorts)",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Name = "runtime"
					cmp.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name:       "http-8080",
							TargetPort: 8080,
						},
						{
							Name:       "debug",
							TargetPort: 5858,
						},
					}

					cmp2 := baseComponent.DeepCopy()
					cmp2.Name = "tools"
					cmp2.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name:       "http-5000",
							TargetPort: 5000,
						},
					}
					_ = data.AddComponents([]v1alpha2.Component{*cmp, *cmp2})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
				customForwardedPorts: []api.ForwardedPort{
					{
						ContainerName: "runtime",
						LocalPort:     8080,
						ContainerPort: 8080,
					},
				},
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Name = "runtime"
				pod.Spec.Containers[0].Ports = []corev1.ContainerPort{
					{
						Name:          "http-8080",
						ContainerPort: 8080,
						HostPort:      8080,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
					{
						Name:          "debug",
						ContainerPort: 5858,
						HostPort:      20001,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				container2 := pod.Spec.Containers[0].DeepCopy()
				container2.Name = "tools"
				container2.Ports = []corev1.ContainerPort{
					{
						Name:          "http-5000",
						ContainerPort: 5000,
						HostPort:      20002,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				pod.Spec.Containers = append(pod.Spec.Containers, *container2)
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "http-8080",
					LocalAddress:  "127.0.0.1",
					LocalPort:     8080,
					ContainerPort: 8080,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 5858,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "tools",
					PortName:      "http-5000",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 5000,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - with debug / custom mapping for port forwarding without container name (customForwardedPorts)",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Name = "runtime"
					cmp.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name: "http-8080",

							TargetPort: 8080,
						},
						{
							Name:       "debug",
							TargetPort: 5858,
						},
					}
					cmp2 := baseComponent.DeepCopy()
					cmp2.Name = "tools"
					cmp2.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name:       "http-5000",
							TargetPort: 5000,
						},
					}

					_ = data.AddComponents([]v1alpha2.Component{*cmp, *cmp2})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
				customForwardedPorts: []api.ForwardedPort{
					{
						LocalPort:     8080,
						ContainerPort: 8080,
					},
					{
						LocalPort:     5000,
						ContainerPort: 5000,
					},
				},
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Name = "runtime"
				pod.Spec.Containers[0].Ports = []corev1.ContainerPort{
					{
						Name:          "http-8080",
						ContainerPort: 8080,
						HostPort:      8080,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					}, {
						Name:          "debug",
						ContainerPort: 5858,
						HostPort:      20001,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				container2 := pod.Spec.Containers[0].DeepCopy()
				container2.Name = "tools"
				container2.Ports = []corev1.ContainerPort{
					{
						Name:          "http-5000",
						ContainerPort: 5000,
						HostPort:      5000,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				pod.Spec.Containers = append(pod.Spec.Containers, *container2)
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "http-8080",
					LocalAddress:  "127.0.0.1",
					LocalPort:     8080,
					ContainerPort: 8080,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 5858,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "tools",
					PortName:      "http-5000",
					LocalAddress:  "127.0.0.1",
					LocalPort:     5000,
					ContainerPort: 5000,
					IsDebug:       false,
				},
			},
		},
		{
			name: "basic component + application endpoint + debug endpoint - with debug / custom mapping for port forwarding with local port in ranged [20001-30001] ports (customForwardedPorts)",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Name = "runtime"
					cmp.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name:       "http-8080",
							TargetPort: 8080,
						},
						{
							Name:       "debug",
							TargetPort: 5858,
						},
					}

					cmp2 := baseComponent.DeepCopy()
					cmp2.Name = "tools"
					cmp2.Container.Endpoints = []v1alpha2.Endpoint{
						{
							Name:       "http-9000",
							TargetPort: 9000,
						},
						{
							Name:       "http-5000",
							TargetPort: 5000,
						},
					}

					_ = data.AddComponents([]v1alpha2.Component{*cmp, *cmp2})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
				customForwardedPorts: []api.ForwardedPort{
					{
						LocalPort:     20001,
						ContainerPort: 8080,
					},
					{
						LocalPort:     20002,
						ContainerPort: 9000,
					},
					{
						LocalPort:     5000,
						ContainerPort: 5000,
					},
				},
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Name = "runtime"
				pod.Spec.Containers[0].Ports = []corev1.ContainerPort{
					{
						Name:          "http-8080",
						ContainerPort: 8080,
						HostPort:      20001,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
					{
						Name:          "debug",
						ContainerPort: 5858,
						HostPort:      20003,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				container2 := pod.Spec.Containers[0].DeepCopy()
				container2.Name = "tools"
				container2.Ports = []corev1.ContainerPort{
					{
						Name:          "http-9000",
						ContainerPort: 9000,
						HostPort:      20002,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
					{
						Name:          "http-5000",
						ContainerPort: 5000,
						HostPort:      5000,
						Protocol:      corev1.ProtocolTCP,
						HostIP:        "127.0.0.1",
					},
				}
				pod.Spec.Containers = append(pod.Spec.Containers, *container2)
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "http-8080",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20001,
					ContainerPort: 8080,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "runtime",
					PortName:      "debug",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20003,
					ContainerPort: 5858,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "tools",
					PortName:      "http-9000",
					LocalAddress:  "127.0.0.1",
					LocalPort:     20002,
					ContainerPort: 9000,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "tools",
					PortName:      "http-5000",
					LocalAddress:  "127.0.0.1",
					LocalPort:     5000,
					ContainerPort: 5000,
					IsDebug:       false,
				},
			},
		},

		{
			name: "basic component + application endpoint + debug endpoint + container ports known - with debug / using customAddress",
			args: args{
				devfileObj: func() parser.DevfileObj {
					data, _ := data.NewDevfileData(string(data.APISchemaVersion200))
					_ = data.AddCommands([]v1alpha2.Command{command})
					cmp := baseComponent.DeepCopy()
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "http",
						TargetPort: 20001,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug",
						TargetPort: 20002,
					})
					cmp.Container.Endpoints = append(cmp.Container.Endpoints, v1alpha2.Endpoint{
						Name:       "debug-1",
						TargetPort: 5858,
					})
					_ = data.AddComponents([]v1alpha2.Component{*cmp})
					return parser.DevfileObj{
						Data: data,
					}
				},
				componentName: devfileName,
				appName:       appName,
				debug:         true,
				customAddress: "192.168.0.1",
			},
			wantPod: func(basePod *corev1.Pod) *corev1.Pod {
				pod := basePod.DeepCopy()
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "http",
					ContainerPort: 20001,
					HostPort:      20003,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "192.168.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "debug",
					ContainerPort: 20002,
					HostPort:      20004,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "192.168.0.1",
				})
				pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
					Name:          "debug-1",
					ContainerPort: 5858,
					HostPort:      20005,
					Protocol:      corev1.ProtocolTCP,
					HostIP:        "192.168.0.1",
				})
				return pod
			},
			wantFwPorts: []api.ForwardedPort{
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "http",
					LocalAddress:  "192.168.0.1",
					LocalPort:     20003,
					ContainerPort: 20001,
					IsDebug:       false,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug",
					LocalAddress:  "192.168.0.1",
					LocalPort:     20004,
					ContainerPort: 20002,
					IsDebug:       true,
				},
				{
					Platform:      "podman",
					ContainerName: "mycomponent",
					PortName:      "debug-1",
					LocalAddress:  "192.168.0.1",
					LocalPort:     20005,
					ContainerPort: 5858,
					IsDebug:       true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			devfileObj := tt.args.devfileObj()
			ctx = astracontext.WithEffectiveDevfileObj(ctx, &devfileObj)
			ctx = astracontext.WithApplication(ctx, tt.args.appName)
			ctx = astracontext.WithComponentName(ctx, tt.args.componentName)
			ctx = astracontext.WithWorkingDirectory(ctx, "/tmp/dir")
			ctrl := gomock.NewController(t)
			podmanClient := podman.NewMockClient(ctrl)
			podmanClient.EXPECT().GetCapabilities().Return(tt.capabilities, nil)
			client := NewDevClient(
				nil, podmanClient, nil, nil, nil, nil, nil, nil,
			)
			got, gotFwPorts, err := client.createPodFromComponent(
				ctx,
				tt.args.debug,
				tt.args.buildCommand,
				tt.args.runCommand,
				tt.args.debugCommand,
				tt.args.forwardLocalhost,
				false,
				tt.args.customForwardedPorts,
				[]int{20001, 20002, 20003, 20004, 20005},
				tt.args.customAddress,
				devfileObj,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPodFromComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			basePod := buildBasePod(tt.args.forwardLocalhost)
			if diff := cmp.Diff(tt.wantPod(basePod), got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("createPodFromComponent() pod mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantFwPorts, gotFwPorts, cmpopts.EquateEmpty(), cmpopts.SortSlices(func(x, y api.ForwardedPort) bool { return x.ContainerName < y.ContainerName })); diff != "" {
				t.Errorf("createPodFromComponent() fwPorts mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
