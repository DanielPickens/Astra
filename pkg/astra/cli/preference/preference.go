package preference

import (
	"context"
	"fmt"

	"github\.com/danielpickens/astra/pkg/astra/cli/preference/add"
	"github\.com/danielpickens/astra/pkg/astra/cli/preference/remove"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	"github\.com/danielpickens/astra/pkg/preference"

	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

// RecommendedCommandName is the recommended preference command name
const RecommendedCommandName = "preference"

var preferenceLongDesc = ktemplates.LongDesc(`Modifies astra specific configuration settings within the global preference file.

%[1]s`)

// NewCmdPreference implements the utils config astra command
func NewCmdPreference(ctx context.Context, name, fullName string, testClientset clientset.Clientset) *cobra.Command {

	// Main Commands
	preferenceViewCmd := NewCmdView(viewCommandName, util.GetFullName(fullName, viewCommandName), testClientset)
	preferenceSetCmd := NewCmdSet(ctx, setCommandName, util.GetFullName(fullName, setCommandName), testClientset)
	preferenceUnsetCmd := NewCmdUnset(unsetCommandName, util.GetFullName(fullName, unsetCommandName), testClientset)
	preferenceAddCmd := add.NewCmdAdd(add.RecommendedCommandName, util.GetFullName(fullName, add.RecommendedCommandName), testClientset)
	preferenceRemoveCmd := remove.NewCmdRemove(remove.RecommendedCommandName, util.GetFullName(fullName, remove.RecommendedCommandName), testClientset)

	// Subcommands

	// Set the examples
	preferenceCmd := &cobra.Command{
		Use:   name,
		Short: "Modifies preference settings",
		Long:  fmt.Sprintf(preferenceLongDesc, preference.FormatSupportedParameters()),
		Example: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
			preferenceViewCmd.Example,
			preferenceSetCmd.Example,
			preferenceUnsetCmd.Example,
			preferenceAddCmd.Example,
			preferenceRemoveCmd.Example,
		),
	}

	// Add the commands, help, usage and annotations
	preferenceCmd.AddCommand(preferenceViewCmd, preferenceSetCmd, preferenceUnsetCmd, preferenceAddCmd, preferenceRemoveCmd)
	preferenceCmd.SetUsageTemplate(util.CmdUsageTemplate)
	util.SetCommandGroup(preferenceCmd, util.UtilityGroup)

	return preferenceCmd
}
