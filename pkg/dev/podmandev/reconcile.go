package podmandev

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	devfilev1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/fatih/color"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/component"
	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	"github\.com/danielpickens/astra/pkg/dev"
	"github\.com/danielpickens/astra/pkg/dev/common"
	"github\.com/danielpickens/astra/pkg/devfile/image"
	"github\.com/danielpickens/astra/pkg/libdevfile"
	"github\.com/danielpickens/astra/pkg/log"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/port"
	"github\.com/danielpickens/astra/pkg/watch"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog"
)

func (o *DevClient) reconcile(
	ctx context.Context,
	parameters common.PushParameters,
	componentStatus *watch.ComponentStatus,
) error {
	var (
		componentName = astracontext.GetComponentName(ctx)
		devfilePath   = astracontext.GetDevfilePath(ctx)
		path          = filepath.Dir(devfilePath)
		options       = parameters.StartOptions
		devfileObj    = parameters.Devfile
	)

	o.warnAboutK8sComponents(devfileObj)

	err := o.buildPushAutoImageComponents(ctx, devfileObj)
	if err != nil {
		return err
	}

	pod, fwPorts, err := o.deployPod(ctx, options, devfileObj)
	if err != nil {
		return err
	}
	o.deployedPod = pod
	componentStatus.SetState(watch.StateReady)

	execRequired, err := o.syncFiles(ctx, options, pod, path)
	if err != nil {
		return err
	}

	// PostStart events from the devfile will only be executed when the component
	// didn't previously exist
	if !componentStatus.PostStartEventsDone && libdevfile.HasPostStartEvents(devfileObj) {
		execHandler := component.NewRunHandler(
			ctx,
			o.podmanClient,
			o.execClient,
			nil, // Tastra(feloy) set this value when we want to support exec on new container on podman
			// Tastra(feloy) set these values when we want to support Apply Image/Kubernetes/OpenShift commands for PostStart commands
			nil, nil,
			component.HandlerOptions{
				PodName:           pod.Name,
				ContainersRunning: component.GetContainersNames(pod),
				Msg:               "Executing post-start command in container",
			},
		)
		err = libdevfile.ExecPostStartEvents(ctx, devfileObj, execHandler)
		if err != nil {
			return err
		}
	}
	componentStatus.PostStartEventsDone = true

	innerLoopWithCommands := !parameters.StartOptions.SkipCommands
	var hasRunOrDebugCmd bool
	if innerLoopWithCommands {
		if execRequired {
			doExecuteBuildCommand := func() error {
				execHandler := component.NewRunHandler(
					ctx,
					o.podmanClient,
					o.execClient,
					nil, // Tastra(feloy) set this value when we want to support exec on new container on podman

					// Tastra(feloy) set these values when we want to support Apply Image/Kubernetes/OpenShift commands for PreStop events
					nil, nil, component.HandlerOptions{
						PodName:           pod.Name,
						ComponentExists:   componentStatus.RunExecuted,
						ContainersRunning: component.GetContainersNames(pod),
						Msg:               "Building your application in container",
					},
				)
				return libdevfile.Build(ctx, devfileObj, options.BuildCommand, execHandler)
			}

			err = doExecuteBuildCommand()
			if err != nil {
				return err
			}

			cmdKind := devfilev1.RunCommandGroupKind
			cmdName := options.RunCommand
			if options.Debug {
				cmdKind = devfilev1.DebugCommandGroupKind
				cmdName = options.DebugCommand
			}
			_, hasRunOrDebugCmd, err = libdevfile.GetCommand(parameters.Devfile, cmdName, cmdKind)
			if err != nil {
				return err
			}

			if hasRunOrDebugCmd {
				cmdHandler := component.NewRunHandler(
					ctx,
					o.podmanClient,
					o.execClient,
					nil, // Tastra(feloy) set this value when we want to support exec on new container on podman

					o.fs,
					image.SelectBackend(ctx),

					// Tastra(feloy) set to deploy Kubernetes/Openshift components
					component.HandlerOptions{
						PodName:           pod.Name,
						ComponentExists:   componentStatus.RunExecuted,
						ContainersRunning: component.GetContainersNames(pod),
					},
				)
				err = libdevfile.ExecuteCommandByNameAndKind(ctx, devfileObj, cmdName, cmdKind, cmdHandler, false)
				if err != nil {
					return err
				}
				componentStatus.RunExecuted = true
			} else {
				msg := fmt.Sprintf("Missing default %v command", cmdKind)
				if cmdName != "" {
					msg = fmt.Sprintf("Missing %v command with name %q", cmdKind, cmdName)
				}
				log.Warning(msg)
			}
		}
	}

	if innerLoopWithCommands && hasRunOrDebugCmd && len(fwPorts) != 0 {
		// Check that the application is actually listening on the ports declared in the Devfile, so we are sure that port-forwarding will work
		appReadySpinner := log.Spinner("Waiting for the application to be ready")
		err = o.checkAppPorts(ctx, pod.Name, fwPorts)
		appReadySpinner.End(err == nil)
		if err != nil {
			log.Warningf("Port forwarding might not work correctly: %v", err)
			log.Warning("Running `astra logs --follow --platform podman` might help in identifying the problem.")
			fmt.Fprintln(options.Out)
		}
	}

	// By default, Podman will not forward to container applications listening on the loopback interface.
	// So we are trying to detect such cases and act accordingly.
	// See https://github\.com/danielpickens/astra/issues/6510#issuecomment-1439986558
	err = o.handleLoopbackPorts(ctx, options, pod, fwPorts)
	if err != nil {
		return err
	}

	if options.ForwardLocalhost {
		// Port-forwarding is enabled by executing dedicated socat commands
		err = o.portForwardClient.StartPortForwarding(ctx, devfileObj, componentName, options.Debug, options.RandomPorts, options.Out, options.ErrOut, fwPorts, options.CustomAddress)
		if err != nil {
			return common.NewErrPortForward(err)
		}
	} // else port-forwarding is done via the main container ports in the pod spec

	for _, fwPort := range fwPorts {
		s := fmt.Sprintf("Forwarding from %s:%d -> %d", fwPort.LocalAddress, fwPort.LocalPort, fwPort.ContainerPort)
		fmt.Fprintf(options.Out, " -  %s", log.SboldColor(color.FgGreen, s))
	}
	err = o.stateClient.SetForwardedPorts(ctx, fwPorts)
	if err != nil {
		return err
	}

	componentStatus.SetState(watch.StateReady)
	return nil
}

