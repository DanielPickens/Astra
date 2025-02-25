package integration

import (
	"path"
	"path/filepath"
	"regexp"

	devfilepkg "github.com/devfile/api/v2/pkg/devfile"
	"github.com/tidwall/gjson"

	"github\.com/danielpickens/astra/tests/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("astra list with devfile", Label(helper.LabelSkipOnOpenShift), func() {
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		helper.Chdir(commonVar.Context)
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})
	Context("listing non-astra managed components", func() {
		When("a non-astra managed component is deployed", func() {
			const (
				// hard coded names from the deployment-app-label.yaml
				deploymentName = "example-deployment"
				managedBy      = "some-tool"
			)
			BeforeEach(func() {
				commonVar.CliRunner.Run("create", "-f", helper.GetExamplePath("manifests", "deployment-app-label.yaml"))
			})
			AfterEach(func() {
				commonVar.CliRunner.Run("delete", "-f", helper.GetExamplePath("manifests", "deployment-app-label.yaml"))
			})
			It("should list the component with astra list", func() {
				output := helper.Cmd("astra", "list", "component").ShouldPass().Out()
				helper.MatchAllInOutput(output, []string{deploymentName, "Unknown", "None", managedBy})
			})
			It("should list the component in JSON", func() {
				output := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass().Out()
				helper.JsonPathContentIs(output, "components.#", "1")
				helper.JsonPathContentIs(output, "components.0.name", deploymentName)
				Expect(gjson.Get(output, "components.0.runningIn").String()).To(BeEmpty())
				helper.JsonPathContentIs(output, "components.0.projectType", "Unknown")
				helper.JsonPathContentIs(output, "components.0.managedBy", managedBy)
			})
		})
		When("a non-astra managed component without the managed-by label is deployed", func() {
			const (
				// hard coded names from the deployment-without-managed-by-label.yaml
				deploymentName = "java-springboot-basic"
			)
			BeforeEach(func() {
				commonVar.CliRunner.Run("create", "-f", helper.GetExamplePath("manifests", "deployment-without-managed-by-label.yaml"))
			})
			AfterEach(func() {
				commonVar.CliRunner.Run("delete", "-f", helper.GetExamplePath("manifests", "deployment-without-managed-by-label.yaml"))
			})
			It("should list the component with astra list", func() {
				output := helper.Cmd("astra", "list", "component").ShouldPass().Out()
				helper.MatchAllInOutput(output, []string{deploymentName, "Unknown", "None", "Unknown"})
			})
			It("should list the component in JSON", func() {
				output := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass().Out()
				helper.JsonPathContentIs(output, "components.#", "1")
				helper.JsonPathContentContain(output, "components.0.name", deploymentName)
				Expect(gjson.Get(output, "components.0.runningIn").String()).To(BeEmpty())
				helper.JsonPathContentContain(output, "components.0.projectType", "Unknown")
				helper.JsonPathContentContain(output, "components.0.managedBy", "")
			})
		})
		When("an operator managed deployment(without instance and managed-by label) is deployed", func() {
			deploymentName := "nginx"
			BeforeEach(func() {
				commonVar.CliRunner.Run("create", "deployment", deploymentName, "--image=nginx")
			})
			AfterEach(func() {
				commonVar.CliRunner.Run("delete", "deployment", deploymentName)
			})
			It("should not be listed in the astra list output", func() {
				output := helper.Cmd("astra", "list", "component").ShouldRun().Out()
				Expect(output).ToNot(ContainSubstring(deploymentName))

			})
		})
	})

	When("a component created in 'app' application", func() {

		var devSession helper.DevSession
		var componentName string

		BeforeEach(func() {
			componentName = helper.RandString(6)
			helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
			helper.CopyExampleDevFile(
				filepath.Join("source", "devfiles", "nodejs", "devfile-deploy.yaml"),
				path.Join(commonVar.Context, "devfile.yaml"),
				componentName)
			helper.Chdir(commonVar.Context)
		})

		for _, label := range []string{
			helper.LabelNoCluster, helper.LabelUnauth,
		} {
			label := label
			It("should list the local component when no authenticated", Label(label), func() {
				By("checking the normal output", func() {
					stdOut := helper.Cmd("astra", "list", "component").ShouldPass().Out()
					Expect(stdOut).To(ContainSubstring(componentName))
				})

				By("checking the JSON output", func() {
					res := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					Expect(stderr).To(BeEmpty())
					helper.JsonPathContentIs(stdout, "componentInDevfile", componentName)
					helper.JsonPathContentIs(stdout, "components.0.name", componentName)
				})
			})
		}

		When("dev is running on cluster", func() {
			BeforeEach(func() {
				var err error
				devSession, err = helper.StartDevMode(helper.DevSessionOpts{})
				Expect(err).ToNot(HaveOccurred())
			})
			AfterEach(func() {
				devSession.Stop()
				devSession.WaitEnd()
			})

			var checkList = func(componentType string) {
				By("checking the normal output", func() {
					stdOut := helper.Cmd("astra", "list", "component").ShouldPass().Out()
					Expect(stdOut).To(ContainSubstring(componentType))
				})
			}

			It("should display platform", func() {
				for _, cmd := range [][]string{
					{"list", "component"},
					{"list"},
				} {
					cmd := cmd
					By("returning platform with json output", func() {
						args := append(cmd, "-o", "json")
						res := helper.Cmd("astra", args...).ShouldPass()
						stdout, stderr := res.Out(), res.Err()
						Expect(stderr).To(BeEmpty())
						Expect(helper.IsJSON(stdout)).To(BeTrue(), "output should be in JSON format")
						helper.JsonPathContentIs(stdout, "components.#", "1")
						helper.JsonPathContentIs(stdout, "components.0.runningOn", "cluster") // Deprecated
						helper.JsonPathContentIs(stdout, "components.0.platform", "cluster")
					})
					By("displaying platform", func() {
						stdout := helper.Cmd("astra", cmd...).ShouldPass().Out()
						Expect(stdout).To(ContainSubstring("PLATFORM"))
					})
				}
			})

			Context("verifying the managedBy Version in the astra list output", func() {
				var version string
				BeforeEach(func() {
					versionOut := helper.Cmd("astra", "version").ShouldPass().Out()
					reastraVersion := regexp.MustCompile(`v[0-9]+.[0-9]+.[0-9]+(?:-\w+)?`)
					version = reastraVersion.FindString(versionOut)

				})
				It("should show managedBy Version", func() {
					By("checking the normal output", func() {
						stdout := helper.Cmd("astra", "list", "component").ShouldPass().Out()
						Expect(stdout).To(ContainSubstring(version))
					})
					By("checking the JSON output", func() {
						stdout := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass().Out()
						helper.JsonPathContentContain(stdout, "components.0.managedByVersion", version)
					})
				})
			})

			It("show an astra deploy or dev in the list", func() {
				By("should display the component as 'Dev' in astra list", func() {
					checkList("Dev")
				})

				By("should display the component as 'Dev' in astra list -o json", func() {
					res := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(stderr).To(BeEmpty())
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					helper.JsonPathContentIs(stdout, "components.#", "1")
					helper.JsonPathContentContain(stdout, "components.0.runningIn.dev", "true")
					helper.JsonPathContentContain(stdout, "components.0.runningIn.deploy", "")
				})

				// Fake the astra deploy image build / push passing in "echo" to PODMAN
				stdout := helper.Cmd("astra", "deploy").AddEnv("PODMAN_CMD=echo").ShouldPass().Out()
				By("building and pushing image to registry", func() {
					Expect(stdout).To(ContainSubstring("build -t quay.io/unknown-account/myimage"))
					Expect(stdout).To(ContainSubstring("push quay.io/unknown-account/myimage"))
				})

				By("should display the component as 'Deploy' in astra list", func() {
					checkList("Dev, Deploy")
				})

				By("should display the component as 'Dev, Deploy' in astra list -o json", func() {
					res := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(stderr).To(BeEmpty())
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					helper.JsonPathContentIs(stdout, "components.#", "1")
					helper.JsonPathContentContain(stdout, "components.0.runningIn.dev", "true")
					helper.JsonPathContentContain(stdout, "components.0.runningIn.deploy", "true")
				})
			})
		})

		When("dev is running on podman", Serial, Label(helper.LabelPodman), func() {
			BeforeEach(func() {
				var err error
				devSession, err = helper.StartDevMode(helper.DevSessionOpts{
					RunOnPodman: true,
				})
				Expect(err).ToNot(HaveOccurred())
			})
			AfterEach(func() {
				devSession.Stop()
				devSession.WaitEnd()
			})

			It("should display component depending on platform flag", func() {
				for _, cmd := range [][]string{
					{"list", "component"},
					{"list"},
				} {
					cmd := cmd
					By("returning component in dev mode with json output", func() {
						args := append(cmd, "-o", "json")
						stdout := helper.Cmd("astra", args...).ShouldPass().Out()
						Expect(helper.IsJSON(stdout)).To(BeTrue(), "output should be in JSON format")
						helper.JsonPathContentIs(stdout, "components.#", "1")
						helper.JsonPathContentIs(stdout, "components.0.name", componentName)
						helper.JsonPathContentIs(stdout, "components.0.runningIn.dev", "true")
						helper.JsonPathContentIs(stdout, "components.0.runningOn", "podman") // Deprecated
						helper.JsonPathContentIs(stdout, "components.0.platform", "podman")
					})
					By("returning component not in dev mode with json output and platform is cluster", func() {
						args := append(cmd, "-o", "json", "--platform", "cluster")
						stdout := helper.Cmd("astra", args...).ShouldPass().Out()
						Expect(helper.IsJSON(stdout)).To(BeTrue(), "output should be in JSON format")
						helper.JsonPathContentIs(stdout, "components.#", "1")
						helper.JsonPathContentIs(stdout, "components.0.name", componentName)
						helper.JsonPathContentIs(stdout, "components.0.runningIn.dev", "false")
						helper.JsonPathDoesNotExist(stdout, "components.0.runningOn") // Deprecated
						helper.JsonPathDoesNotExist(stdout, "components.0.platform")
					})
					By("displaying component in dev mode", func() {
						stdout := helper.Cmd("astra", cmd...).ShouldPass().Out()
						Expect(stdout).To(ContainSubstring(componentName))
						Expect(stdout).To(ContainSubstring("PLATFORM"))
						Expect(stdout).To(ContainSubstring("podman"))
						Expect(stdout).To(ContainSubstring("Dev"))
					})
				}
			})
		})
	})

	Context("devfile has missing metadata", func() {
		// Note: We will be using SpringBoot example here because it helps to distinguish between language and projectType.
		// In terms of SpringBoot, spring is the projectType and java is the language; see https://github\.com/danielpickens/astra/issues/4815
		BeforeEach(func() {
			helper.CopyExample(filepath.Join("source", "devfiles", "springboot", "project"), commonVar.Context)
		})
		var metadata devfilepkg.DevfileMetadata

		// checkList checks the list output (both normal and json) to see if it contains the expected componentType
		var checkList = func(componentType string) {
			By("checking the normal output", func() {
				stdOut := helper.Cmd("astra", "list", "component").ShouldPass().Out()
				Expect(stdOut).To(ContainSubstring(componentType))
			})

			By("checking the JSON output", func() {
				res := helper.Cmd("astra", "list", "component", "-o", "json").ShouldPass()
				stdout, stderr := res.Out(), res.Err()
				Expect(stderr).To(BeEmpty())
				Expect(helper.IsJSON(stdout)).To(BeTrue())
				helper.JsonPathContentIs(stdout, "components.#", "1")
				helper.JsonPathContentContain(stdout, "components.0.projectType", componentType)
			})
		}

		When("projectType is missing", func() {
			BeforeEach(func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", helper.GetExamplePath("source", "devfiles", "springboot", "devfile-with-missing-projectType-metadata.yaml")).ShouldPass()
				helper.CreateLocalEnv(commonVar.Context, "aname", commonVar.Project)
				metadata = helper.GetMetadataFromDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
			})

			It("should show the language for 'Type' in astra list", Label(helper.LabelNoCluster), func() {
				checkList(metadata.Language)
			})

			It("should show the language for 'Type' in astra list", Label(helper.LabelUnauth), func() {
				checkList(metadata.Language)
			})

			When("the component is pushed in dev mode", func() {
				var devSession helper.DevSession
				BeforeEach(func() {
					var err error
					devSession, err = helper.StartDevMode(helper.DevSessionOpts{})
					Expect(err).ToNot(HaveOccurred())
				})
				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})

				It("should show the language for 'Type' in astra list", func() {
					checkList(metadata.Language)
				})
			})
		})

		When("projectType and language is missing", func() {
			BeforeEach(func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", helper.GetExamplePath("source", "devfiles", "springboot", "devfile-with-missing-projectType-and-language-metadata.yaml")).ShouldPass()
				helper.CreateLocalEnv(commonVar.Context, "aname", commonVar.Project)
				metadata = helper.GetMetadataFromDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
			})
			It("should show 'Not available' for 'Type' in astra list", Label(helper.LabelNoCluster), func() {
				checkList("Not available")
			})
			It("should show 'Not available' for 'Type' in astra list", Label(helper.LabelUnauth), func() {
				checkList("Not available")
			})
			When("the component is pushed", func() {
				var devSession helper.DevSession
				BeforeEach(func() {
					var err error
					devSession, err = helper.StartDevMode(helper.DevSessionOpts{})
					Expect(err).ToNot(HaveOccurred())
				})
				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})
				It("should show 'nodejs' for 'Type' in astra list", func() {
					checkList("Not available")
				})
			})
		})
	})
})
