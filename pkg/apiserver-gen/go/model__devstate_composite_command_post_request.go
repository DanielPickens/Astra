/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DevstateCompositeCommandPostRequest struct {

	// Name of the command
	Name string `json:"name,omitempty"`

	Parallel bool `json:"parallel,omitempty"`

	Commands []string `json:"commands,omitempty"`
}

// AssertDevstateCompositeCommandPostRequestRequired checks if the required fields are not zero-ed
func AssertDevstateCompositeCommandPostRequestRequired(obj DevstateCompositeCommandPostRequest) error {
	return nil
}

// AssertRecurseDevstateCompositeCommandPostRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DevstateCompositeCommandPostRequest (e.g. [][]DevstateCompositeCommandPostRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDevstateCompositeCommandPostRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDevstateCompositeCommandPostRequest, ok := obj.(DevstateCompositeCommandPostRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDevstateCompositeCommandPostRequestRequired(aDevstateCompositeCommandPostRequest)
	})
}
