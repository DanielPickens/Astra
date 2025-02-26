/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Resource struct {
	Name string `json:"name"`

	Inlined string `json:"inlined,omitempty"`

	Uri string `json:"uri,omitempty"`

	DeployByDefault string `json:"deployByDefault"`

	// true if the resource is not referenced in any command
	Orphan bool `json:"orphan"`
}

// AssertResourceRequired checks if the required fields are not zero-ed
func AssertResourceRequired(obj Resource) error {
	elements := map[string]interface{}{
		"name":            obj.Name,
		"deployByDefault": obj.DeployByDefault,
		"orphan":          obj.Orphan,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseResourceRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Resource (e.g. [][]Resource), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseResourceRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aResource, ok := obj.(Resource)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertResourceRequired(aResource)
	})
}
