//go:build !race
// +build !race

package integration

import (
	"testing"

	"github\.com/danielpickens/astra/tests/helper"
)

func TestIntegration(t *testing.T) {
	helper.RunTestSpecs(t, "Integration Suite")
}
