package validate

import (
	"testing"

	"github.com/devfile/library/v2/pkg/devfile/parser/data"

	devfilev1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	devfileParser "github.com/devfile/library/v2/pkg/devfile/parser"
	parsercommon "github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"github.com/devfile/library/v2/pkg/testingutil"

	"github\.com/danielpickens/astra/pkg/util"
)

func Test_getCommands(t *testing.T) {

	component := []devfilev1.Component{
		testingutil.GetFakeContainerComponent("alias1"),
	}

	tests := []struct {
		name             string
		execCommands     []devfilev1.Command
		compCommands     []devfilev1.Command
		expectedCommands []devfilev1.Command
	}{
		{
			name: "Case 1: One command",
			execCommands: []devfilev1.Command{
				{
					Id: "somecommand",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
			},
			expectedCommands: []devfilev1.Command{
				{
					Id: "somecommand",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
			},
		},
		{
			name: "Case 2: Multiple commands",
			execCommands: []devfilev1.Command{
				{
					Id: "somecommand",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
				{
					Id: "somecommand2",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
			},
			compCommands: []devfilev1.Command{
				{
					Id: "mycomposite",
					CommandUnion: devfilev1.CommandUnion{
						Composite: &devfilev1.CompositeCommand{
							Commands: []string{},
						},
					},
				},
			},
			expectedCommands: []devfilev1.Command{
				{
					Id: "somecommand",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
				{
					Id: "somecommand2",
					CommandUnion: devfilev1.CommandUnion{
						Exec: &devfilev1.ExecCommand{
							HotReloadCapable: util.GetBool(false),
						},
					},
				},
				{
					Id: "mycomposite",
					CommandUnion: devfilev1.CommandUnion{
						Composite: &devfilev1.CompositeCommand{
							Commands: []string{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devObj := devfileParser.DevfileObj{
				Data: func() data.DevfileData {
					devfileData, err := data.NewDevfileData(string(data.APISchemaVersion200))
					if err != nil {
						t.Error(err)
					}
					err = devfileData.AddComponents(component)
					if err != nil {
						t.Error(err)
					}
					err = devfileData.AddCommands(tt.execCommands)
					if err != nil {
						t.Error(err)
					}
					err = devfileData.AddCommands(tt.compCommands)
					if err != nil {
						t.Error(err)
					}
					return devfileData
				}(),
			}

			commands, err := devObj.Data.GetCommands(parsercommon.DevfileOptions{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			commandsMap := getCommandsMap(commands)
			if len(commandsMap) != len(tt.expectedCommands) {
				t.Errorf("TestGetCommands error: number of returned commands don't match: %v got: %v", len(tt.expectedCommands), len(commandsMap))
			}
			for _, command := range tt.expectedCommands {
				_, ok := commandsMap[command.Id]
				if !ok {
					t.Errorf("TestGetCommands error: command %v not found in map", command.Id)
				}
			}
		})
	}

}
