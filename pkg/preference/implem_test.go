package preference

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github\.com/danielpickens/astra/pkg/config"
	envcontext "github\.com/danielpickens/astra/pkg/config/context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNew(t *testing.T) {

	tempConfigFile, err := os.CreateTemp("", "astraconfig")
	if err != nil {
		t.Fatal(err)
	}
	defer tempConfigFile.Close()
	tempConfigFileName := tempConfigFile.Name()

	tests := []struct {
		name    string
		output  *preferenceInfo
		success bool
	}{
		{
			name: "Test filename is being set",
			output: &preferenceInfo{
				Filename: tempConfigFile.Name(),
				Preference: Preference{
					TypeMeta: metav1.TypeMeta{
						Kind:       preferenceKind,
						APIVersion: preferenceAPIVersion,
					},
					astraSettings: astraSettings{
						RegistryList: &[]Registry{
							{
								Name:   DefaultDevfileRegistryName,
								URL:    DefaultDevfileRegistryURL,
								Secure: false,
							},
						},
					},
				},
			},
			success: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = envcontext.WithEnvConfig(ctx, config.Configuration{
				Globalastraconfig: &tempConfigFileName,
			})
			cfi, err := newPreferenceInfo(ctx)
			switch test.success {
			case true:
				if err != nil {
					t.Errorf("expected test to pass, but it failed with error: %v", err)
				}
			case false:
				if err == nil {
					t.Errorf("expected test to fail, but it passed!")
				}
			}
			if diff := cmp.Diff(test.output, cfi); diff != "" {
				t.Errorf("newPreferenceInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetPushTimeout(t *testing.T) {

	nonzeroValue := 5 * time.Second

	tests := []struct {
		name           string
		existingConfig Preference
		want           time.Duration
	}{
		{
			name:           "Validating default value from test case",
			existingConfig: Preference{},
			want:           240,
		},
		{
			name: "Validating value 5 from configuration",
			existingConfig: Preference{
				astraSettings: astraSettings{
					PushTimeout: &nonzeroValue,
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = envcontext.WithEnvConfig(ctx, config.Configuration{})
			cfg, err := newPreferenceInfo(ctx)
			if err != nil {
				t.Error(err)
			}
			cfg.Preference = tt.existingConfig

			output := cfg.GetPushTimeout()
			if output != (tt.want * time.Second) {
				t.Errorf("GetPushTimeout returned unexpected value\ngot: %d \nexpected: %d\n", output, tt.want)
			}
		})
	}
}

func TestGetTimeout(t *testing.T) {
	zeroValue := 0 * time.Second
	nonzeroValue := 5 * time.Second
	tests := []struct {
		name           string
		existingConfig Preference
		want           time.Duration
	}{
		{
			name:           "validating value 1 from config in default case",
			existingConfig: Preference{},
			want:           1,
		},

		{
			name: "validating value 0 from config",
			existingConfig: Preference{
				astraSettings: astraSettings{
					Timeout: &zeroValue,
				},
			},
			want: 0,
		},

		{
			name: "validating value 5 from config",
			existingConfig: Preference{
				astraSettings: astraSettings{
					Timeout: &nonzeroValue,
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = envcontext.WithEnvConfig(ctx, config.Configuration{})
			cfg, err := newPreferenceInfo(ctx)
			if err != nil {
				t.Error(err)
			}
			cfg.Preference = tt.existingConfig

			output := cfg.GetTimeout()
			if output != (tt.want * time.Second) {
				t.Errorf("GetTimeout returned unexpected value\ngot: %d \nexpected: %d\n", output, tt.want)
			}
		})
	}
}

func TestSetConfiguration(t *testing.T) {
	trueValue := true
	falseValue := false
	minValue := minimumDurationValue

	tests := []struct {
		name           string
		parameter      string
		value          string
		existingConfig Preference
		wantErr        bool
		want           interface{}
	}{
		// update notification
		{
			name:           fmt.Sprintf("%s set nil to true", UpdateNotificationSetting),
			parameter:      UpdateNotificationSetting,
			value:          "true",
			existingConfig: Preference{},
			want:           true,
			wantErr:        false,
		},
		{
			name:      fmt.Sprintf("%s set true to false", UpdateNotificationSetting),
			parameter: UpdateNotificationSetting,
			value:     "false",
			existingConfig: Preference{
				astraSettings: astraSettings{
					UpdateNotification: &trueValue,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      fmt.Sprintf("%s set false to true", UpdateNotificationSetting),
			parameter: UpdateNotificationSetting,
			value:     "true",
			existingConfig: Preference{
				astraSettings: astraSettings{
					UpdateNotification: &falseValue,
				},
			},
			want:    true,
			wantErr: false,
		},

		{
			name:           fmt.Sprintf("%s invalid value", UpdateNotificationSetting),
			parameter:      UpdateNotificationSetting,
			value:          "invalid_value",
			existingConfig: Preference{},
			wantErr:        true,
		},
		// time out
		{
			name:      fmt.Sprintf("%s set to 5 from 0", TimeoutSetting),
			parameter: TimeoutSetting,
			value:     "5s",
			existingConfig: Preference{
				astraSettings: astraSettings{
					Timeout: &minValue,
				},
			},
			want:    5 * time.Second,
			wantErr: false,
		},
		{
			name:           fmt.Sprintf("%s set to 300", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "5m",
			existingConfig: Preference{},
			want:           5 * time.Second,
			wantErr:        false,
		},
		{
			name:           fmt.Sprintf("%s set to 0", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "0s",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s invalid value", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "this",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s set to 300 with mixed case in parameter name", TimeoutSetting),
			parameter:      "TimeOut",
			value:          "5m",
			existingConfig: Preference{},
			want:           5 * time.Minute,
			wantErr:        false,
		},
		// invalid parameter
		{
			name:           "invalid parameter",
			parameter:      "invalid_parameter",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s set to 0", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "0s",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s invalid value", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "invalid",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s negative value", TimeoutSetting),
			parameter:      TimeoutSetting,
			value:          "-5s",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("%s set to 99 with mixed case in parameter name", TimeoutSetting),
			parameter:      "PushTimeout",
			value:          "99s",
			existingConfig: Preference{},
			want:           99 * time.Second,
			wantErr:        false,
		},
		{
			name:           "set RegistryCacheTime to 1 minutes",
			parameter:      "RegistryCacheTime",
			value:          "1m",
			existingConfig: Preference{},
			want:           1 * time.Minute,
			wantErr:        false,
		},
		{
			name:           "set RegistryCacheTime to non int value",
			parameter:      "RegistryCacheTime",
			value:          "a",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("set %s to non bool value", ConsentTelemetrySetting),
			parameter:      ConsentTelemetrySetting,
			value:          "123",
			existingConfig: Preference{},
			wantErr:        true,
		},
		{
			name:           fmt.Sprintf("set %s from nil to true", ConsentTelemetrySetting),
			parameter:      ConsentTelemetrySetting,
			value:          "true",
			existingConfig: Preference{},
			wantErr:        false,
			want:           true,
		},
		{
			name:      fmt.Sprintf("set %s from true to false", ConsentTelemetrySetting),
			parameter: ConsentTelemetrySetting,
			value:     "false",
			existingConfig: Preference{
				astraSettings: astraSettings{
					ConsentTelemetry: &trueValue,
				},
			},
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = envcontext.WithEnvConfig(ctx, config.Configuration{})
			cfg, err := newPreferenceInfo(ctx)
			if err != nil {
				t.Error(err)
			}
			cfg.Preference = tt.existingConfig

			err = cfg.SetConfiguration(tt.parameter, tt.value)

			if !tt.wantErr && err == nil {
				// validating the value after executing Serconfiguration
				// according to component in positive cases
				switch tt.parameter {
				case "updatenotification":
					if *cfg.astraSettings.UpdateNotification != tt.want {
						t.Errorf("unexpected value after execution of SetConfiguration\ngot: %t \nexpected: %t\n", *cfg.astraSettings.UpdateNotification, tt.want)
					}
				case "timeout":
					if *cfg.astraSettings.Timeout != tt.want {
						t.Errorf("unexpected value after execution of SetConfiguration\ngot: %v \nexpected: %d\n", cfg.astraSettings.Timeout, tt.want)
					}
				case "registrycachetime":
					if *cfg.astraSettings.RegistryCacheTime != tt.want {
						t.Errorf("unexpected value after execution of SetConfiguration\ngot: %v \nexpected: %d\n", *cfg.astraSettings.RegistryCacheTime, tt.want)
					}
				}
			} else if tt.wantErr && err != nil {
				// negative cases
				switch tt.parameter {
				case "updatenotification":
				case "timeout":
					typedval, e := strconv.Atoi(tt.value)
					// if err is found in cases other than value <0 or !ok
					if !(typedval < 0 || e != nil) {
						t.Error(e)
					}
				}
			} else {
				t.Error(err)
			}

		})
	}
}

func TestConsentTelemetry(t *testing.T) {
	trueValue := true
	falseValue := false

	tests := []struct {
		name           string
		existingConfig Preference
		want           bool
	}{
		{
			name:           fmt.Sprintf("%s nil", ConsentTelemetrySetting),
			existingConfig: Preference{},
			want:           false,
		},
		{
			name: fmt.Sprintf("%s true", ConsentTelemetrySetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					ConsentTelemetry: &trueValue,
				},
			},
			want: true,
		},
		{
			name: fmt.Sprintf("%s false", ConsentTelemetrySetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					ConsentTelemetry: &falseValue,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := preferenceInfo{
				Preference: tt.existingConfig,
			}
			output := cfg.GetConsentTelemetry()

			if output != tt.want {
				t.Errorf("ConsentTelemetry returned unexpected value\ngot: %t \nexpected: %t\n", output, tt.want)
			}

		})
	}
}

func TestGetupdateNotification(t *testing.T) {
	trueValue := true
	falseValue := false

	tests := []struct {
		name           string
		existingConfig Preference
		want           bool
	}{
		{
			name:           fmt.Sprintf("%s nil", UpdateNotificationSetting),
			existingConfig: Preference{},
			want:           true,
		},
		{
			name: fmt.Sprintf("%s true", UpdateNotificationSetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					UpdateNotification: &trueValue,
				},
			},
			want: true,
		},
		{
			name: fmt.Sprintf("%s false", UpdateNotificationSetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					UpdateNotification: &falseValue,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := preferenceInfo{
				Preference: tt.existingConfig,
			}
			output := cfg.GetUpdateNotification()

			if output != tt.want {
				t.Errorf("GetUpdateNotification returned unexpected value\ngot: %t \nexpected: %t\n", output, tt.want)
			}

		})
	}
}

func TestIsSupportedParameter(t *testing.T) {
	tests := []struct {
		testName      string
		param         string
		expectedLower string
		expected      bool
	}{
		{
			testName:      "existing, lower case",
			param:         "timeout",
			expectedLower: "timeout",
			expected:      true,
		},
		{
			testName:      "existing, from description",
			param:         "Timeout",
			expectedLower: "timeout",
			expected:      true,
		},
		{
			testName:      "existing, mixed case",
			param:         "TimeOut",
			expectedLower: "timeout",
			expected:      true,
		},
		{
			testName: "empty",
			param:    "",
			expected: false,
		},
		{
			testName: "unexisting",
			param:    "foo",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Log("Running test: ", tt.testName)
		t.Run(tt.testName, func(t *testing.T) {
			actual, ok := asSupportedParameter(tt.param)
			if tt.expected != ok && tt.expectedLower != actual {
				t.Fail()
			}
		})
	}
}

func TestPreferenceIsntCreatedWhenastraIsUsed(t *testing.T) {
	// cleaning up old astra files if any
	ctx := context.Background()
	ctx = envcontext.WithEnvConfig(ctx, config.Configuration{})
	filename, err := getPreferenceFile(ctx)
	if err != nil {
		t.Error(err)
	}
	os.RemoveAll(filename)

	conf, err := newPreferenceInfo(ctx)
	if err != nil {
		t.Errorf("error while creating global preference %v", err)
	}
	if _, err = os.Stat(conf.Filename); !os.IsNotExist(err) {
		t.Errorf("preference file shouldn't exist yet")
	}
}

func TestMetaTypePopulatedInPreference(t *testing.T) {
	ctx := context.Background()
	ctx = envcontext.WithEnvConfig(ctx, config.Configuration{})
	pi, err := newPreferenceInfo(ctx)

	if err != nil {
		t.Error(err)
	}
	if pi.APIVersion != preferenceAPIVersion || pi.Kind != preferenceKind {
		t.Error("the api version and kind in preference are incorrect")
	}
}

func TestHandleWithoutRegistryExist(t *testing.T) {
	tests := []struct {
		name         string
		registryList []Registry
		operation    string
		registryName string
		registryURL  string
		want         []Registry
	}{
		{
			name:         "Add registry",
			registryList: []Registry{},
			operation:    "add",
			registryName: "testName",
			registryURL:  "testURL",
			want: []Registry{
				{
					Name: "testName",
					URL:  "testURL",
				},
			},
		},
		{
			name:         "Delete registry",
			registryList: []Registry{},
			operation:    "remove",
			registryName: "testName",
			registryURL:  "testURL",
			want:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleWithoutRegistryExist(tt.registryList, tt.operation, tt.registryName, tt.registryURL, false)
			if err != nil {
				t.Logf("Error message is %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("handleWithoutRegistryExist() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleWithRegistryExist(t *testing.T) {
	tests := []struct {
		name         string
		index        int
		registryList []Registry
		operation    string
		registryName string
		forceFlag    bool
		want         []Registry
	}{
		{
			name:  "Add registry",
			index: 0,
			registryList: []Registry{
				{
					Name: "testName",
					URL:  "testURL",
				},
			},
			operation:    "add",
			registryName: "testName",
			forceFlag:    false,
			want:         nil,
		},
		{
			name:  "Delete registry",
			index: 0,
			registryList: []Registry{
				{
					Name: "testName",
					URL:  "testURL",
				},
			},
			operation:    "remove",
			registryName: "testName",
			forceFlag:    true,
			want:         []Registry{},
		},
	}

	for _, tt := range tests {
		got, err := handleWithRegistryExist(tt.index, tt.registryList, tt.operation, tt.registryName, tt.forceFlag)
		if err != nil {
			t.Logf("Error message is %v", err)
		}

		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("handleWithRegistryExist() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGetConsentTelemetry(t *testing.T) {
	trueValue := true
	falseValue := false

	tests := []struct {
		name           string
		existingConfig Preference
		want           bool
	}{{
		name:           fmt.Sprintf("%s nil", ConsentTelemetrySetting),
		existingConfig: Preference{},
		want:           false,
	},
		{
			name: fmt.Sprintf("%s true", ConsentTelemetrySetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					ConsentTelemetry: &trueValue,
				},
			},
			want: true,
		},
		{
			name: fmt.Sprintf("%s false", ConsentTelemetrySetting),
			existingConfig: Preference{
				astraSettings: astraSettings{
					ConsentTelemetry: &falseValue,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := preferenceInfo{
				Preference: tt.existingConfig,
			}
			output := cfg.GetConsentTelemetry()

			if output != tt.want {
				t.Errorf("GetConsentTelemetry returned unexpected value\ngot: %t \nexpected: %t\n", output, tt.want)
			}

		})
	}
}
