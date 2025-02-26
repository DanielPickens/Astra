package genericclioptions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/devfile/library/v2/pkg/devfile/parser"
	"gopkg.in/AlecAivazis/survey.v1/terminal"

	"github\.com/danielpickens/astra/pkg/machineoutput"
	"github\.com/danielpickens/astra/pkg/version"

	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/commonflags"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	commonutil "github\.com/danielpickens/astra/pkg/util"

	"gopkg.in/AlecAivazis/survey.v1"

	"k8s.io/klog"
	"k8s.io/utils/pointer"

	"github\.com/danielpickens/astra/pkg/astra/cli/feature"

	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	fcontext "github\.com/danielpickens/astra/pkg/astra/commonflags/context"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/segment"
	scontext "github\.com/danielpickens/astra/pkg/segment/context"

	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/log"
)

type Runnable interface {
	SetClientset(clientset *clientset.Clientset)
	Complete(ctx context.Context, cmdline cmdline.Cmdline, args []string) error
	Validate(ctx context.Context) error
	Run(ctx context.Context) error
}

type SignalHandler interface {
	HandleSignal(ctx context.Context, cancelFunc context.CancelFunc) error
}

type Cleanuper interface {
	Cleanup(ctx context.Context, commandError error) error
}

// A PreIniter command is a command that will run `init` command if no file is present in current directory
// Commands implementing this interfaec must add FILESYSTEM and INIT dependencies
type PreIniter interface {
	// PreInit indicates a command will run `init`, and display the message returned by the method
	PreInit() string
}

// JsonOutputter must be implemented by commands with JSON output
// For these commands, the `-o json` flag will be added
// when err is not nil, the text of the error will be returned in a `message` field on stderr with an exit status of 1
// when err is nil, the result of RunForJsonOutput will be returned in JSON format on stdout with an exit status of 0
type JsonOutputter interface {
	RunForJsonOutput(ctx context.Context) (result interface{}, err error)
}

// DevfileUser must be implemented by commands that use Devfile depending on arguments, or
// commands that depend on FS but do not use Devfile.
// If the interface is not implemented and the command depends on FS, the command is expected to use Devfile
type DevfileUser interface {
	// UseDevfile returns true if the command with the specified cmdline and args needs to have access to the Devfile
	UseDevfile(ctx context.Context, cmdline cmdline.Cmdline, args []string) bool
}

const (
	// defaultAppName is the default name of the application when an application name is not provided
	defaultAppName = "app"
)

