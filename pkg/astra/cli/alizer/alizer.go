package alizer

import (
	"context"
	"errors"
	"fmt"

	"github\.com/danielpickens/astra/pkg/alizer"
	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/commonflags"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	astrautil "github\.com/danielpickens/astra/pkg/astra/util"

	"github.com/spf13/cobra"
)

const RecommendedCommandName = "analyze"

type AlizerOptions struct {
	clientset *clientset.Clientset
}

var _ genericclioptions.Runnable = (*AlizerOptions)(nil)
var _ genericclioptions.JsonOutputter = (*AlizerOptions)(nil)

// NewAlizerOptions creates a new AlizerOptions instance
func NewAlizerOptions() *AlizerOptions {
	return &AlizerOptions{}
}

func (o *AlizerOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

func (o *AlizerOptions) UseDevfile(ctx context.Context, cmdline cmdline.Cmdline, args []string) bool {
	return false
}

func (o *AlizerOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) (err error) {
	return nil
}

func (o *AlizerOptions) Validate(ctx context.Context) error {
	return nil
}

func (o *AlizerOptions) Run(ctx context.Context) (err error) {
	return errors.New("this command can be run with json output only, please use the flag: -o json")
}

// RunForJsonOutput contains the logic for the astra command
func (o *AlizerOptions) RunForJsonOutput(ctx context.Context) (out interface{}, err error) {
	workingDir := astracontext.GetWorkingDirectory(ctx)
	detected, err := o.clientset.AlizerClient.DetectFramework(ctx, workingDir)
	if err != nil {
		//revive:disable:error-strings This is a top-level error message displayed as is to the end user
		return nil, fmt.Errorf("No valid devfile found for project in %s: %w", workingDir, err)
	}
	appPorts, err := o.clientset.AlizerClient.DetectPorts(workingDir)
	if err != nil {
		return nil, err
	}
	name, err := o.clientset.AlizerClient.DetectName(workingDir)
	if err != nil {
		return nil, err
	}
	result := alizer.NewDetectionResult(detected.Type, detected.Registry, appPorts, detected.DefaultVersion, name)
	return []api.DetectionResult{*result}, nil
}

func NewCmdAlizer(name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	o := NewAlizerOptions()
	alizerCmd := &cobra.Command{
		Use:         name,
		Short:       "Detect devfile to use based on files present in current directory",
		Long:        "Detect devfile to use based on files present in current directory",
		Args:        cobra.MaximumNArgs(0),
		Annotations: map[string]string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
	}
	clientset.Add(alizerCmd, clientset.ALIZER, clientset.FILESYSTEM)
	util.SetCommandGroup(alizerCmd, util.UtilityGroup)
	commonflags.UseOutputFlag(alizerCmd)
	alizerCmd.SetUsageTemplate(astrautil.CmdUsageTemplate)
	return alizerCmd
}
