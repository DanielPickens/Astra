package docautomation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github\.com/danielpickens/astra/tests/helper"
)

var _ = Describe("doc command reference astra set namespace", func() {
	var commonVar helper.CommonVar
	var commonPath = filepath.Join("command-reference", "docs-mdx", "set-namespace")
	var outputStringFormat = "```console\n$ astra %s\n%s```\n"

	BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		helper.Chdir(commonVar.Context)
		Expect(helper.VerifyFileExists(".astra/env/env.yaml")).To(BeFalse())
	})

	AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	Context("To set an active namespace resource", func() {

		It("Sets a namespace resource to be current active on a kubernetes cluster", func() {
			args := []string{"set", "namespace", "astra-dev"}
			out := helper.Cmd("astra", args...).ShouldPass().Out()
			got := fmt.Sprintf(outputStringFormat, strings.Join(args, " "), helper.StripSpinner(out))
			file := "set_namespace.mdx"
			want := helper.GetMDXContent(filepath.Join(commonPath, file))
			diff := cmp.Diff(want, got)
			Expect(diff).To(BeEmpty(), file)
		})

		It("Sets a project resource to be current active on a openshift cluster", func() {
			args := []string{"set", "project", "astra-dev"}
			out := helper.Cmd("astra", args...).ShouldPass().Out()
			got := fmt.Sprintf(outputStringFormat, strings.Join(args, " "), helper.StripSpinner(out))
			file := "set_project.mdx"
			want := helper.GetMDXContent(filepath.Join(commonPath, file))
			diff := cmp.Diff(want, got)
			Expect(diff).To(BeEmpty(), file)
		})
	})

})
