package init

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"

	"k8s.io/kubectl/pkg/util/templates"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/component"
	"github\.com/danielpickens/astra/pkg/devfile"
	"github\.com/danielpickens/astra/pkg/devfile/location"
	"github\.com/danielpickens/astra/pkg/init/backend"
	"github\.com/danielpickens/astra/pkg/libdevfile"
	"github\.com/danielpickens/astra/pkg/log"
	"github\.com/danielpickens/astra/pkg/astra/cli/files"
	"github\.com/danielpickens/astra/pkg/astra/cli/messages"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/commonflags"
	fcontext "github\.com/danielpickens/astra/pkg/astra/commonflags/context"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	astrautil "github\.com/danielpickens/astra/pkg/astra/util"
	scontext "github\.com/danielpickens/astra/pkg/segment/context"
)

// RecommendedCommandName is the recommended command name
const RecommendedCommandName = "init"

var initExample = templates.Examples(`
  # Boostrap a new component in interactive mode
  %[1]s

  # Bootstrap a new component with a specific devfile from registry;
  # if several versions of the devfile exists then it will download the default version 
  %[1]s --name my-app --devfile nodejs
  
  # Bootstrap a new component with a specific versioned devfile from registry
  %[1]s --name my-app --devfile nodejs --devfile-version 2.1.0

 # Bootstrap a new component with the latest devfile from registry
  %[1]s --name my-app --devfile nodejs --devfile-version latest

  # Bootstrap a new component with a specific devfile from a specific registry
  %[1]s --name my-app --devfile nodejs --devfile-registry MyRegistry
  
  # Bootstrap a new component with a specific devfile from the local filesystem
  %[1]s --name my-app --devfile-path $HOME/devfiles/nodejs/devfile.yaml
  
  # Bootstrap a new component with a specific devfile from the web
  %[1]s --name my-app --devfile-path https://devfiles.example.com/nodejs/devfile.yaml

  # Bootstrap a new component and download a starter project
  %[1]s --name my-app --devfile nodejs --starter nodejs-starter

  # Bootstrap a new component with a specific devfile from registry for a specific architecture
  %[1]s --name my-app --devfile nodejs --architecture s390x
  `)

type InitOptions struct {
	// Clients
	clientset *clientset.Clientset

	// Flags passed to the command
	flags map[string]string
}

var _ genericclioptions.Runnable = (*InitOptions)(nil)
var _ genericclioptions.JsonOutputter = (*InitOptions)(nil)

// NewInitOptions creates a new InitOptions instance
func NewInitOptions() *InitOptions {
	return &InitOptions{}
}

func (o *InitOptions) SetClientset(clientset *clientset.Clientset) {
	o.clientset = clientset
}

func (o *InitOptions) UseDevfile(ctx context.Context, cmdline cmdline.Cmdline, args []string) bool {
	return false
}

// Complete will build the parameters for init, using different backends based on the flags set,
// either by using flags or interactively if no flag is passed
// Complete will return an error immediately if the current working directory is not empty
func (o *InitOptions) Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) (err error) {

	o.flags = o.clientset.InitClient.GetFlags(cmdline.GetFlags())

	scontext.SetInteractive(cmdline.Context(), len(o.flags) == 0)

	return nil
}

// Validate validates the InitOptions based on completed values
func (o *InitOptions) Validate(ctx context.Context) error {

	workingDir := astracontext.GetWorkingDirectory(ctx)

	devfilePresent, err := location.DirectoryContainsDevfile(o.clientset.FS, workingDir)
	if err != nil {
		return err
	}
	if devfilePresent {
		return errors.New("a devfile already exists in the current directory")
	}

	err = o.clientset.InitClient.Validate(o.flags, o.clientset.FS, workingDir)
	if err != nil {
		return err
	}

	if len(o.flags) == 0 && fcontext.IsJsonOutput(ctx) {
		return errors.New("parameters are expected to select a devfile")
	}
	return nil
}

