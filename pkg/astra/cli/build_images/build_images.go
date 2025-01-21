package build_images

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github\.com/danielpickens/astra/pkg/devfile/image"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/commonflags"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	astrautil "github\.com/danielpickens/astra/pkg/astra/util"
)

// RecommendedCommandName is the recommended command name
const RecommendedCommandName = "build-images"

// BuildImagesOptions encapsulates the options for the astra command
type BuildImagesOptions struct {
	// Clients
	clientset *clientset.Clientset

	// Flags
	pushFlag bool
}

var _ genericclioptions.Runnable = (*BuildImagesOptions)(nil)

var buildImagesExample = templates.Examples(`
  # Build images defined in the devfile
  %[1]s

  # Build images and push them to their registries
  %[1]s --push
`)

// NewBuildImagesOptions creates a new BuildImagesOptions instance
func NewBuildImagesOptions() *BuildImagesOptions {
	return &BuildImagesOptions{}
}

func (o *BuildImagesOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

// Complete completes LoginOptions after they've been created
func (o *BuildImagesOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) (err error) {
	return nil
}

// Validate validates the LoginOptions based on completed values
func (o *BuildImagesOptions) Validate(ctx context.Context) (err error) {
	devfileObj := astracontext.GetEffectiveDevfileObj(ctx)
	if devfileObj == nil {
		return genericclioptions.NewNoDevfileError(astracontext.GetWorkingDirectory(ctx))
	}
	return nil
}

// Run contains the logic for the astra command
func (o *BuildImagesOptions) Run(ctx context.Context) (err error) {
	return image.BuildPushImages(ctx, image.SelectBackend(ctx), o.clientset.FS, o.pushFlag)
}

// NewCmdBuildImages implements the astra command
func NewCmdBuildImages(name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	o := NewBuildImagesOptions()
	buildImagesCmd := &cobra.Command{
		Use:     name,
		Short:   "Build images",
		Long:    "Build images defined in the devfile",
		Example: fmt.Sprintf(buildImagesExample, fullName),
		Args:    cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
	}

	util.SetCommandGroup(buildImagesCmd, util.MainGroup)
	buildImagesCmd.SetUsageTemplate(astrautil.CmdUsageTemplate)
	commonflags.UseVariablesFlags(buildImagesCmd)
	buildImagesCmd.Flags().BoolVar(&o.pushFlag, "push", false, "If true, build and push the images")
	clientset.Add(buildImagesCmd, clientset.FILESYSTEM)

	return buildImagesCmd
}
