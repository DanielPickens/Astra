package helper

import (
	"os"

	_ "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github\.com/danielpickens/astra/pkg/astra/cli/feature"
)

// EnableExperimentalMode enables the experimental mode, so that experimental features of astra can be used.
func EnableExperimentalMode() {
	err := os.Setenv(feature.astraExperimentalModeEnvVar, "true")
	Expect(err).ShouldNot(HaveOccurred())
}

// ResetExperimentalMode disables the experimental mode.
//
// Note that calling any experimental feature of astra right is expected to error out if experimental mode is not enabled.
func ResetExperimentalMode() {
	if _, ok := os.LookupEnv(feature.astraExperimentalModeEnvVar); ok {
		err := os.Unsetenv(feature.astraExperimentalModeEnvVar)
		Expect(err).ShouldNot(HaveOccurred())
	}
}
