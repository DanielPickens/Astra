/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Endpoint struct {
	Name string `json:"name"`

	Exposure string `json:"exposure,omitempty"`

	Path string `json:"path,omitempty"`

	Protocol string `json:"protocol,omitempty"`

	Secure bool `json:"secure,omitempty"`

	TargetPort int32 `json:"targetPort"`
}

// AssertEndpointRequired checks if the required fields are not zero-ed
func AssertEndpointRequired(obj Endpoint) error {
	elements := map[string]interface{}{
		"name":       obj.Name,
		"targetPort": obj.TargetPort,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseEndpointRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Endpoint (e.g. [][]Endpoint), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseEndpointRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aEndpoint, ok := obj.(Endpoint)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertEndpointRequired(aEndpoint)
	})
}