// Run contains the logic for the astra command
func (o *InitOptions) Run(ctx context.Context) (err error) {

	devfileObj, _, name, devfileLocation, starterInfo, err := o.run(ctx)
	if err != nil {
		return err
	}

	exitMessage := fmt.Sprintf(`
Your new component '%s' is ready in the current directory.
To start editing your component, use 'astra dev' and open this folder in your favorite IDE.
Changes will be directly reflected on the cluster.`, devfileObj.Data.GetMetadata().Name)

	if len(o.flags) == 0 {
		automateCommand := fmt.Sprintf("astra init --name %s --devfile %s --devfile-registry %s", name, devfileLocation.Devfile, devfileLocation.DevfileRegistry)
		if devfileLocation.DevfileVersion != "" {
			automateCommand = fmt.Sprintf("%s --devfile-version %s", automateCommand, devfileLocation.DevfileVersion)
		}
		if starterInfo != nil {
			automateCommand = fmt.Sprintf("%s --starter %s", automateCommand, starterInfo.Name)
		}

		klog.V(2).Infof("Port configuration using flag is currently not supported")
		log.Infof("\nYou can automate this command by executing:\n   %s", automateCommand)
	}

	if libdevfile.HasDeployCommand(devfileObj.Data) {
		exitMessage += "\nTo deploy your component to a cluster use \"astra deploy\"."
	}
	log.Info(exitMessage)

	return nil
}

// RunForJsonOutput is executed instead of Run when -o json flag is given
func (o *InitOptions) RunForJsonOutput(ctx context.Context) (out interface{}, err error) {
	devfileObj, devfilePath, _, _, _, err := o.run(ctx)
	if err != nil {
		return nil, err
	}
	devfileData, err := api.GetDevfileData(devfileObj)
	if err != nil {
		return nil, err
	}
	return api.Component{
		DevfilePath:       devfilePath,
		DevfileData:       devfileData,
		DevForwardedPorts: []api.ForwardedPort{},
		RunningIn:         api.NewRunningModes(),
		ManagedBy:         "astra",
	}, nil
}

// run downloads the devfile and starter project and returns the content of the devfile, path of the devfile, name of the component, api.DetectionResult object for DevfileRegistry info and StarterProject object
func (o *InitOptions) run(ctx context.Context) (devfileObj parser.DevfileObj, path string, name string, devfileLocation *api.DetectionResult, starterInfo *v1alpha2.StarterProject, err error) {
	var starterDownloaded bool

	workingDir := astracontext.GetWorkingDirectory(ctx)

	defer func() {
		if err == nil {
			return
		}
		if starterDownloaded {
			err = fmt.Errorf("%w\nthe command failed after downloading the starter project. By security, the directory is not cleaned up", err)
		} else {
			_ = o.clientset.FS.Remove("devfile.yaml")
			err = fmt.Errorf("%w\nthe command failed, the devfile has been removed from current directory", err)
		}
	}()

	isEmptyDir, err := location.DirIsEmpty(o.clientset.FS, workingDir)
	if err != nil {
		return parser.DevfileObj{}, "", "", nil, nil, err
	}

	// Show a welcome message for when you initially run `astra init`.

	var infoOutput string
	if isEmptyDir && len(o.flags) == 0 {
		infoOutput = messages.NoSourceCodeDetected
	} else if len(o.flags) == 0 {
		infoOutput = messages.SourceCodeDetected
	}
	log.Title(messages.InitializingNewComponent, infoOutput)
	log.Println()
	if len(o.flags) == 0 {
		log.Info(messages.InteractiveModeEnabled)
	}

	devfileObj, devfilePath, devfileLocation, err := o.clientset.InitClient.SelectAndPersonalizeDevfile(ctx, o.flags, workingDir)
	if err != nil {
		return parser.DevfileObj{}, "", "", nil, nil, err
	}

	starterInfo, err = o.clientset.InitClient.SelectStarterProject(devfileObj, o.flags, isEmptyDir)
	if err != nil {
		return parser.DevfileObj{}, "", "", nil, nil, err
	}

	// Set the name in the devfile but do not write it yet to disk,
	// because the starter project downloaded at the end might come bundled with a specific Devfile.
	name, err = o.clientset.InitClient.PersonalizeName(devfileObj, o.flags)
	if err != nil {
		return parser.DevfileObj{}, "", "", nil, nil, fmt.Errorf("failed to update the devfile's name: %w", err)
	}

	if starterInfo != nil {
		var containsDevfile bool
		// WARNING: this will remove all the content of the destination directory, ie the devfile.yaml file
		containsDevfile, err = o.clientset.InitClient.DownloadStarterProject(starterInfo, workingDir)
		if err != nil {
			return parser.DevfileObj{}, "", "", nil, nil, fmt.Errorf("unable to download starter project %q: %w", starterInfo.Name, err)
		}
		starterDownloaded = true

		// in case the starter project contains a devfile, read it again
		if containsDevfile {
			devfileObj, err = devfile.ParseAndValidateFromFile(devfilePath, "", false)
			if err != nil {
				return parser.DevfileObj{}, "", "", nil, nil, err
			}
		}
	}
	// WARNING: SetMetadataName writes the Devfile to disk
	if err = devfileObj.SetMetadataName(name); err != nil {
		return parser.DevfileObj{}, "", "", nil, nil, err
	}

	err = files.ReportLocalFileGeneratedByastra(o.clientset.FS, workingDir, filepath.Base(devfilePath))
	if err != nil {
		klog.V(4).Infof("error trying to report local file generated: %v", err)
	}

	scontext.SetComponentType(ctx, component.GetComponentTypeFromDevfileMetadata(devfileObj.Data.GetMetadata()))
	scontext.SetLanguage(ctx, devfileObj.Data.GetMetadata().Language)
	scontext.SetProjectType(ctx, devfileObj.Data.GetMetadata().ProjectType)
	scontext.SetDevfileName(ctx, devfileObj.GetMetadataName())

	return devfileObj, devfilePath, name, devfileLocation, starterInfo, nil
}

