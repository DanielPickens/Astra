/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ImageCommand struct {
	Component string `json:"component"`
}

// AssertImageCommandRequired checks if the required fields are not zero-ed
func AssertImageCommandRequired(obj ImageCommand) error {
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

// AssertRecurseImageCommandRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ImageCommand (e.g. [][]ImageCommand), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseImageCommandRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aImageCommand, ok := obj.(ImageCommand)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertImageCommandRequired(aImageCommand)
	})
}
