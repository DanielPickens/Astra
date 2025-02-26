package validate

import (
	"fmt"

	devfilev1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// NoComponentsError returns an error if no component is found
type NoComponentsError struct {
}

func (e *NoComponentsError) Error() string {
	return "no components present"
}

// NoContainerComponentError returns an error if no container component is found
type NoContainerComponentError struct {
}

func (e *NoContainerComponentError) Error() string {
	return fmt.Sprintf("astra requires atleast one component of type '%s' in devfile", devfilev1.ContainerComponentType)
}

// UnsupportedastraCommandError returns an error if the command is neither exec nor composite
type UnsupportedastraCommandError struct {
	commandId string
}

func (e *UnsupportedastraCommandError) Error() string {
	return fmt.Sprintf("command %q must be of type \"exec\" or \"composite\"", e.commandId)
}
