package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/astra/util"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	RecommendedCommandName = "completion"
)

var (
	completionExample = ktemplates.Examples(`   # BASH

	## Load into your current shell environment
  source <(%[1]s bash)

	## Load persistently

	### Save the completion to a file
	%[1]s bash > ~/.astra/completion.bash.inc

	### Load the completion from within your $HOME/.bash_profile
	source ~/.astra/completion.bash.inc

  # ZSH

	## Load into your current shell environment
  source <(%[1]s zsh)

	## Load persistently
	%[1]s zsh > "${fpath[1]}/_astra"

	# FISH

	## Load into your current shell environment
	source <(%[1]s fish)

	## Load persistently
	%[1]s fish > ~/.config/fish/completions/astra.fish

	# POWERSHELL

	## Load into your current shell environment
	%[1]s powershell | Out-String | Invoke-Expression

	## Load persistently
	%[1]s powershell >> $PROFILE
`)
	completionLongDesc = ktemplates.LongDesc(`Add astra completion support to your development environment.
This will append your PS1 environment variable with astra component and application information.`)
)

// NewCmdCompletion implements the utils completion astra command
func NewCmdCompletion(name, fullName string) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:                   name,
		Short:                 "Add astra completion support to your development environment",
		Long:                  completionLongDesc,
		Example:               fmt.Sprintf(completionExample, fullName),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			// Below we ignore the error returns in order to pass golint validation
			// We will handle the error in the main function / output when the user inputs `astra completion`.
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				// Due to https://github.com/spf13/cobra/issues/1529 we cannot load zsh
				// via using source, so we need to add compdef to the beginning of the output so we can easily do:
				// source <(astra completion zsh)
				zsh := "#compdef astra\ncompdef _astra astra\n"
				out := os.Stdout
				_, _ = out.Write([]byte(zsh))
				_ = cmd.Root().GenZshCompletion(out)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	completionCmd.SetUsageTemplate(util.CmdUsageTemplate)
	util.SetCommandGroup(completionCmd, util.UtilityGroup)
	return completionCmd
}
