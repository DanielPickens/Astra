//go:build !race
// +build !race

package docautomation

import (
	"testing"

	"github\.com/danielpickens/astra/tests/helper"
)

func TestUserGuidesDocAutomation(t *testing.T) {
	helper.RunTestSpecs(t, "User Guides Doc Automation Suite")
}
