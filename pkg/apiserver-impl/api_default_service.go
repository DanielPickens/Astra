package apiserver_impl

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	openapi "github\.com/danielpickens/astra/pkg/apiserver-gen/go"
	"github\.com/danielpickens/astra/pkg/apiserver-impl/devstate"
	"github\.com/danielpickens/astra/pkg/component/describe"
	"github\.com/danielpickens/astra/pkg/devfile"
	"github\.com/danielpickens/astra/pkg/devfile/validate"
	"github\.com/danielpickens/astra/pkg/kclient"
	fcontext "github\.com/danielpickens/astra/pkg/astra/commonflags/context"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/podman"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/segment"
	scontext "github\.com/danielpickens/astra/pkg/segment/context"
	"github\.com/danielpickens/astra/pkg/state"
	"k8s.io/klog"
)

// DefaultApiService is a service that implements the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
	cancel           context.CancelFunc
	pushWatcher      chan<- struct{}
	kubeClient       kclient.ClientInterface
	podmanClient     podman.Client
	stateClient      state.Client
	preferenceClient preference.Client

	devfileState devstate.DevfileState
}

// NewDefaultApiService creates a default api service
func NewDefaultApiService(
	cancel context.CancelFunc,
	pushWatcher chan<- struct{},
	kubeClient kclient.ClientInterface,
	podmanClient podman.Client,
	stateClient state.Client,
	preferenceClient preference.Client,
) openapi.DefaultApiServicer {
	return &DefaultApiService{
		cancel:           cancel,
		pushWatcher:      pushWatcher,
		kubeClient:       kubeClient,
		podmanClient:     podmanClient,
		stateClient:      stateClient,
		preferenceClient: preferenceClient,

		devfileState: devstate.NewDevfileState(),
	}
}

// ComponentCommandPost -
func (s *DefaultApiService) ComponentCommandPost(ctx context.Context, componentCommandPostRequest openapi.ComponentCommandPostRequest) (openapi.ImplResponse, error) {
	switch componentCommandPostRequest.Name {
	case "push":
		select {
		case s.pushWatcher <- struct{}{}:
			return openapi.Response(http.StatusOK, openapi.GeneralSuccess{
				Message: "push was successfully executed",
			}), nil
		default:
			return openapi.Response(http.StatusTooManyRequests, openapi.GeneralError{
				Message: "a push operation is not possible at this time. Please retry later",
			}), nil
		}

	default:
		return openapi.Response(http.StatusBadRequest, openapi.GeneralError{
			Message: fmt.Sprintf("command name %q not supported. Supported values are: %q", componentCommandPostRequest.Name, "push"),
		}), nil
	}
// ComponentGet -
func (s *DefaultApiService) ComponentGet(ctx context.Context) (openapi.ImplResponse, error) {
	value, _, err := describe.DescribeDevfileComponent(ctx, s.kubeClient, s.podmanClient, s.stateClient)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
			Message: fmt.Sprintf("error getting the description of the component: %s", err),
		}), nil
	}

	// Check for another status response in openapi.Response property object
	if value.Status != "success" {
		return openapi.Response(http.StatusConflict, openapi.GeneralError{
			Message: fmt.Sprintf("component description status is not successful: %s", value.Status),
		}), nil
	}

	return openapi.Response(http.StatusOK, value), nil
}

// InstanceDelete -
func (s *DefaultApiService) InstanceDelete(ctx context.Context) (openapi.ImplResponse, error) {
	s.cancel()
	return openapi.Response(http.StatusOK, openapi.GeneralSuccess{
		Message: fmt.Sprintf("'astra dev' instance with pid: %d is shutting down.", astracontext.GetPID(ctx)),
	}), nil
}

// InstanceGet -
func (s *DefaultApiService) InstanceGet(ctx context.Context) (openapi.ImplResponse, error) {
	response := openapi.InstanceGet200Response{
		Pid:                int32(astracontext.GetPID(ctx)),
		ComponentDirectory: astracontext.GetWorkingDirectory(ctx),
	}
	return openapi.Response(http.StatusOK, response), nil
}

func (s *DefaultApiService) DevfileGet(ctx context.Context) (openapi.ImplResponse, error) {
	devfilePath := astracontext.GetDevfilePath(ctx)
	content, err := os.ReadFile(devfilePath)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
			Message: fmt.Sprintf("error getting Devfile content: %s", err),
		}), nil
	}
	return openapi.Response(http.StatusOK, openapi.DevfileGet200Response{
		Content: string(content),
	}), nil

}

func (s *DefaultApiService) DevfilePut(ctx context.Context, params openapi.DevfilePutRequest) (openapi.ImplResponse, error) {

	tmpdir, err := func() (string, error) {
		dir, err := os.MkdirTemp("", "astra")
		if err != nil {
			return "", err
		}
		return dir, os.WriteFile(filepath.Join(dir, "devfile.yaml"), []byte(params.Content), 0600)
	}()
	defer func() {
		if tmpdir != "" {
			err = os.RemoveAll(tmpdir)
			if err != nil {
				klog.V(1).Infof("Error deleting temp directory %q: %s", tmpdir, err)
			}
		}
	}()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
			Message: fmt.Sprintf("error saving temp Devfile: %s", err),
		}), nil
	}

	err = s.validateDevfile(ctx, tmpdir)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
			Message: fmt.Sprintf("error validating Devfile: %s", err),
		}), nil
	}

	devfilePath := astracontext.GetDevfilePath(ctx)
	err = os.WriteFile(devfilePath, []byte(params.Content), 0600)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
			Message: fmt.Sprintf("error writing Devfile content to %q: %s", devfilePath, err),
		}), nil
	}

	return openapi.Response(http.StatusOK, openapi.GeneralSuccess{
		Message: "devfile has been successfully written to disk",
	}), nil

}

func (s *DefaultApiService) validateDevfile(ctx context.Context, dir string) error {
	var (
		variables     = fcontext.GetVariables(ctx)
		imageRegistry = s.preferenceClient.GetImageRegistry()
	)
	devObj, err := devfile.ParseAndValidateFromFileWithVariables(dir, variables, imageRegistry, false)
	if err != nil {
		return fmt.Errorf("failed to parse the devfile: %w", err)
	}
	return validate.ValidateDevfileData(devObj.Data)
}

func (s *DefaultApiService) TelemetryGet(ctx context.Context) (openapi.ImplResponse, error) {
	var (
		enabled = scontext.GetTelemetryStatus(ctx)
		apikey  string
		userid  string
	)
	if enabled {
		apikey = segment.GetApikey()
		var err error
		userid, err = segment.GetUserIdentity(segment.GetTelemetryFilePath())
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, openapi.GeneralError{
				Message: fmt.Sprintf("error getting telemetry data: %s", err),
			}), nil
		}
	}

	return openapi.Response(http.StatusOK, openapi.TelemetryResponse{
		Enabled: enabled,
		Apikey:  apikey,
		Userid:  userid,
	}), nil
}
