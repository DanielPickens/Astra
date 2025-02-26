package remove

import (
	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/astra/cli/remove/binding"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
)

// RecommendedCommandName is the recommended remove command name
const RecommendedCommandName = "remove"

// NewCmdRemove implements the astra remove command
func NewCmdRemove(name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	var removeCmd = &cobra.Command{
		Use:   name,
		Short: "Remove resources from devfile",
	}

	bindingCmd := binding.NewCmdBinding(binding.BindingRecommendedCommandName, util.GetFullName(fullName, binding.BindingRecommendedCommandName), testClientset)
	removeCmd.AddCommand(bindingCmd)
	util.SetCommandGroup(removeCmd, util.ManagementGroup)
	removeCmd.SetUsageTemplate(util.CmdUsageTemplate)

	return removeCmd
}
