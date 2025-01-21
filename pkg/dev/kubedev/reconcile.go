package kubedev

import (
	"context"

	"github\.com/danielpickens/astra/pkg/dev/common"
	"github\.com/danielpickens/astra/pkg/watch"
)

// reconcile updates the component if a matching component exists or creates one if it doesn't exist
// Once the component has started, it will sync the source code to it.
// The componentStatus will be modified to reflect the status of the component when the function returns
func (o *DevClient) reconcile(ctx context.Context, parameters common.PushParameters, componentStatus *watch.ComponentStatus) (err error) {

	// pastraK indicates if the pod is ready to use for the inner loop
	var pastraK bool
	pastraK, err = o.createComponents(ctx, parameters, componentStatus)
	if err != nil {
		return err
	}
	if !pastraK {
		return nil
	}

	return o.innerloop(ctx, parameters, componentStatus)
}
