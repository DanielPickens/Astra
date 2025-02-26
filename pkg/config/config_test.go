package config

import (
	"testing"

	"github.com/sethvargo/go-envconfig"
)

func TestDefaultValues(t *testing.T) {
	cfg, err := GetConfigurationWith(envconfig.MapLookuper(nil))
	if err != nil {
		t.Errorf("Error is not expected: %v", err)
	}

	checkDefaultStringValue(t, "DockerCmd", cfg.DockerCmd, "docker")
	checkDefaultStringValue(t, "PodmanCmd", cfg.PodmanCmd, "podman")
	checkDefaultStringValue(t, "TelemetryCaller", cfg.TelemetryCaller, "")
	checkDefaultBoolValue(t, "astraExperimentalMode", cfg.astraExperimentalMode, false)

	// Use noinit to set non initialized value as nil instead of zero-value
	checkNilString(t, "Globalastraconfig", cfg.Globalastraconfig)
	checkNilString(t, "Globalastraconfig", cfg.Globalastraconfig)
	checkNilString(t, "astraDebugTelemetryFile", cfg.astraDebugTelemetryFile)
	checkNilBool(t, "astraDisableTelemetry", cfg.astraDisableTelemetry)
	checkNilString(t, "astraTrackingConsent", cfg.astraTrackingConsent)

}

func checkDefaultStringValue(t *testing.T, fieldName string, field string, def string) {
	if field != def {
		t.Errorf("default value for %q should be %q but is %q", fieldName, def, field)
	}

}

func checkDefaultBoolValue(t *testing.T, fieldName string, field bool, def bool) {
	if field != def {
		t.Errorf("default value for %q should be %v but is %v", fieldName, def, field)
	}

}

func checkNilString(t *testing.T, fieldName string, field *string) {
	if field != nil {
		t.Errorf("value for non specified env var %q should be nil but is %q", fieldName, *field)

	}
}

func checkNilBool(t *testing.T, fieldName string, field *bool) {
	if field != nil {
		t.Errorf("value for non specified env var %q should be nil but is %v", fieldName, *field)

	}
}
