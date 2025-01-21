package component

import (
	"context"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	devfilev1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"k8s.io/klog"

	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	"github\.com/danielpickens/astra/pkg/configAutomount"
	"github\.com/danielpickens/astra/pkg/devfile/image"
	"github\.com/danielpickens/astra/pkg/exec"
	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/libdevfile"
	"github\.com/danielpickens/astra/pkg/log"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/platform"
	"github\.com/danielpickens/astra/pkg/remotecmd"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
)

type runHandler struct {
	ctx                   context.Context
	platformClient        platform.Client
	execClient            exec.Client
	configAutomountClient configAutomount.Client
	podName               string
	ComponentExists       bool
	containersRunning     []string
	msg                   string
	directRun             bool

	fs           filesystem.Filesystem
	imageBackend image.Backend

	devfile parser.DevfileObj
	path    string
}

var _ libdevfile.Handler = (*runHandler)(nil)

type HandlerOptions struct {
	PodName           string
	ComponentExists   bool
	ContainersRunning []string
	Msg               string
	DirectRun         bool

	// For apply Kubernetes / Openshift
	Devfile parser.DevfileObj
	Path    string
}

func NewRunHandler(
	ctx context.Context,
	platformClient platform.Client,
	execClient exec.Client,
	configAutomountClient configAutomount.Client,

	// For building images
	fs filesystem.Filesystem,
	imageBackend image.Backend,

	options HandlerOptions,

) *runHandler {
	return &runHandler{
		ctx:                   ctx,
		platformClient:        platformClient,
		execClient:            execClient,
		configAutomountClient: configAutomountClient,
		podName:               options.PodName,
		ComponentExists:       options.ComponentExists,
		containersRunning:     options.ContainersRunning,
		msg:                   options.Msg,
		directRun:             options.DirectRun,

		fs:           fs,
		imageBackend: imageBackend,

		devfile: options.Devfile,
		path:    options.Path,
	}
}

func (a *runHandler) ApplyImage(img devfilev1.Component) error {
	return image.BuildPushSpecificImage(a.ctx, a.imageBackend, a.fs, img, envcontext.GetEnvConfig(a.ctx).PushImages)
}

func (a *runHandler) ApplyKubernetes(kubernetes devfilev1.Component, kind v1alpha2.CommandGroupKind) error {
	var (
		componentName = astracontext.GetComponentName(a.ctx)
		appName       = astracontext.GetApplication(a.ctx)
	)
	mode := astralabels.ComponentDevMode
	if kind == v1alpha2.DeployCommandGroupKind {
		mode = astralabels.ComponentDeployMode
	}
	switch platform := a.platformClient.(type) {
	case kclient.ClientInterface:
		return ApplyKubernetes(mode, appName, componentName, a.devfile, kubernetes, platform, a.path)
	default:
		klog.V(4).Info("apply kubernetes/Openshift commands are not implemented on podman")
		log.Warningf("Apply Kubernetes/Openshift components are not supported on Podman. Skipping: %v.", kubernetes.Name)
		return nil
	}
}

func (a *runHandler) ApplyOpenShift(openshift devfilev1.Component, kind v1alpha2.CommandGroupKind) error {
	return a.ApplyKubernetes(openshift, kind)
}

func (a *runHandler) ExecuteNonTerminatingCommand(ctx context.Context, command devfilev1.Command) error {
	var (
		componentName = astracontext.GetComponentName(a.ctx)
		appName       = astracontext.GetApplication(a.ctx)
	)
	if isContainerRunning(command.Exec.Component, a.containersRunning) {
		return ExecuteRunCommand(ctx, a.execClient, a.platformClient, command, a.ComponentExists, a.podName, appName, componentName)
	}
	switch platform := a.platformClient.(type) {
	case kclient.ClientInterface:
		return ExecuteInNewContainer(ctx, platform, a.configAutomountClient, a.devfile, componentName, appName, command)
	default:
		klog.V(4).Info("executing a command in a new container is not implemented on podman")
		log.Warningf("executing a command in a new container is not implemented on podman. Skipping: %v.", command.Id)
		return nil
	}
}

func (a *runHandler) ExecuteTerminatingCommand(ctx context.Context, command devfilev1.Command) error {
	var (
		componentName = astracontext.GetComponentName(a.ctx)
		appName       = astracontext.GetApplication(a.ctx)
	)
	if isContainerRunning(command.Exec.Component, a.containersRunning) {
		return ExecuteTerminatingCommand(ctx, a.execClient, a.platformClient, command, a.ComponentExists, a.podName, appName, componentName, a.msg, a.directRun)
	}
	switch platform := a.platformClient.(type) {
	case kclient.ClientInterface:
		return ExecuteInNewContainer(ctx, platform, a.configAutomountClient, a.devfile, componentName, appName, command)
	default:
		klog.V(4).Info("executing a command in a new container is not implemented on podman")
		log.Warningf("executing a command in a new container is not implemented on podman. Skipping: %v.", command.Id)
		return nil
	}
}

// IsRemoteProcessForCommandRunning returns true if the command is running
func (a *runHandler) IsRemoteProcessForCommandRunning(ctx context.Context, command devfilev1.Command, podName string) (bool, error) {
	remoteProcess, err := remotecmd.NewKubeExecProcessHandler(a.execClient).GetProcessInfoForCommand(ctx, remotecmd.CommandDefinition{Id: command.Id}, podName, command.Exec.Component)
	if err != nil {
		return false, err
	}

	return remoteProcess.Status == remotecmd.Running, nil
}

func isContainerRunning(container string, containers []string) bool {
	for _, cnt := range containers {
		if container == cnt {
			return true
		}
	}
	return false
}