func GenericRun(o Runnable, testClientset clientset.Clientset, cmd *cobra.Command, args []string) error {
	var (
		err             error
		startTime       = time.Now()
		ctx, cancelFunc = context.WithCancel(cmd.Context())
	)

	defer func() {
		if err != nil {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
		}
		cancelFunc()
	}()

	userConfig, _ := preference.NewClient(ctx)
	envConfig := envcontext.GetEnvConfig(ctx)

	//lint:ignore SA1019 We deprecated this env var, but until it is removed, we still need to support it
	disableTelemetryEnvSet := envConfig.astraDisableTelemetry != nil
	var disableTelemetry bool
	if disableTelemetryEnvSet {
		disableTelemetry = *envConfig.astraDisableTelemetry
	}
	debugTelemetry := pointer.StringDeref(envConfig.astraDebugTelemetryFile, "")
	trackingConsentValue, isTrackingConsentEnabled, trackingConsentEnvSet, trackingConsentErr := segment.IsTrackingConsentEnabled(&envConfig)

	// check for conflicting settings
	if trackingConsentErr == nil && disableTelemetryEnvSet && trackingConsentEnvSet && disableTelemetry == isTrackingConsentEnabled {
		//lint:ignore SA1019 We deprecated this env var, but we really want users to know there is a conflict here
		return fmt.Errorf("%[1]s and %[2]s values are in conflict. %[1]s is deprecated, please use only %[2]s",
			segment.DisableTelemetryEnv, segment.TrackingConsentEnv)
	}

	// Prompt the user to consent for telemetry if a value is not set already
	// Skip prompting if the preference command is called
	// This prompt has been placed here so that it does not prompt the user when they call --help
	if !userConfig.IsSet(preference.ConsentTelemetrySetting) && cmd.Parent().Name() != "preference" {
		if !segment.RunningInTerminal() {
			klog.V(4).Infof("Skipping telemetry question because there is no terminal (tty)\n")
		} else {
			var askConsent bool
			if trackingConsentErr != nil {
				klog.V(4).Infof("error in determining value of tracking consent env var: %v", trackingConsentErr)
				askConsent = true
			} else if trackingConsentEnvSet {
				if isTrackingConsentEnabled {
					klog.V(4).Infof("Skipping telemetry question due to %s=%s\n", segment.TrackingConsentEnv, trackingConsentValue)
					klog.V(4).Info("Telemetry is enabled!\n")
					if err1 := userConfig.SetConfiguration(preference.ConsentTelemetrySetting, "true"); err1 != nil {
						klog.V(4).Info(err1.Error())
					}
				} else {
					klog.V(4).Infof("Skipping telemetry question due to %s=%s\n", segment.TrackingConsentEnv, trackingConsentValue)
				}
			} else if disableTelemetry {
				//lint:ignore SA1019 We deprecated this env var, but until it is removed, we still need to support it
				klog.V(4).Infof("Skipping telemetry question due to %s=%t\n", segment.DisableTelemetryEnv, disableTelemetry)
			} else {
				askConsent = true
			}
			if askConsent {
				var consentTelemetry bool
				prompt := &survey.Confirm{Message: "Help astra improve by allowing it to collect usage data. Read about our privacy statement: https://developers.redhat.com/article/tool-data-collection. You can change your preference later by changing the ConsentTelemetry preference.", Default: true}
				err = survey.AskOne(prompt, &consentTelemetry, nil)
				if err != nil {
					return err
				}
				err = userConfig.SetConfiguration(preference.ConsentTelemetrySetting, strconv.FormatBool(consentTelemetry))
				if err != nil {
					klog.V(4).Info(err.Error())
				}
			}
		}
	}
	if len(debugTelemetry) > 0 {
		klog.V(4).Infof("WARNING: debug telemetry, if enabled, will be logged in %s", debugTelemetry)
	}

	isTelemetryEnabled := segment.IsTelemetryEnabled(userConfig, envcontext.GetEnvConfig(ctx))
	scontext.SetTelemetryStatus(ctx, isTelemetryEnabled)

	// We can dereference as there is a default value defined for this config field
	err = scontext.SetCaller(ctx, envConfig.TelemetryCaller)
	if err != nil {
		klog.V(3).Infof("error handling caller property for telemetry: %v", err)
		err = nil
	}

	scontext.SetFlags(ctx, cmd.Flags())
	// set value for telemetry status in context so that we do not need to call IsTelemetryEnabled every time to check its status
	scontext.SetPreviousTelemetryStatus(ctx, segment.IsTelemetryEnabled(userConfig, envConfig))

	scontext.SetExperimentalMode(ctx, envConfig.astraExperimentalMode)

	// Send data to telemetry in case of user interrupt
	captureSignals := []os.Signal{syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt}
	go commonutil.StartSignalWatcher(captureSignals, func(receivedSignal os.Signal) {
		err = fmt.Errorf("user interrupted the command execution: %w", terminal.InterruptErr)
		if handler, ok := o.(SignalHandler); ok {
			err = handler.HandleSignal(ctx, cancelFunc)
			if err != nil {
				log.Errorf("error handling interrupt signal : %v", err)
			}
		}
		scontext.SetSignal(ctx, receivedSignal)
		startTelemetry(cmd, err, startTime)
	})

	err = commonflags.CheckMachineReadableOutputCommand(&envConfig, cmd)
	if err != nil {
		return err
	}
	err = commonflags.CheckPlatformCommand(cmd)
	if err != nil {
		return err
	}
	err = commonflags.CheckVariablesCommand(cmd)
	if err != nil {
		return err
	}

	cmdLineObj := cmdline.NewCobra(cmd)
	platform := commonflags.GetPlatformValue(cmdLineObj)
	deps, err := clientset.Fetch(cmd, platform, testClientset)
	if err != nil {
		return err

	}
	o.SetClientset(deps)

	if feature.IsExperimentalModeEnabled(ctx) {
		log.DisplayExperimentalWarning()
	}

	ctx = fcontext.WithJsonOutput(ctx, commonflags.GetJsonOutputValue(cmdLineObj))
	if platform != "" {
		ctx = fcontext.WithPlatform(ctx, platform)
	}
	ctx = astracontext.WithApplication(ctx, defaultAppName)

	if deps.KubernetesClient != nil {
		namespace := deps.KubernetesClient.GetCurrentNamespace()
		ctx = astracontext.WithNamespace(ctx, namespace)
	}

	if deps.FS != nil {
		var cwd string
		cwd, err = deps.FS.Getwd()
		if err != nil {
			startTelemetry(cmd, err, startTime)
			return err
		}
		ctx = astracontext.WithWorkingDirectory(ctx, cwd)

		var variables map[string]string
		variables, err = commonflags.GetVariablesValues(cmdLineObj)
		if err != nil {
			return err

		}
		ctx = fcontext.WithVariables(ctx, variables)

		if preiniter, ok := o.(PreIniter); ok {
			msg := preiniter.PreInit()
			err = runPreInit(ctx, cwd, deps, cmdLineObj, msg)
			if err != nil {
				startTelemetry(cmd, err, startTime)
				return err
			}
		}

		useDevfile := true
		if devfileUser, ok := o.(DevfileUser); ok {
			useDevfile = devfileUser.UseDevfile(ctx, cmdLineObj, args)
		}

		if useDevfile {
			var devfilePath, componentName string
			var devfileObj *parser.DevfileObj
			devfilePath, devfileObj, componentName, err = getDevfileInfo(cmd, deps.FS, cwd, variables, userConfig.GetImageRegistry())
			if err != nil {
				startTelemetry(cmd, err, startTime)
				return err
			}
			ctx = astracontext.WithDevfilePath(ctx, devfilePath)
			ctx = astracontext.WithEffectiveDevfileObj(ctx, devfileObj)
			ctx = astracontext.WithComponentName(ctx, componentName)
		}
	}

	// Run completion, validation and run.
	// Only upload data to segment for completion and validation if a non-nil error is returned.
	err = o.Complete(ctx, cmdLineObj, args)
	if err != nil {
		startTelemetry(cmd, err, startTime)
		return err
	}

	err = o.Validate(ctx)
	if err != nil {
		startTelemetry(cmd, err, startTime)
		return err
	}

	if jsonOutputter, ok := o.(JsonOutputter); ok && log.IsJSON() {
		var out interface{}
		out, err = jsonOutputter.RunForJsonOutput(ctx)
		if err == nil {
			machineoutput.OutputSuccess(testClientset.Stdout, testClientset.Stderr, out)
		}
	} else {
		err = o.Run(ctx)
	}
	startTelemetry(cmd, err, startTime)
	if cleanuper, ok := o.(Cleanuper); ok {
		err = cleanuper.Cleanup(ctx, err)
	}
	return err
}

