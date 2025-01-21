package logs

import (
	"fmt"

	astralabels "github\.com/danielpickens/astra/pkg/labels"
)

type InvalidModeError struct {
	mode string
}

func (e InvalidModeError) Error() string {
	return fmt.Sprintf("invalid mode %q; valid modes are %q, %q, and %q", e.mode, astralabels.ComponentDevMode, astralabels.ComponentDeployMode, astralabels.ComponentAnyMode)
}
