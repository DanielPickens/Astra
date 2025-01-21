package preference

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/kclient"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/registry"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
)

func TestView(t *testing.T) {
	ctrl := gomock.NewController(t)
	prefClient := preference.NewMockClient(ctrl)
	kubeClient := kclient.NewMockClientInterface(ctrl)
	registryClient := registry.NewRegistryClient(filesystem.NewFakeFs(), prefClient, kubeClient)
	opts := NewViewOptions()
	opts.SetClientset(&clientset.Clientset{
		PreferenceClient: prefClient,
		RegistryClient:   registryClient,
	})

	cmdline := cmdline.NewMockCmdline(ctrl)

	args := []string{}
	err := opts.Complete(context.Tastra(), cmdline, args)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
		return
	}

	err = opts.Validate(context.Tastra())
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
		return
	}
	boolValue := true
	intValue := 5
	var intNilValue *int = nil
	var boolNilValue *bool = nil

	preferenceList := api.PreferenceList{
		Items: []api.PreferenceItem{
			{
				Name:    preference.UpdateNotificationSetting,
				Value:   boolNilValue,
				Default: false,
			},
			{
				Name:    preference.PushTimeoutSetting,
				Value:   &intValue,
				Default: preference.DefaultPushTimeout,
			},
			{
				Name:    preference.RegistryCacheTimeSetting,
				Value:   intNilValue,
				Default: preference.DefaultRegistryCacheTime,
			},
			{
				Name:    preference.ConsentTelemetrySetting,
				Value:   &boolValue,
				Default: preference.DefaultConsentTelemetrySetting,
			},
			{
				Name:    preference.TimeoutSetting,
				Value:   intNilValue,
				Default: preference.DefaultTimeout,
			},
			{
				Name:    preference.EphemeralSetting,
				Value:   &boolValue,
				Default: preference.DefaultEphemeralSetting,
			},
		},
	}
	registryList := []api.Registry{
		{
			Name:   preference.DefaultDevfileRegistryName,
			URL:    preference.DefaultDevfileRegistryURL,
			Secure: false,
		},
		{
			Name:   "StagingRegistry",
			URL:    "https://registry.staging.devfile.io",
			Secure: true,
		},
	}
	prefClient.EXPECT().NewPreferenceList().Return(preferenceList)
	prefClient.EXPECT().RegistryList().Return(registryList)
	// Tastra(rm3l): test with different data returned by GetDevfileRegistries
	kubeClient.EXPECT().GetRegistryList().Return(nil, nil)

	err = opts.Run(context.Background())
	if err != nil {
		t.Errorf(`Expected nil error, got %s`, err)
	}
}
