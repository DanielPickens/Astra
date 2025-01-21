package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	DockerCmd                     string        `env:"DOCKER_CMD,default=docker"`
	Globalastraconfig               *string       `env:"GLOBALastraCONFIG,noinit"`
	astraDebugTelemetryFile         *string       `env:"astra_DEBUG_TELEMETRY_FILE,noinit"`
	astraDisableTelemetry           *bool         `env:"astra_DISABLE_TELEMETRY,noinit"`
	astraLogLevel                   *int          `env:"astra_LOG_LEVEL,noinit"`
	astraTrackingConsent            *string       `env:"astra_TRACKING_CONSENT,noinit"`
	PodmanCmd                     string        `env:"PODMAN_CMD,default=podman"`
	PodmanCmdInitTimeout          time.Duration `env:"PODMAN_CMD_INIT_TIMEOUT,default=1s"`
	TelemetryCaller               string        `env:"TELEMETRY_CALLER,default="`
	astraExperimentalMode           bool          `env:"astra_EXPERIMENTAL_MODE,default=false"`
	PushImages                    bool          `env:"astra_PUSH_IMAGES,default=true"`
	astraContainerBackendGlobalArgs []string      `env:"astra_CONTAINER_BACKEND_GLOBAL_ARGS,noinit,delimiter=;"`
	astraImageBuildArgs             []string      `env:"astra_IMAGE_BUILD_ARGS,noinit,delimiter=;"`
	astraContainerRunArgs           []string      `env:"astra_CONTAINER_RUN_ARGS,noinit,delimiter=;"`
}

// GetConfiguration initializes a Configuration for astra by using the system environment.
// See GetConfigurationWith for a more configurable version.
func GetConfiguration() (*Configuration, error) {
	return GetConfigurationWith(envconfig.OsLookuper())
}

// GetConfigurationWith initializes a Configuration for astra by using the specified envconfig.Lookuper to resolve values.
// It is recommended to use this function (instead of GetConfiguration) if you don't need to depend on the current system environment,
// typically in unit tests.
func GetConfigurationWith(lookuper envconfig.Lookuper) (*Configuration, error) {
	var s Configuration
	c := envconfig.Config{
		Target:   &s,
		Lookuper: lookuper,
	}
	err := envconfig.ProcessWith(context.Background(), &c)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
