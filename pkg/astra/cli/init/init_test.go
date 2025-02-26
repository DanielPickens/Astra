package init

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	_init "github\.com/danielpickens/astra/pkg/init"
	"github\.com/danielpickens/astra/pkg/astra/cmdline"
	"github\.com/danielpickens/astra/pkg/astra/genericclioptions/clientset"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
)

func TestInitOptions_Complete(t *testing.T) {
	tests := []struct {
		name           string
		cmdlineExpects func(*cmdline.MockCmdline)
		initExpects    func(*_init.MockClient)
		fsysPopulate   func(fsys filesystem.Filesystem)
		wantErr        bool
	}{
		{
			name: "directory not empty",
			cmdlineExpects: func(mock *cmdline.MockCmdline) {
				mock.EXPECT().Context().Return(context.Background()).AnyTimes()
				mock.EXPECT().GetFlags().Times(1)
			},
			fsysPopulate: func(fsys filesystem.Filesystem) {
				_ = fsys.WriteFile(".emptyfile", []byte(""), 0644)
			},
			wantErr: false,
		},
		{
			name: "directory empty",
			cmdlineExpects: func(mock *cmdline.MockCmdline) {
				mock.EXPECT().Context().Return(context.Background()).AnyTimes()
				mock.EXPECT().GetFlags().Times(1)
			},
			initExpects: func(mock *_init.MockClient) {
			},
			fsysPopulate: func(fsys filesystem.Filesystem) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := filesystem.NewFakeFs()
			if tt.fsysPopulate != nil {
				tt.fsysPopulate(fsys)
			}
			ctrl := gomock.NewController(t)
			prefClient := preference.NewMockClient(ctrl)
			initClient := _init.NewMockClient(ctrl)
			initClient.EXPECT().GetFlags(gomock.Any()).Return(map[string]string{})
			o := NewInitOptions()
			o.SetClientset(&clientset.Clientset{
				PreferenceClient: prefClient,
				InitClient:       initClient,
				FS:               fsys,
			})
			cmdline := cmdline.NewMockCmdline(ctrl)
			if tt.cmdlineExpects != nil {
				tt.cmdlineExpects(cmdline)
			}
			if tt.initExpects != nil {
				tt.initExpects(initClient)
			}
			if err := o.Complete(context.Tastra(), cmdline, []string{}); (err != nil) != tt.wantErr {
				t.Errorf("InitOptions.Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
