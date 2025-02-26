package logout

import (
	"context"
	"fmt"
	"os"

	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	astrautil "github\.com/danielpickens/astra/pkg/astra/util"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

// RecommendedCommandName is the recommended command name
const RecommendedCommandName = "logout"

var example = templates.Examples(`  # Logout
  %[1]s
`)

// LogoutOptions encapsulates the options for the astra logout command
type LogoutOptions struct {
	// Clients
	clientset *clientset.Clientset
}

var _ genericclioptions.Runnable = (*LogoutOptions)(nil)

// NewLogoutOptions creates a new LogoutOptions instance
func NewLogoutOptions() *LogoutOptions {
	return &LogoutOptions{}
}

func (o *LogoutOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

// Complete completes LogoutOptions after they've been created
func (o *LogoutOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) (err error) {
	return nil
}

// Validate validates the LogoutOptions based on completed values
func (o *LogoutOptions) Validate(ctx context.Context) (err error) {
	return nil
}

// Run contains the logic for the astra logout command
func (o *LogoutOptions) Run(ctx context.Context) (err error) {
	return o.clientset.KubernetesClient.RunLogout(os.Stdout)
}

// NewCmdLogout implements the logout astra command
func NewCmdLogout(name, fullName string, testClientset clientset.Clientset) *cobra.Command {
	o := NewLogoutOptions()
	logoutCmd := &cobra.Command{
		Use:     name,
		Short:   "Logout of the cluster",
		Long:    "Logout of the cluster.",
		Example: fmt.Sprintf(example, fullName),
		Args:    genericclioptions.NoArgsAndSilenceJSON,
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
	}

	// Add a defined annotation in order to appear in the help menu
	util.SetCommandGroup(logoutCmd, util.OpenshiftGroup)
	logoutCmd.SetUsageTemplate(astrautil.CmdUsageTemplate)

	clientset.Add(logoutCmd, clientset.KUBERNETES)

	return logoutCmd
}
