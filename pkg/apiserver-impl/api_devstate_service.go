package apiserver_impl

import (
	"context"

	openapi "github\.com/danielpickens/astra/pkg/apiserver-gen/go"
	"github\.com/danielpickens/astra/pkg/apiserver-impl/devstate"
	"github\.com/danielpickens/astra/pkg/kclient"
	"github\.com/danielpickens/astra/pkg/podman"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/state"
)

// DevstateApiService is a service that implements the logic for the DevstateApiServicer
// This service should implement the business logic for every endpoint for the DevstateApi API.
// Include any external packages or services that will be required by this service.
type DevstateApiService struct {
	cancel           context.CancelFunc
	pushWatcher      chan<- struct{}
	kubeClient       kclient.ClientInterface
	podmanClient     podman.Client
	stateClient      state.Client
	preferenceClient preference.Client

	devfileState devstate.DevfileState
}

// NewDevstateApiService creates a devstate api service
func NewDevstateApiService(
	cancel context.CancelFunc,
	pushWatcher chan<- struct{},
	kubeClient kclient.ClientInterface,
	podmanClient podman.Client,
	stateClient state.Client,
	preferenceClient preference.Client,
) openapi.DevstateApiServicer {
	return &DevstateApiService{
		cancel:           cancel,
		pushWatcher:      pushWatcher,
		kubeClient:       kubeClient,
		podmanClient:     podmanClient,
		stateClient:      stateClient,
		preferenceClient: preferenceClient,

		devfileState: devstate.NewDevfileState(),
	}
}
