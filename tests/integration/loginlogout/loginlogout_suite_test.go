package integration

import (
	"testing"

	"github\.com/danielpickens/astra/tests/helper"
)

func TestLoginlogout(t *testing.T) {
	helper.RunTestSpecs(t, "Login Logout Suite")
}
