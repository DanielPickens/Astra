/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DevfileContent struct {
	Content string `json:"content"`

	Version string `json:"version"`

	Commands []Command `json:"commands"`

	Containers []Container `json:"containers"`

	Images []Image `json:"images"`

	Resources []Resource `json:"resources"`

	Volumes []Volume `json:"volumes"`

	Events Events `json:"events"`

	Metadata Metadata `json:"metadata"`
}

// AssertDevfileContentRequired checks if the required fields are not zero-ed
func AssertDevfileContentRequired(obj DevfileContent) error {
	elements := map[string]interface{}{
		"content":    obj.Content,
		"version":    obj.Version,
		"commands":   obj.Commands,
		"containers": obj.Containers,
		"images":     obj.Images,
		"resources":  obj.Resources,
		"volumes":    obj.Volumes,
		"events":     obj.Events,
		"metadata":   obj.Metadata,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Commands {
		if err := AssertCommandRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Containers {
		if err := AssertContainerRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Images {
		if err := AssertImageRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Resources {
		if err := AssertResourceRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Volumes {
		if err := AssertVolumeRequired(el); err != nil {
			return err
		}
	}
	if err := AssertEventsRequired(obj.Events); err != nil {
		return err
	}
	if err := AssertMetadataRequired(obj.Metadata); err != nil {
		return err
	}
	return nil
}

// AssertRecurseDevfileContentRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DevfileContent (e.g. [][]DevfileContent), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDevfileContentRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDevfileContent, ok := obj.(DevfileContent)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDevfileContentRequired(aDevfileContent)
	})
}
