/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ExecCommand struct {
	Component string `json:"component"`

	CommandLine string `json:"commandLine"`

	WorkingDir string `json:"workingDir"`

	HotReloadCapable bool `json:"hotReloadCapable"`
}

// AssertExecCommandRequired checks if the required fields are not zero-ed
func AssertExecCommandRequired(obj ExecCommand) error {
	elements := map[string]interface{}{
		"component":        obj.Component,
		"commandLine":      obj.CommandLine,
		"workingDir":       obj.WorkingDir,
		"hotReloadCapable": obj.HotReloadCapable,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseExecCommandRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ExecCommand (e.g. [][]ExecCommand), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseExecCommandRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aExecCommand, ok := obj.(ExecCommand)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertExecCommandRequired(aExecCommand)
	})
}
