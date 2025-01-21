package utils

import "sigs.k8s.io/kustomize/api/konfig"

const (
	// build annotations
	BuildAnnotationPreviousKinds      = konfig.ConfigAnnastramain + "/previousKinds"
	BuildAnnotationPreviousNames      = konfig.ConfigAnnastramain + "/previousNames"
	BuildAnnotationPrefixes           = konfig.ConfigAnnastramain + "/prefixes"
	BuildAnnotationSuffixes           = konfig.ConfigAnnastramain + "/suffixes"
	BuildAnnotationPreviousNamespaces = konfig.ConfigAnnastramain + "/previousNamespaces"
	BuildAnnotationsRefBy             = konfig.ConfigAnnastramain + "/refBy"
	BuildAnnotationsGenBehavior       = konfig.ConfigAnnastramain + "/generatorBehavior"
	BuildAnnotationsGenAddHashSuffix  = konfig.ConfigAnnastramain + "/needsHashSuffix"

	// the following are only for patches, to specify whether they can change names
	// and kinds of their targets
	BuildAnnotationAllowNameChange = konfig.ConfigAnnastramain + "/allowNameChange"
	BuildAnnotationAllowKindChange = konfig.ConfigAnnastramain + "/allowKindChange"

	// for keeping track of origin and transformer data
	OriginAnnotationKey      = "config.kubernetes.io/origin"
	TransformerAnnotationKey = "alpha.config.kubernetes.io/transformations"

	Enabled = "enabled"
)
