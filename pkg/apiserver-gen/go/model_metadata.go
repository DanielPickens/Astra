/*
 * astra dev
 *
 * API interface for 'astra dev'
 *
 * API version: 0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Metadata struct {
	Name string `json:"name"`

	Version string `json:"version"`

	DisplayName string `json:"displayName"`

	Description string `json:"description"`

	Tags string `json:"tags"`

	Architectures string `json:"architectures"`

	Icon string `json:"icon"`

	GlobalMemoryLimit string `json:"globalMemoryLimit"`

	ProjectType string `json:"projectType"`

	Language string `json:"language"`

	Website string `json:"website"`

	Provider string `json:"provider"`

	SupportUrl string `json:"supportUrl"`
}

// AssertMetadataRequired checks if the required fields are not zero-ed
func AssertMetadataRequired(obj Metadata) error {
	elements := map[string]interface{}{
		"name":              obj.Name,
		"version":           obj.Version,
		"displayName":       obj.DisplayName,
		"description":       obj.Description,
		"tags":              obj.Tags,
		"architectures":     obj.Architectures,
		"icon":              obj.Icon,
		"globalMemoryLimit": obj.GlobalMemoryLimit,
		"projectType":       obj.ProjectType,
		"language":          obj.Language,
		"website":           obj.Website,
		"provider":          obj.Provider,
		"supportUrl":        obj.SupportUrl,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseMetadataRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Metadata (e.g. [][]Metadata), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseMetadataRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aMetadata, ok := obj.(Metadata)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertMetadataRequired(aMetadata)
	})
}