// startTelemetry uploads the data to segment if user has consented to usage data collection and the command is not telemetry
// Tastra: move this function to a more suitable place, preferably pkg/segment
func startTelemetry(cmd *cobra.Command, err error, startTime time.Time) {
	if strings.Contains(cmd.CommandPath(), "telemetry") {
		return
	}
	ctx := cmd.Context()
	// Re-read the preferences, so that we can get the real settings in case a command updated the preferences file
	currentUserConfig, prefErr := preference.NewClient(ctx)
	if prefErr != nil {
		klog.V(2).Infof("unable to build preferences client to get telemetry consent status: %v", prefErr)
		return
	}
	isTelemetryEnabled := segment.IsTelemetryEnabled(currentUserConfig, envcontext.GetEnvConfig(ctx))
	if !(scontext.GetPreviousTelemetryStatus(ctx) || isTelemetryEnabled) {
		return
	}
	scontext.SetTelemetryStatus(ctx, isTelemetryEnabled)
	uploadData := &segment.TelemetryData{
		Event: cmd.CommandPath(),
		Properties: segment.TelemetryProperties{
			Duration:      time.Since(startTime).Milliseconds(),
			Success:       err == nil,
			Tty:           segment.RunningInTerminal(),
			Version:       fmt.Sprintf("astra %v (%v)", version.VERSION, version.GITCOMMIT),
			CmdProperties: scontext.GetContextProperties(ctx),
		},
	}
	if err != nil {
		uploadData.Properties.Error = segment.SetError(err)
		uploadData.Properties.ErrorType = segment.ErrorType(err)
	}
	data, err1 := json.Marshal(uploadData)
	if err1 != nil {
		klog.V(4).Infof("Failed to marshall telemetry data. %q", err1.Error())
	}
	command := exec.Command(os.Args[0], "telemetry", string(data))
	if err1 = command.Start(); err1 != nil {
		klog.V(4).Infof("Failed to start the telemetry process. Error: %q", err1.Error())
		return
	}
	if err1 = command.Process.Release(); err1 != nil {
		klog.V(4).Infof("Failed to release the process. %q", err1.Error())
		return
	}
}

// NoArgsAndSilenceJSON returns the NoArgs value, and silence output when JSON output is activated
func NoArgsAndSilenceJSON(cmd *cobra.Command, args []string) error {
	if log.IsJSON() {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
	}
	return cobra.NoArgs(cmd, args)
}
