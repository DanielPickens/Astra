package set

import (
	"fmt"

	"github\.com/danielpickens/astra/pkg/astra/cli/set/namespace"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"

	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended namespace command name
const RecommendedCommandName = "set"

// NewCmdSet implements the namespace astra command
func NewCmdSet(name, fullName string, testClientset clientset.Clientset) *cobra.Command {

	namespaceSetCmd := namespace.NewCmdNamespaceSet(namespace.RecommendedCommandName,
		util.GetFullName(fullName, namespace.RecommendedCommandName), testClientset)
	setCmd := &cobra.Command{
		Use:   name + " [options]",
		Short: "Perform set operation",
		Long:  "Perform set operation",
		Example: fmt.Sprintf("%s\n",
			namespaceSetCmd.Example,
		),
	}

	setCmd.AddCommand(namespaceSetCmd)

	util.SetCommandGroup(setCmd, util.ManagementGroup)
	setCmd.SetUsageTemplate(util.CmdUsageTemplate)

	return setCmd
}
