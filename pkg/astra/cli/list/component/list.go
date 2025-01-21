package component

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/kclient"
	"github\.com/danielpickens/astra/pkg/astra/cli/feature"
	"github\.com/danielpickens/astra/pkg/astra/cli/ui"
	"github\.com/danielpickens/astra/pkg/astra/commonflags"
	"github\.com/danielpickens/astra/pkg/astra/util"

	"github\.com/danielpickens/astra/pkg/component"

	"github\.com/danielpickens/astra/pkg/log"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	fcontext "github\.com/danielpickens/astra/pkg/astra/commonflags/context"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

// RecommendedCommandName is the recommended list name
const RecommendedCommandName = "component"

var listExample = ktemplates.Examples(`  # List all components in the application
%[1]s
  `)

// ListOptions ...
type ListOptions struct {
	// Clients
	clientset *clientset.Clientset

	// Local variables
	namespaceFilter string

	// Flags
	namespaceFlag string
}

var _ genericclioptions.Runnable = (*ListOptions)(nil)
var _ genericclioptions.JsonOutputter = (*ListOptions)(nil)

// NewListOptions ...
func NewListOptions() *ListOptions {
	return &ListOptions{}
}

func (o *ListOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

// Complete ...
func (lo *ListOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) (err error) {
	// If the namespace flag has been passed, we will search there.
	// if it hasn't, we will search from the default project / namespace.
	if lo.namespaceFlag != "" {
		if lo.clientset.KubernetesClient == nil {
			return kclient.NewNoConnectionError()
		}
		lo.namespaceFilter = lo.namespaceFlag
	} else if lo.clientset.KubernetesClient != nil {
		lo.namespaceFilter = astracontext.GetNamespace(ctx)
	}

	return nil
}

// Validate ...
func (lo *ListOptions) Validate(ctx context.Context) (err error) {
	if lo.clientset.KubernetesClient == nil {
		log.Warning(kclient.NewNoConnectionError())
	}
	return nil
}

// Run has the logic to perform the required actions as part of command
func (lo *ListOptions) Run(ctx context.Context) error {
	listSpinner := log.Spinnerf("Listing components from namespace '%s'", lo.namespaceFilter)
	defer listSpinner.End(false)

	list, err := lo.run(ctx)
	if err != nil {
		return err
	}

	listSpinner.End(true)

	HumanReadableOutput(ctx, list)
	return nil
}

// Run contains the logic for the astra command
func (lo *ListOptions) RunForJsonOutput(ctx context.Context) (out interface{}, err error) {
	return lo.run(ctx)
}

func (lo *ListOptions) run(ctx context.Context) (api.ResourcesList, error) {
	var (
		devfileObj    = astracontext.GetEffectiveDevfileObj(ctx)
		componentName = astracontext.GetComponentName(ctx)

		kubeClient   = lo.clientset.KubernetesClient
		podmanClient = lo.clientset.PodmanClient
	)

	switch fcontext.GetPlatform(ctx, "") {
	case commonflags.PlatformCluster:
		podmanClient = nil
	case commonflags.PlatformPodman:
		kubeClient = nil
	}

	allComponents, componentInDevfile, err := component.ListAllComponents(
		kubeClient, podmanClient, lo.namespaceFilter, devfileObj, componentName)
	if err != nil {
		return api.ResourcesList{}, err
	}

	// RunningOn is displayed only when Platform is active
	if !feature.IsEnabled(ctx, feature.GenericPlatformFlag) {
		for i := range allComponents {
			//lint:ignore SA1019 we need to output the deprecated value, before to remove it in a future release
			allComponents[i].RunningOn = ""
			allComponents[i].Platform = ""
		}
	}
	return api.ResourcesList{
		ComponentInDevfile: componentInDevfile,
		Components:         allComponents,
	}, nil
}

// NewCmdList implements the list astra command
func NewCmdComponentList(ctx context.Context, name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	o := NewListOptions()

	var listCmd = &cobra.Command{
		Use:     name,
		Short:   "List all components in the current namespace",
		Long:    "List all components in the current namespace.",
		Example: fmt.Sprintf(listExample, fullName),
		Args:    genericclioptions.NoArgsAndSilenceJSON,
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
		Aliases: []string{"components"},
	}
	clientset.Add(listCmd, clientset.KUBERNETES_NULLABLE, clientset.FILESYSTEM)
	if feature.IsEnabled(ctx, feature.GenericPlatformFlag) {
		clientset.Add(listCmd, clientset.PODMAN_NULLABLE)
	}
	listCmd.Flags().StringVar(&o.namespaceFlag, "namespace", "", "Namespace for astra to scan for components")

	util.SetCommandGroup(listCmd, util.ManagementGroup)
	commonflags.UseOutputFlag(listCmd)
	commonflags.UsePlatformFlag(listCmd)

	return listCmd
}

func HumanReadableOutput(ctx context.Context, list api.ResourcesList) {
	components := list.Components
	if len(components) == 0 {
		log.Error("There are no components deployed.")
		return
	}

	t := ui.NewTable()

	// Create the header and then sort accordingly
	headers := table.Row{"NAME", "PROJECT TYPE", "RUNNING IN", "MANAGED"}
	if feature.IsEnabled(ctx, feature.GenericPlatformFlag) {
		headers = append(headers, "PLATFORM")
	}
	t.AppendHeader(headers)
	t.SortBy([]table.SortBy{
		{Name: "MANAGED", Mode: table.Asc},
		{Name: "NAME", Mode: table.Dsc},
	})

	// Go through each component and add it to the table
	for _, comp := range components {

		// Mark the name as yellow in the index to it's easier to see.
		name := text.Colors{text.FgHiYellow}.Sprint(comp.Name)

		// Get the managed by label
		managedBy := comp.ManagedBy
		if managedBy == "" {
			managedBy = api.TypeUnknown
		}

		// Get the mode (dev or deploy)
		mode := comp.RunningIn.String()

		// Get the type of the component
		componentType := comp.Type
		if componentType == "" {
			componentType = api.TypeUnknown
		}

		// If we find our local unpushed component, let's change the output appropriately.
		if list.ComponentInDevfile == comp.Name {
			name = fmt.Sprintf("* %s", name)
		}
		if comp.ManagedByVersion != "" {
			managedBy += fmt.Sprintf(" (%s)", comp.ManagedByVersion)
		}
		// If we are managing that component, output it as blue (our logo colour) to indicate it's used by astra
		if comp.ManagedBy == "astra" {
			managedBy = text.Colors{text.FgBlue}.Sprintf(managedBy)
		}

		row := table.Row{name, componentType, mode, managedBy}

		if feature.IsEnabled(ctx, feature.GenericPlatformFlag) {
			platform := comp.Platform
			if platform == "" {
				platform = "None"
			}
			row = append(row, platform)
		}

		t.AppendRow(row)
	}
	t.Render()

}
