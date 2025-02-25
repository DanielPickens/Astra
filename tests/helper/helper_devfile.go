package helper

import (
	"fmt"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/v2/pkg/attributes"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	parsercommon "github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	. "github.com/onsi/gomega"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/pointer"

	"github\.com/danielpickens/astra/pkg/devfile"
)

// DevfileUpdater is a helper type that can mutate a Devfile object.
// It is intended to be used in conjunction with the UpdateDevfileContent function.
type DevfileUpdater func(*parser.DevfileObj) error

// DevfileMetadataNameSetter sets the 'metadata.name' field into the given Devfile
var DevfileMetadataNameSetter = func(name string) DevfileUpdater {
	return func(d *parser.DevfileObj) error {
		return d.SetMetadataName(name)
	}
}

// DevfileMetadataNameRemover removes the 'metadata.name' field from the given Devfile
var DevfileMetadataNameRemover = DevfileMetadataNameSetter("")

// DevfileCommandGroupUpdater updates the group definition of the specified command.
// It returns an error if the command was not found in the Devfile, or if there are multiple commands with the same name and type.
var DevfileCommandGroupUpdater = func(cmdName string, cmdType v1alpha2.CommandType, group *v1alpha2.CommandGroup) DevfileUpdater {
	return func(d *parser.DevfileObj) error {
		cmds, err := d.Data.GetCommands(parsercommon.DevfileOptions{
			CommandOptions: parsercommon.CommandOptions{
				CommandType: cmdType,
			},
			FilterByName: cmdName,
		})
		if err != nil {
			return err
		}
		if len(cmds) != 1 {
			return fmt.Errorf("found %v command(s) with name %q", len(cmds), cmdName)
		}
		cmd := cmds[0]
		switch cmdType {
		case v1alpha2.ApplyCommandType:
			cmd.Apply.Group = group
		case v1alpha2.CompositeCommandType:
			cmd.Composite.Group = group
		case v1alpha2.CustomCommandType:
			cmd.Custom.Group = group
		case v1alpha2.ExecCommandType:
			cmd.Exec.Group = group
		default:
			return fmt.Errorf("command type not handled: %q", cmdType)
		}
		return nil
	}
}

// UpdateDevfileContent parses the Devfile at the given path, then updates its content using the given handlers, and writes the updated Devfile to the given path.
//
// The handlers are invoked in the order they are provided.
//
// No operation is performed if no handler function is specified.
//
// See DevfileMetadataNameRemover for an example of handler function that can operate on the Devfile content.
func UpdateDevfileContent(path string, handlers []DevfileUpdater) {
	if len(handlers) == 0 {
		//Nothing to do => skip
		return
	}

	d, err := parser.ParseDevfile(parser.ParserArgs{
		Path:               path,
		FlattenedDevfile:   pointer.Bool(false),
		SetBooleanDefaults: pointer.Bool(false),
	})
	Expect(err).NotTo(HaveOccurred())
	for _, h := range handlers {
		err = h(&d)
		Expect(err).NotTo(HaveOccurred())
	}
	err = d.WriteYamlDevfile()
	Expect(err).NotTo(HaveOccurred())
}

// SetFsGroup is a DevfileUpdater which sets an attribute to a container
// to set a specific fsGroup for the container's pod
func SetFsGroup(containerName string, fsGroup int) DevfileUpdater {
	return func(d *parser.DevfileObj) error {
		containers, err := d.Data.GetComponents(parsercommon.DevfileOptions{
			FilterByName: containerName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(containers)).To(Equal(1))
		containers[0].Attributes = attributes.Attributes{
			"pod-overrides": apiext.JSON{
				Raw: []byte(fmt.Sprintf(`{"spec": {"securityContext": {"fsGroup": %d}}}`, fsGroup)),
			},
		}
		err = d.Data.UpdateComponent(containers[0])
		Expect(err).NotTo(HaveOccurred())
		return nil
	}
}

// ReadRawDevfile parses and validates the Devfile specified and returns its raw content.
func ReadRawDevfile(devfilePath string) parser.DevfileObj {
	d, err := devfile.ParseAndValidateFromFile(devfilePath, "", false)
	Expect(err).ToNot(HaveOccurred())
	return d
}