// warnAboutApplyComponents prints a warning if the Devfile contains standalone K8s components (not referenced by any Apply commands). These resources are currently applied when running in the cluster mode, but not on Podman.
func (o *DevClient) warnAboutK8sComponents(devfileObj parser.DevfileObj) {
	var components []string
	// get all standalone k8s components for a given commandGK
	k8sComponents, _ := libdevfile.GetK8sAndOcComponentsToPush(devfileObj, false)

	if len(k8sComponents) == 0 {
		return
	}

	for _, comp := range k8sComponents {
		components = append(components, comp.Name)
	}

	log.Warningf("Kubernetes components are not supported on Podman. Skipping: %v.", strings.Join(components, ", "))
}

func (o *DevClient) buildPushAutoImageComponents(ctx context.Context, devfileObj parser.DevfileObj) error {
	components, err := libdevfile.GetImageComponentsToPushAutomatically(devfileObj)
	if err != nil {
		return err
	}

	for _, c := range components {
		err = image.BuildPushSpecificImage(ctx, image.SelectBackend(ctx), o.fs, c, envcontext.GetEnvConfig(ctx).PushImages)
		if err != nil {
			return err
		}
	}
	return nil
}

// deployPod deploys the component as a Pod in podman
func (o *DevClient) deployPod(ctx context.Context, options dev.StartOptions, devfileObj parser.DevfileObj) (*corev1.Pod, []api.ForwardedPort, error) {

	spinner := log.Spinner("Deploying pod")
	defer spinner.End(false)

	pod, fwPorts, err := o.createPodFromComponent(
		ctx,
		options.Debug,
		options.BuildCommand,
		options.RunCommand,
		options.DebugCommand,
		options.ForwardLocalhost,
		options.RandomPorts,
		options.CustomForwardedPorts,
		o.usedPorts,
		options.CustomAddress,
		devfileObj,
	)
	if err != nil {
		return nil, nil, err
	}
	o.usedPorts = getUsedPorts(fwPorts)

	if equality.Semantic.DeepEqual(o.deployedPod, pod) {
		klog.V(4).Info("pod is already deployed as required")
		spinner.End(true)
		return o.deployedPod, fwPorts, nil
	}

	// Delete previous pod, if running
	if o.deployedPod != nil {
		err = o.podmanClient.CleanupPodResources(o.deployedPod, false)
		if err != nil {
			return nil, nil, err
		}
	} else {
		err = o.checkVolumesFree(pod)
		if err != nil {
			return nil, nil, err
		}
	}

	err = o.podmanClient.PlayKube(pod)
	if err != nil {
		// there are cases when pod is created even if there is an error with the pod def; for e.g. incorrect image
		if podMap, _ := o.podmanClient.PodLs(); podMap[pod.Name] {
			o.deployedPod = &corev1.Pod{}
			o.deployedPod.SetName(pod.Name)
		}
		return nil, nil, err
	}

	spinner.End(true)
	return pod, fwPorts, nil
}

