package create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/astra/cli/create/namespace"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	astrautil "github\.com/danielpickens/astra/pkg/astra/util"
)

// RecommendedCommandName is the recommended namespace command name
const RecommendedCommandName = "create"

// NewCmdCreate implements the namespace astra command
func NewCmdCreate(name, fullName string, testClientset clientset.Clientset) *cobra.Command {

	namespaceCreateCmd := namespace.NewCmdNamespaceCreate(namespace.RecommendedCommandName, astrautil.GetFullName(fullName, namespace.RecommendedCommandName), testClientset)
	createCmd := &cobra.Command{
		Use:   name + " [options]",
		Short: "Perform create operation",
		Long:  "Perform create operation",
		Example: fmt.Sprintf("%s\n",
			namespaceCreateCmd.Example,
		),
	}

	createCmd.AddCommand(namespaceCreateCmd)

	// Add a defined annotation in order to appear in the help menu
	util.SetCommandGroup(createCmd, util.ManagementGroup)
	createCmd.SetUsageTemplate(astrautil.CmdUsageTemplate)

	return createCmd
}
