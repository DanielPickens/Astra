/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ApplyCommand struct {
	Component string `json:"component"`
}

// AssertApplyCommandRequired checks if the required fields are not zero-ed
func AssertApplyCommandRequired(obj ApplyCommand) error {
	elements := map[string]interface{}{
		"component": obj.Component,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseApplyCommandRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ApplyCommand (e.g. [][]ApplyCommand), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseApplyCommandRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aApplyCommand, ok := obj.(ApplyCommand)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertApplyCommandRequired(aApplyCommand)
	})
}
