package alizer

import (
	"context"

	"github.com/devfile/alizer/pkg/apis/model"

	"github\.com/danielpickens/astra/pkg/api"
)

type DetectedFramework struct {
	Type           model.DevfileType
	DefaultVersion string
	Registry       api.Registry
	Architectures  []string
}

type Client interface {
	DetectFramework(ctx context.Context, path string) (DetectedFramework, error)
	DetectName(path string) (string, error)
	DetectPorts(path string) ([]int, error)
}
