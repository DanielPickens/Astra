package e2escenarios

import (
	"testing"

	"github\.com/danielpickens/astra/tests/helper"
)

func TestE2eScenarios(t *testing.T) {
	helper.RunTestSpecs(t, "astra e2e scenarios")
}