// NewCmdInit implements the astra command
func NewCmdInit(name, fullName string, testClientset clientset.Clientset) *cobra.Command {

	o := NewInitOptions()
	initCmd := &cobra.Command{
		Use:     name,
		Short:   "Init bootstraps a new project",
		Long:    "Bootstraps a new project",
		Example: fmt.Sprintf(initExample, fullName),
		Args:    cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return genericclioptions.GenericRun(o, testClientset, cmd, args)
		},
	}
	clientset.Add(initCmd, clientset.PREFERENCE, clientset.FILESYSTEM, clientset.REGISTRY, clientset.INIT)

	initCmd.Flags().String(backend.FLAG_NAME, "", "name of the component to create; it must follow the RFC 1123 Label Names standard and not be all-numeric")
	initCmd.Flags().String(backend.FLAG_DEVFILE, "", "name of the devfile in devfile registry")
	initCmd.Flags().String(backend.FLAG_DEVFILE_REGISTRY, "", "name of the devfile registry (as configured in \"astra preference view\"). It can be used in combination with --devfile, but not with --devfile-path")
	initCmd.Flags().String(backend.FLAG_STARTER, "", "name of the starter project")
	initCmd.Flags().String(backend.FLAG_DEVFILE_PATH, "", "path to a devfile. This is an alternative to using devfile from Devfile registry. It can be local filesystem path or http(s) URL")
	initCmd.Flags().String(backend.FLAG_DEVFILE_VERSION, "", "version of the devfile stack; use \"latest\" to dowload the latest stack")
	initCmd.Flags().StringArray(backend.FLAG_ARCHITECTURE, []string{}, "Architecture supported. Can be one or multiple values from amd64, arm64, ppc64le, s390x. Default is amd64.")
	initCmd.Flags().StringArray(backend.FLAG_RUN_PORT, []string{}, "ports used by the application (via the 'run' command)")

	commonflags.UseOutputFlag(initCmd)
	// Add a defined annotation in order to appear in the help menu
	util.SetCommandGroup(initCmd, util.MainGroup)
	initCmd.SetUsageTemplate(astrautil.CmdUsageTemplate)
	return initCmd
}
