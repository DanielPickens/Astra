/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DevstateCompositeCommandCommandNamePatchRequest struct {
	Parallel bool `json:"parallel,omitempty"`

	Commands []string `json:"commands,omitempty"`
}

// AssertDevstateCompositeCommandCommandNamePatchRequestRequired checks if the required fields are not zero-ed
func AssertDevstateCompositeCommandCommandNamePatchRequestRequired(obj DevstateCompositeCommandCommandNamePatchRequest) error {
	return nil
}

// AssertRecurseDevstateCompositeCommandCommandNamePatchRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DevstateCompositeCommandCommandNamePatchRequest (e.g. [][]DevstateCompositeCommandCommandNamePatchRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDevstateCompositeCommandCommandNamePatchRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDevstateCompositeCommandCommandNamePatchRequest, ok := obj.(DevstateCompositeCommandCommandNamePatchRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDevstateCompositeCommandCommandNamePatchRequestRequired(aDevstateCompositeCommandCommandNamePatchRequest)
	})
}