func (o *DevClient) checkAppPorts(ctx context.Context, podName string, portsToFwd []api.ForwardedPort) error {
	containerPortsMapping := make(map[string][]int)
	for _, p := range portsToFwd {
		containerPortsMapping[p.ContainerName] = append(containerPortsMapping[p.ContainerName], p.ContainerPort)
	}
	return port.CheckAppPortsListening(ctx, o.execClient, podName, containerPortsMapping, 1*time.Minute)
}

// handleLoopbackPorts tries to detect if any of the ports to forward (in fwPorts) is actually bound to the loopback interface within the specified pod.
// If that is the case, it will either return an error if options.IgnoreLocalhost is false, or no error otherwise.
//
// Note that this method should be called after the process representing the application (run or debug command) is actually started in the pod.
func (o *DevClient) handleLoopbackPorts(ctx context.Context, options dev.StartOptions, pod *corev1.Pod, fwPorts []api.ForwardedPort) error {
	if len(pod.Spec.Containers) == 0 {
		return nil
	}

	loopbackPorts, err := port.DetectRemotePortsBoundOnLoopback(ctx, o.execClient, pod.Name, pod.Spec.Containers[0].Name, fwPorts)
	if err != nil {
		log.Warningf("unable to detect container ports bound on the loopback interface: %v", err)
	}

	if len(loopbackPorts) == 0 {
		return nil
	}

	klog.V(5).Infof("detected %d ports bound on the loopback interface in the pod: %v", len(loopbackPorts), loopbackPorts)
	list := make([]string, 0, len(loopbackPorts))
	for _, p := range loopbackPorts {
		list = append(list, fmt.Sprintf("%s (%d)", p.PortName, p.ContainerPort))
	}
	msg := fmt.Sprintf(`Detected that the following port(s) can be reached only via the container loopback interface: %s.
Port forwarding on Podman currently does not work with applications listening on the loopback interface.
Either change the application to make those port(s) reachable on all interfaces (0.0.0.0), or rerun 'astra dev' with `, strings.Join(list, ", "))
	if options.IgnoreLocalhost {
		msg += "'--forward-localhost' to make port-forwarding work with such ports."
	} else {
		msg += `any of the following options:
- --ignore-localhost: no error will be returned by astra, but forwarding to those ports might not work on Podman.
- --forward-localhost: astra will inject a dedicated side container to redirect traffic to such ports.`
	}
	if options.IgnoreLocalhost {
		// ForwardLocalhost should not be true at this point.
		log.Warningf(msg)
	} else if !options.ForwardLocalhost {
		log.Errorf(msg)
		return errors.New("cannot make port forwarding work with ports bound to the loopback interface only")
	}

	return nil
}
