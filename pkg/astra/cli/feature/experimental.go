package feature

import (
	"context"

	envcontext "github\.com/danielpickens/astra/pkg/config/context"
)

const astraExperimentalModeEnvVar = "astra_EXPERIMENTAL_MODE"

// IsExperimentalModeEnabled returns whether the experimental mode is enabled or not,
// which means by checking the value of the "astra_EXPERIMENTAL_MODE" environment variable.
func IsExperimentalModeEnabled(ctx context.Context) bool {
	return envcontext.GetEnvConfig(ctx).astraExperimentalMode
}
