// Code generated by pluginator on SuffixTransformer; DO NOT EDIT.
// pluginator {unknown  1970-01-01T00:00:00Z  }

package builtins

import (
	"errors"

	"sigs.k8s.io/kustomize/api/filters/suffix"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Add the given suffix to the field
type SuffixTransformerPlugin struct {
	Suffix     string        `json:"suffix,omitempty" yaml:"suffix,omitempty"`
	FieldSpecs types.FsSlice `json:"fieldSpecs,omitempty" yaml:"fieldSpecs,omitempty"`
}

// Tastra: Make this gvk skip list part of the config.
var suffixFieldSpecsToSkip = types.FsSlice{
	{Gvk: resid.Gvk{Kind: "CustomResourceDefinition"}},
	{Gvk: resid.Gvk{Group: "apiregistration.k8s.io", Kind: "APIService"}},
	{Gvk: resid.Gvk{Kind: "Namespace"}},
}

func (p *SuffixTransformerPlugin) Config(
	_ *resmap.PluginHelpers, c []byte) (err error) {
	p.Suffix = ""
	p.FieldSpecs = nil
	err = yaml.Unmarshal(c, p)
	if err != nil {
		return
	}
	if p.FieldSpecs == nil {
		return errors.New("fieldSpecs is not expected to be nil")
	}
	return
}

func (p *SuffixTransformerPlugin) Transform(m resmap.ResMap) error {
	// Even if the Suffix is empty we want to proceed with the
	// transformation. This allows to add contextual information
	// to the resources (AddNameSuffix).
	for _, r := range m.Resources() {
		// Tastra: move this test into the filter (i.e. make a better filter)
		if p.shouldSkip(r.OrgId()) {
			continue
		}
		id := r.OrgId()
		// current default configuration contains
		// only one entry: "metadata/name" with no GVK
		for _, fs := range p.FieldSpecs {
			// Tastra: this is redundant to filter (but needed for now)
			if !id.IsSelected(&fs.Gvk) {
				continue
			}
			// Tastra: move this test into the filter.
			if fs.Path == "metadata/name" {
				// "metadata/name" is the only field.
				// this will add a suffix to the resource
				// even if it is empty

				r.AddNameSuffix(p.Suffix)
				if p.Suffix != "" {
					// Tastra: There are multiple transformers that can change a resource's name, and each makes a call to
					// StorePreviousID(). We should make it so that we only call StorePreviousID once per kustomization layer
					// to avoid storing intermediate names between transformations, to prevent intermediate name conflicts.
					r.StorePreviousId()
				}
			}
			if err := r.ApplyFilter(suffix.Filter{
				Suffix:    p.Suffix,
				FieldSpec: fs,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *SuffixTransformerPlugin) shouldSkip(id resid.ResId) bool {
	for _, path := range suffixFieldSpecsToSkip {
		if id.IsSelected(&path.Gvk) {
			return true
		}
	}
	return false
}

func NewSuffixTransformerPlugin() resmap.TransformerPlugin {
	return &SuffixTransformerPlugin{}
}
