package delete

import (
	"context"

	"github\.com/danielpickens/astra/pkg/astra/cli/delete/component"
	"github\.com/danielpickens/astra/pkg/astra/cli/delete/namespace"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended delete command name
const RecommendedCommandName = "delete"

// NewCmdDelete implements the delete astra command
func NewCmdDelete(ctx context.Context, name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   name,
		Short: "Delete resources",
	}

	componentCmd := component.NewCmdComponent(ctx, component.ComponentRecommendedCommandName,
		util.GetFullName(fullName, component.ComponentRecommendedCommandName), testClientset)
	deleteCmd.AddCommand(componentCmd)

	namespaceDeleteCmd := namespace.NewCmdNamespaceDelete(namespace.RecommendedCommandName,
		util.GetFullName(fullName, namespace.RecommendedCommandName), testClientset)
	deleteCmd.AddCommand(namespaceDeleteCmd)

	util.SetCommandGroup(deleteCmd, util.ManagementGroup)
	deleteCmd.SetUsageTemplate(util.CmdUsageTemplate)

	return deleteCmd
}
