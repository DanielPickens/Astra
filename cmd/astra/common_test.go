package main

import (
	"bytes"
	"context"
	"flag"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"k8s.io/klog"

	"github.com/sethvargo/go-envconfig"

	"github\.com/danielpickens/astra/pkg/config"
	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	"github\.com/danielpickens/astra/pkg/astra/cli"
	astracontext "github\.com/danielpickens/astra/pkg/astra/context"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
)

func resetGlobalFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	klog.InitFlags(nil)
}

type runOptions struct {
	env    map[string]string
	config map[string]string
}

func runCommand(
	t *testing.T,
	args []string,
	options runOptions,
	clientset clientset.Clientset,
	populateFS func(fs filesystem.Filesystem) error,
	f func(err error, stdout, stderr string),
) {

	// We are running the test on a new and empty directory (on real filesystem)
	originWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(originWd)
	}()
	cwd := t.TempDir()
	err = os.Chdir(cwd)
	if err != nil {
		t.Fatal(err)
	}

	if populateFS != nil {
		err = populateFS(clientset.FS)
		if err != nil {
			t.Fatal(err)
		}
	}

	ctx := context.Background()
	envConfig, err := config.GetConfigurationWith(envconfig.MapLookuper(options.config))

	if err != nil {
		t.Fatal(err)
	}
	ctx = envcontext.WithEnvConfig(ctx, *envConfig)
	ctx = astracontext.WithPID(ctx, 101)

	for k, v := range options.env {
		t.Setenv(k, v)
	}

	resetGlobalFlags()

	var stdoutB, stderrB bytes.Buffer

	clientset.Stdout = &stdoutB
	clientset.Stderr = &stderrB
	root, err := cli.NewCmdastra(ctx, cli.astraRecommendedName, cli.astraRecommendedName, nil, clientset)
	if err != nil {
		t.Fatal(err)
	}

	root.SetOut(&stdoutB)
	root.SetErr(&stderrB)

	root.SetArgs(args)

	err = root.ExecuteContext(ctx)

	stdout := stdoutB.String()
	stderr := stderrB.String()

	f(err, stdout, stderr)
}

func checkEqual[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Errorf("Name should be \"%v\" but is \"%v\"", b, a)
	}
}

func checkerr[T comparable](t *testing.T, err error, expected T) {
	if err != nil {
		t.Errorf("Error should be nil but is \"%v\"", err)
	}
}

