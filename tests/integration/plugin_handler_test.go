package integration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github\.com/danielpickens/astra/pkg/astra/cli/plugins"
	"github\.com/danielpickens/astra/tests/helper"
)

var sampleScript = []byte(`
#!/bin/sh
echo 'hello'
`)

var _ = Describe("astra plugin functionality", Label(helper.LabelSkipOnOpenShift), func() {
	var tempDir string
	var origPath = os.Getenv("PATH")
	var handler plugins.PluginHandler
	var _ = BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp(os.TempDir(), "astra")
		Expect(err).NotTo(HaveOccurred())
		os.Setenv("PATH", fmt.Sprintf("%s:%s", origPath, tempDir))
		var baseScriptName = "tst-script"
		scriptName := path.Join(tempDir, baseScriptName)
		err = os.WriteFile(scriptName, sampleScript, 0755)
		Expect(err).NotTo(HaveOccurred())
		handler = plugins.NewExecHandler("tst")
	})

	var _ = AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
		os.Setenv("PATH", origPath)
	})

	Context("when an executable with the correct prefix exists on the path", func() {
		It("finds the plugin", func() {
			if runtime.GOOS == "windows" {
				Skip("doesn't find scripts on Windows platform")
			}
			found := handler.Lookup("script")
			Expect(found).Should(Equal(filepath.Join(tempDir, "tst-script")))
		})
	})

	Context("when no executable with the correct prefix exists on the path", func() {
		It("does not find the plugin", func() {
			if runtime.GOOS == "windows" {
				Skip("doesn't find scripts on Windows platform")
			}
			found := handler.Lookup("unknown")
			Expect(found).Should(Equal(""))
		})
	})
})
