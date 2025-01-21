package apiserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	apiserver_impl "github\.com/danielpickens/astra/pkg/apiserver-impl"
	"github\.com/danielpickens/astra/pkg/libdevfile"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
)

const (
	RecommendedCommandName = "api-server"
)

type ApiServerOptions struct {
	clientset *clientset.Clientset

	// Flags
	randomPortsFlag bool
	portFlag        int
}

func NewApiServerOptions() *ApiServerOptions {
	return &ApiServerOptions{}
}

var _ genericclioptions.Runnable = (*ApiServerOptions)(nil)
var _ genericclioptions.SignalHandler = (*ApiServerOptions)(nil)
var _ genericclioptions.Cleanuper = (*ApiServerOptions)(nil)

func (o *ApiServerOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

func (o *ApiServerOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) error {
	return nil
}

func (o *ApiServerOptions) Validate(ctx context.Context) error {
	if o.randomPortsFlag && o.portFlag != 0 {
		return errors.New("--random-ports and --port cannot be used together")
	}
	return nil
}

func (o *ApiServerOptions) Run(ctx context.Context) (err error) {
	err = o.clientset.StateClient.Init(ctx)
	if err != nil {
		err = fmt.Errorf("unable to save state file: %w", err)
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		devfileObj  = astracontext.GetEffectiveDevfileObj(ctx)
		devfilePath = astracontext.GetDevfilePath(ctx)
	)

	devfileFiles, err := libdevfile.GetReferencedLocalFiles(*devfileObj)
	if err != nil {
		return err
	}
	devfileFiles = append(devfileFiles, devfilePath)
	_, err = apiserver_impl.StartServer(
		ctx,
		cancel,
		o.randomPortsFlag,
		o.portFlag,
		devfilePath,
		devfileFiles,
		o.clientset.FS,
		nil,
		nil,
		o.clientset.StateClient,
		o.clientset.PreferenceClient,
		o.clientset.InformerClient,
	)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (o *ApiServerOptions) Cleanup(ctx context.Context, commandError error) error {
	err := o.clientset.StateClient.SaveExit(ctx)
	if err != nil {
		klog.V(1).Infof("unable to persist dev state: %v", err)
	}
	return nil
}

func (o *ApiServerOptions) HandleSignal(ctx context.Context, cancelFunc context.CancelFunc) error {
	cancelFunc()
	return nil
}

// NewCmdApiServer implements the astra api-server command
func NewCmdApiServer(ctx context.Context, name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	o := NewApiServerOptions()
	apiserverCmd := &cobra.Command{
		Use:   name,
		Short: "Start the API server",
		Long:  "Start the API server",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
	}
	clientset.Add(apiserverCmd,
		clientset.FILESYSTEM,
		clientset.INFORMER,
		clientset.STATE,
		clientset.PREFERENCE,
	)
	apiserverCmd.Flags().BoolVar(&o.randomPortsFlag, "random-ports", false, "Assign a random API Server port.")
	apiserverCmd.Flags().IntVar(&o.portFlag, "port", 0, "Define custom port for API Server.")
	return apiserverCmd
}
