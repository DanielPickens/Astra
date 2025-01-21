package common

import (
	"context"
	"fmt"

	"github\.com/danielpickens/astra/pkg/component"
	"github\.com/danielpickens/astra/pkg/configAutomount"
	"github\.com/danielpickens/astra/pkg/devfile/image"
	"github\.com/danielpickens/astra/pkg/exec"
	"github\.com/danielpickens/astra/pkg/libdevfile"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/platform"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
)

func Run(
	ctx context.Context,
	commandName string,
	platformClient platform.Client,
	execClient exec.Client,
	configAutomountClient configAutomount.Client,
	filesystem filesystem.Filesystem,
) error {
	var (
		componentName = astracontext.GetComponentName(ctx)
		devfileObj    = astracontext.GetEffectiveDevfileObj(ctx)
		devfilePath   = astracontext.GetDevfilePath(ctx)
	)

	pod, err := platformClient.GetPodUsingComponentName(componentName)
	if err != nil {
		return fmt.Errorf("unable to get pod for component %s: %w. Please check the command 'astra dev' is running", componentName, err)
	}

	handler := component.NewRunHandler(
		ctx,
		platformClient,
		execClient,
		configAutomountClient,
		filesystem,
		image.SelectBackend(ctx),
		component.HandlerOptions{
			PodName:           pod.Name,
			ContainersRunning: component.GetContainersNames(pod),
			Msg:               "Executing command in container",
			DirectRun:         true,
			Devfile:           *devfileObj,
			Path:              devfilePath,
		},
	)

	return libdevfile.ExecuteCommandByName(ctx, *devfileObj, commandName, handler, false)
}
