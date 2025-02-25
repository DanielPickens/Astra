package kubedev

import (
	"context"
	"fmt"

	devfilev1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"

	"github\.com/danielpickens/astra/pkg/binding"
	_delete "github\.com/danielpickens/astra/pkg/component/delete"
	"github\.com/danielpickens/astra/pkg/configAutomount"
	"github\.com/danielpickens/astra/pkg/dev"
	"github\.com/danielpickens/astra/pkg/dev/common"
	"github\.com/danielpickens/astra/pkg/devfile"
	"github\.com/danielpickens/astra/pkg/devfile/location"
	"github\.com/danielpickens/astra/pkg/exec"
	"github\.com/danielpickens/astra/pkg/kclient"
	"github\.com/danielpickens/astra/pkg/portForward"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/sync"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
	"github\.com/danielpickens/astra/pkg/watch"

	"k8s.io/klog"
)

type DevClient struct {
	kubernetesClient      kclient.ClientInterface
	prefClient            preference.Client
	portForwardClient     portForward.Client
	watchClient           watch.Client
	bindingClient         binding.Client
	syncClient            sync.Client
	filesystem            filesystem.Filesystem
	execClient            exec.Client
	deleteClient          _delete.Client
	configAutomountClient configAutomount.Client

	// deploymentExists is true when the deployment is already created when calling createComponents
	deploymentExists bool
	// portsChanged is true of ports have changed since the last call to createComponents
	portsChanged bool
	// portsToForward lists the port to forward during inner loop (Tastra move port forward to createComponents)
	portsToForward map[string][]devfilev1.Endpoint
}

var _ dev.Client = (*DevClient)(nil)

func NewDevClient(
	kubernetesClient kclient.ClientInterface,
	prefClient preference.Client,
	portForwardClient portForward.Client,
	watchClient watch.Client,
	bindingClient binding.Client,
	syncClient sync.Client,
	filesystem filesystem.Filesystem,
	execClient exec.Client,
	deleteClient _delete.Client,
	configAutomountClient configAutomount.Client,
) *DevClient {
	return &DevClient{
		kubernetesClient:      kubernetesClient,
		prefClient:            prefClient,
		portForwardClient:     portForwardClient,
		watchClient:           watchClient,
		bindingClient:         bindingClient,
		syncClient:            syncClient,
		filesystem:            filesystem,
		execClient:            execClient,
		deleteClient:          deleteClient,
		configAutomountClient: configAutomountClient,
	}
}

func (o *DevClient) Start(
	ctx context.Context,
	options dev.StartOptions,
) error {
	klog.V(4).Infoln("Creating new adapter")

	var (
		componentStatus = watch.ComponentStatus{
			ImageComponentsAutoApplied: make(map[string]devfilev1.ImageComponent),
		}
	)

	klog.V(4).Infoln("Creating inner-loop resources for the component")

	watchParameters := watch.WatchParameters{
		StartOptions:        options,
		DevfileWatchHandler: o.regenerateAdapterAndPush,
		WatchCluster:        true,
	}

	return o.watchClient.WatchAndPush(ctx, watchParameters, componentStatus)
}

// RegenerateAdapterAndPush get the new devfile and pushes the files to remote pod
func (o *DevClient) regenerateAdapterAndPush(ctx context.Context, pushParams common.PushParameters, componentStatus *watch.ComponentStatus) error {

	devObj, err := devfile.ParseAndValidateFromFileWithVariables(location.DevfileLocation(o.filesystem, ""), pushParams.StartOptions.Variables, o.prefClient.GetImageRegistry(), true)
	if err != nil {
		return fmt.Errorf("unable to read devfile: %w", err)
	}

	pushParams.Devfile = devObj

	err = o.reconcile(ctx, pushParams, componentStatus)
	if err != nil {
		return fmt.Errorf("watch command was unable to push component: %w", err)
	}
	return nil
}
