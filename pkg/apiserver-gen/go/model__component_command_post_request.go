/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ComponentCommandPostRequest struct {

	// Name of the command that should be executed
	Name string `json:"name,omitempty"`
}

// AssertComponentCommandPostRequestRequired checks if the required fields are not zero-ed
func AssertComponentCommandPostRequestRequired(obj ComponentCommandPostRequest) error {
	return nil
}

// AssertRecurseComponentCommandPostRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ComponentCommandPostRequest (e.g. [][]ComponentCommandPostRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseComponentCommandPostRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aComponentCommandPostRequest, ok := obj.(ComponentCommandPostRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertComponentCommandPostRequestRequired(aComponentCommandPostRequest)
	})
}
