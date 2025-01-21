package podmandev

import (
	"context"

	"github\.com/danielpickens/astra/pkg/dev/common"
	"k8s.io/klog"
)

func (o *DevClient) Run(
	ctx context.Context,
	commandName string,
) error {
	klog.V(4).Infof("running command %q on podman", commandName)
	return common.Run(
		ctx,
		commandName,
		o.podmanClient,
		o.execClient,
		nil, // Tastra(feloy) set when running on new container is supported on podman
		o.fs,
	)
}
