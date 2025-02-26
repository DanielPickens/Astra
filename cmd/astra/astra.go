package main

import (
	"context"
	"flag"
	"os"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github\.com/danielpickens/astra/pkg/config"
	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	"github\.com/danielpickens/astra/pkg/log"
	"github\.com/danielpickens/astra/pkg/astra/cli"
	"github\.com/danielpickens/astra/pkg/astra/cli/version"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/astra/util"
	"github\.com/danielpickens/astra/pkg/astra/util/completion"
	"github\.com/danielpickens/astra/pkg/preference"
	segment "github\.com/danielpickens/astra/pkg/segment/context"

	"k8s.io/klog"
)

func main() {
	// We need to reinitialize this global variable in case flags are defined by third-party packages
	// (for example vendor/sigs.k8s.io/controller-runtime/pkg/client/config/config.go)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Create a context ready for receiving telemetry data
	// and save into it configuration based on environment variables
	ctx := segment.NewContext(context.Background())
	envConfig, err := config.GetConfiguration()
	if err != nil {
		util.LogErrorAndExit(err, "")
	}
	ctx = envcontext.WithEnvConfig(ctx, *envConfig)
	ctx = astracontext.WithPID(ctx, os.Getpid())

	// create the complete command
	klog.InitFlags(nil)

	root, err := cli.NewCmdastra(ctx, cli.astraRecommendedName, cli.astraRecommendedName, func(err error) {
		util.LogErrorAndExit(err, "")
	}, clientset.Clientset{})
	if err != nil {
		util.LogErrorAndExit(err, "")
	}
	rootCmp := createCompletion(root)
	cmp := complete.New("astra", rootCmp)

	// AddFlags adds the completion flags to the program flags, specifying custom names
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)

	// add the completion flags to the root command, though they won't appear in completions
	root.Flags().AddGoFlagSet(flag.CommandLine)

	// override usage so that flag.Parse uses root command's usage instead of default one when invoked with -h
	flag.CommandLine.Usage = func() {
		_ = root.Help()
	}

	// parse the flags but hack around to avoid exiting with error code 2 on help
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	args := os.Args[1:]
	if err = flag.CommandLine.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return
		}
	}

	// run the completion, in case that the completion was invoked
	// and ran as a completion script or handled a flag that passed
	// as argument, the Run method will return true,
	// in that case, our program have nothing to do and should return.
	if cmp.Complete() {
		return
	}

	cfg, err := preference.NewClient(ctx)
	if err != nil {
		util.LogErrorAndExit(err, "")
	}

	// Call commands
	// checking the value of updatenotification in config
	// before proceeding with fetching the latest version
	if cfg.GetUpdateNotification() {
		updateInfo := make(chan string)
		go version.GetLatestReleaseInfo(updateInfo)

		util.LogErrorAndExit(root.ExecuteContext(ctx), "")
		select {
		case message := <-updateInfo:
			log.Info(message)
		default:
			klog.V(4).Info("Could not get the latest release information in time. Never mind, exiting gracefully :)")
		}
	} else {
		util.LogErrorAndExit(root.ExecuteContext(ctx), "")
	}
}

func createCompletion(root *cobra.Command) complete.Command {
	rootCmp := complete.Command{}
	rootCmp.Flags = make(complete.Flags)
	addFlags := func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		var handler complete.Predictor
		handler, ok := completion.GetCommandFlagHandler(root, flag.Name)
		if !ok {
			handler = complete.PredictAnything
		}

		if len(flag.Shorthand) > 0 {
			rootCmp.Flags["-"+flag.Shorthand] = handler
		}

		rootCmp.Flags["--"+flag.Name] = handler
	}
	root.LocalFlags().VisitAll(addFlags)
	root.InheritedFlags().VisitAll(addFlags)
	if root.HasAvailableSubCommands() {
		rootCmp.Sub = make(complete.Commands)
		for _, c := range root.Commands() {
			if !c.Hidden {
				rootCmp.Sub[c.Name()] = createCompletion(c)
			}
		}
	}

	var handler complete.Predictor
	handler, ok := completion.GetCommandHandler(root)
	if !ok {
		handler = complete.PredictNothing
	}
	rootCmp.Args = handler

	return rootCmp
}
