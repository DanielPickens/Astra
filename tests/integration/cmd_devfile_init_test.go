package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"

	"github\.com/danielpickens/astra/pkg/config"
	envcontext "github\.com/danielpickens/astra/pkg/config/context"
	"github\.com/danielpickens/astra/pkg/astra/cli/messages"
	"github\.com/danielpickens/astra/pkg/preference"
	"github\.com/danielpickens/astra/pkg/segment"
	segmentContext "github\.com/danielpickens/astra/pkg/segment/context"
	"github\.com/danielpickens/astra/pkg/util"

	"github\.com/danielpickens/astra/tests/helper"
)

var _ = Describe("astra devfile init command tests", Label(helper.LabelSkipOnOpenShift), func() {
	var commonVar helper.CommonVar

	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		helper.Chdir(commonVar.Context)
		Expect(helper.VerifyFileExists(".astra/env/env.yaml")).To(BeFalse())
	})

	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	for _, label := range []string{
		helper.LabelNoCluster, helper.LabelUnauth,
	} {
		label := label
		var _ = Context("label "+label, Label(label), func() {

			It("astra init should fail", func() {
				By("running astra init with incomplete flags", func() {
					helper.Cmd("astra", "init", "--name", "aname").ShouldFail()
				})

				By("using an invalid component name", func() {
					helper.Cmd("astra", "init", "--devfile", "go", "--name", "123").ShouldFail()
				})

				By("running astra init with json and no other flags", func() {
					res := helper.Cmd("astra", "init", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
					helper.JsonPathContentIs(stderr, "message", "parameters are expected to select a devfile")
				})

				By("running astra init with incomplete flags and JSON output", func() {
					res := helper.Cmd("astra", "init", "--name", "aname", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
					helper.JsonPathContentContain(stderr, "message", "either --devfile or --devfile-path parameter should be specified")
				})

				By("keeping an empty directory when running astra init with wrong starter name", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile", "go", "--starter", "wrongname").ShouldFail()
					files := helper.ListFilesInDir(commonVar.Context)
					Expect(len(files)).To(Equal(0))
				})
				By("using an invalid devfile name", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile", "invalid").ShouldFail()
				})

				for _, devfileName := range []string{"devfile.yaml", ".devfile.yaml", "devfile.yml", ".devfile.yml"} {
					devfileName := devfileName
					By("running astra init in a directory containing a "+devfileName, func() {
						helper.CopyExampleDevFile(
							filepath.Join("source", "devfiles", "nodejs", "devfile-registry.yaml"),
							filepath.Join(commonVar.Context, devfileName),
							"")
						defer os.Remove(filepath.Join(commonVar.Context, devfileName))
						err := helper.Cmd("astra", "init").ShouldFail().Err()
						Expect(err).To(ContainSubstring("a devfile already exists in the current directory"))
					})
				}

				By("running astra init with wrong local file path given to --devfile-path", func() {
					err := helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", "/some/path/devfile.yaml").ShouldFail().Err()
					Expect(err).To(ContainSubstring("unable to download devfile"))
				})
				By("running astra init with wrong URL path given to --devfile-path", func() {
					err := helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", "https://github.com/path/to/devfile.yaml").ShouldFail().Err()
					Expect(err).To(ContainSubstring("unable to download devfile"))
				})
				By("running astra init multiple times", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile", "nodejs").ShouldPass()
					defer helper.DeleteFile(filepath.Join(commonVar.Context, "devfile.yaml"))
					output := helper.Cmd("astra", "init", "--name", "aname", "--devfile", "nodejs").ShouldFail().Err()
					Expect(output).To(ContainSubstring("a devfile already exists in the current directory"))
				})

				By("running astra init with --devfile-path and --devfile-registry", func() {
					errOut := helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", "https://github.com/path/to/devfile.yaml", "--devfile-registry", "DefaultDevfileRegistry").ShouldFail().Err()
					Expect(errOut).To(ContainSubstring("--devfile-registry parameter cannot be used with --devfile-path"))
				})
				By("running astra init with invalid --devfile-registry value", func() {
					fakeRegistry := "fake"
					errOut := helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", "https://github.com/path/to/devfile.yaml", "--devfile-registry", fakeRegistry).ShouldFail().Err()
					Expect(errOut).To(ContainSubstring(fmt.Sprintf("%q not found", fakeRegistry)))
				})
			})

			Context("running astra init with valid flags", func() {
				When("using --devfile flag", func() {
					compName := "aname"
					var output string
					BeforeEach(func() {
						output = helper.Cmd("astra", "init", "--name", compName, "--devfile", "go").ShouldPass().Out()
					})

					It("should download a devfile.yaml file and correctly set the component name in it", func() {
						By("not showing the interactive mode notice message", func() {
							Expect(output).ShouldNot(ContainSubstring(messages.InteractiveModeEnabled))
						})
						files := helper.ListFilesInDir(commonVar.Context)
						Expect(files).To(SatisfyAll(
							HaveLen(2),
							ContainElements("devfile.yaml", util.DotastraDirectory)))
						metadata := helper.GetMetadataFromDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
						Expect(metadata.Name).To(BeEquivalentTo(compName))
					})
				})
				When("using --devfile flag and JSON output", func() {
					compName := "aname"
					var res *helper.CmdWrapper
					BeforeEach(func() {
						res = helper.Cmd("astra", "init", "--name", compName, "--devfile", "go", "-o", "json").ShouldPass()
					})

					It("should return correct values in output", func() {
						stdout, stderr := res.Out(), res.Err()
						Expect(stderr).To(BeEmpty())
						Expect(helper.IsJSON(stdout)).To(BeTrue())
						helper.JsonPathContentIs(stdout, "devfilePath", filepath.Join(commonVar.Context, "devfile.yaml"))
						helper.JsonPathContentIs(stdout, "devfileData.devfile.metadata.name", compName)
						helper.JsonPathContentIs(stdout, "devfileData.supportedastraFeatures.dev", "true")
						helper.JsonPathContentIs(stdout, "devfileData.supportedastraFeatures.debug", "true")
						helper.JsonPathContentIs(stdout, "devfileData.supportedastraFeatures.deploy", "false")
						helper.JsonPathContentIs(stdout, "managedBy", "astra")
					})
				})

				const (
					devfileName = "go"
				)
				for _, ctx := range []struct {
					title, devfileVersion, requiredVersion string
					checkVersion                           func(metadataVersion string)
				}{
					{
						title:          "to download the latest version",
						devfileVersion: "latest",
						checkVersion: func(metadataVersion string) {
							reg := helper.NewRegistry(commonVar.GetDevfileRegistryURL())
							stack, err := reg.GetStack(devfileName)
							Expect(err).ToNot(HaveOccurred())
							Expect(len(stack.Versions)).ToNot(BeZero())
							lastVersion := stack.Versions[0]
							Expect(metadataVersion).To(BeEquivalentTo(lastVersion.Version))
						},
					},
					{
						title:          "to download a specific version",
						devfileVersion: "1.0.2",
						checkVersion: func(metadataVersion string) {
							Expect(metadataVersion).To(BeEquivalentTo("1.0.2"))
						},
					},
				} {
					ctx := ctx
					When(fmt.Sprintf("using --devfile-version flag %s", ctx.title), func() {
						BeforeEach(func() {
							helper.Cmd("astra", "init", "--name", "aname", "--devfile", devfileName, "--devfile-version", ctx.devfileVersion).ShouldPass()
						})

						It("should download the devfile with the requested version", func() {
							files := helper.ListFilesInDir(commonVar.Context)
							Expect(files).To(ContainElements("devfile.yaml"))
							metadata := helper.GetMetadataFromDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
							ctx.checkVersion(metadata.Version)
						})
					})

					When(fmt.Sprintf("using --devfile-version flag and JSON output %s", ctx.title), func() {
						var res *helper.CmdWrapper
						BeforeEach(func() {
							res = helper.Cmd("astra", "init", "--name", "aname", "--devfile", devfileName, "--devfile-version", ctx.devfileVersion, "-o", "json").ShouldPass()
						})

						It("should show the requested devfile version", func() {
							stdout := res.Out()
							Expect(helper.IsJSON(stdout)).To(BeTrue())
							ctx.checkVersion(gjson.Get(stdout, "devfileData.devfile.metadata.version").String())
						})
					})
				}

				When("using --devfile-path flag with a local devfile", func() {
					var newContext string
					BeforeEach(func() {
						newContext = helper.CreateNewContext()
						newDevfilePath := filepath.Join(newContext, "devfile.yaml")
						helper.CopyExampleDevFile(
							filepath.Join("source", "devfiles", "nodejs", "devfile-registry.yaml"), newDevfilePath, "")
						helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", newDevfilePath).ShouldPass()
					})
					AfterEach(func() {
						helper.DeleteDir(newContext)
					})
					It("should copy the devfile.yaml file", func() {
						files := helper.ListFilesInDir(commonVar.Context)
						Expect(files).To(SatisfyAll(
							HaveLen(2),
							ContainElements(util.DotastraDirectory, "devfile.yaml")))
					})
				})
				When("using --devfile-path flag with a URL", func() {
					BeforeEach(func() {
						helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", "https://raw.githubusercontent.com/astra-devfiles/registry/master/devfiles/nodejs/devfile.yaml").ShouldPass()
					})
					It("should copy the devfile.yaml file", func() {
						files := helper.ListFilesInDir(commonVar.Context)
						Expect(files).To(SatisfyAll(
							HaveLen(2),
							ContainElements("devfile.yaml", util.DotastraDirectory)))
					})
				})
				When("using --devfile-registry flag", func() {
					It("should successfully run astra init if specified registry is valid", func() {
						helper.Cmd("astra", "init", "--name", "aname", "--devfile", "go", "--devfile-registry", "DefaultDevfileRegistry").ShouldPass()
					})

				})
			})
			When("a dangling env file exists in the working directory", func() {
				BeforeEach(func() {
					helper.CreateLocalEnv(commonVar.Context, "aname", commonVar.Project)
				})
				It("should successfully create a devfile component and remove the dangling env file", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile", "go").ShouldPass()
				})
			})

			When("running astra init with a devfile that has a subDir starter project", func() {
				BeforeEach(func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", helper.GetExamplePath("source", "devfiles", "springboot", "devfile-with-subDir.yaml"), "--starter", "springbootproject").ShouldPass()
				})

				It("should successfully extract the project in the specified subDir path", func() {
					var found, notToBeFound int
					pathsToValidate := map[string]bool{
						filepath.Join(commonVar.Context, "java", "com"):                                            true,
						filepath.Join(commonVar.Context, "java", "com", "example"):                                 true,
						filepath.Join(commonVar.Context, "java", "com", "example", "demo"):                         true,
						filepath.Join(commonVar.Context, "java", "com", "example", "demo", "DemoApplication.java"): true,
						filepath.Join(commonVar.Context, "resources", "application.properties"):                    true,
					}
					pathsNotToBePresent := map[string]bool{
						filepath.Join(commonVar.Context, "src"):  true,
						filepath.Join(commonVar.Context, "main"): true,
					}
					err := filepath.Walk(commonVar.Context, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if ok := pathsToValidate[path]; ok {
							found++
						}
						if ok := pathsNotToBePresent[path]; ok {
							notToBeFound++
						}
						return nil
					})
					Expect(err).To(BeNil())

					Expect(found).To(Equal(len(pathsToValidate)))
					Expect(notToBeFound).To(Equal(0))
				})
			})

			It("should successfully run astra init for devfile with starter project from the specified branch", func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-branch.yaml"), "--starter", "nodejs-starter").ShouldPass()
				expectedFiles := []string{"package.json", "package-lock.json", "README.md", "devfile.yaml", "test"}
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements(expectedFiles))
			})

			It("should successfully run astra init for devfile with starter project from the specified tag", func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-tag.yaml"), "--starter", "nodejs-starter").ShouldPass()
				expectedFiles := []string{"package.json", "package-lock.json", "README.md", "devfile.yaml", "app"}
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements(expectedFiles))
			})

			It("should successfully run astra init for devfile with starter project on git with main default branch", func() {
				helper.Cmd("astra", "init",
					"--name", "vertx",
					"--devfile-path", helper.GetExamplePath("source", "devfiles", "java", "devfile-with-git-main-branch.yaml"),
					"--starter", "vertx-http-example-redhat",
				).ShouldPass()
			})

			When("running astra init from a directory with sources", func() {
				BeforeEach(func() {
					helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
				})
				It("should work without --starter flag", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile", "nodejs").ShouldPass()
				})
				It("should not accept --starter flag", func() {
					err := helper.Cmd("astra", "init", "--name", "aname", "--devfile", "nodejs", "--starter", "nodejs-starter").ShouldFail().Err()
					Expect(err).To(ContainSubstring("--starter parameter cannot be used when the directory is not empty"))
				})
			})
			Context("checking astra init final output message", func() {
				var newContext, devfilePath string

				BeforeEach(func() {
					newContext = helper.CreateNewContext()
					devfilePath = filepath.Join(newContext, "devfile.yaml")
				})

				AfterEach(func() {
					helper.DeleteDir(newContext)
				})

				When("the devfile used by `astra init` does not contain a deploy command", func() {
					var out string

					BeforeEach(func() {
						helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile.yaml"), devfilePath, "")
						out = helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", devfilePath).ShouldPass().Out()
					})

					It("should only show information about `astra dev`, and not `astra deploy`", func() {
						Expect(out).To(ContainSubstring("astra dev"))
						Expect(out).ToNot(ContainSubstring("astra deploy"))
					})

					It("should not show the interactive mode notice message", func() {
						Expect(out).ShouldNot(ContainSubstring(messages.InteractiveModeEnabled))
					})
				})

				When("the devfile used by `astra init` contains a deploy command", func() {
					var out string

					BeforeEach(func() {
						helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-deploy.yaml"), devfilePath, "")
						out = helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", devfilePath).ShouldPass().Out()
					})

					It("should show information about both `astra dev`, and `astra deploy`", func() {
						Expect(out).To(ContainSubstring("astra dev"))
						Expect(out).To(ContainSubstring("astra deploy"))
					})

					It("should not show the interactive mode notice message", func() {
						Expect(out).ShouldNot(ContainSubstring(messages.InteractiveModeEnabled))
					})
				})
			})

			When("devfile contains parent URI", func() {
				var originalKeyList []string
				var srcDevfile string

				BeforeEach(func() {
					var err error
					srcDevfile = helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-parent.yaml")
					originalDevfileContent, err := os.ReadFile(srcDevfile)
					Expect(err).To(BeNil())
					var content map[string]interface{}
					Expect(yaml.Unmarshal(originalDevfileContent, &content)).To(BeNil())
					for k := range content {
						originalKeyList = append(originalKeyList, k)
					}
				})

				It("should not replace the original devfile", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile-path", srcDevfile).ShouldPass()
					devfileContent, err := os.ReadFile(filepath.Join(commonVar.Context, "devfile.yaml"))
					Expect(err).To(BeNil())
					var content map[string]interface{}
					Expect(yaml.Unmarshal(devfileContent, &content)).To(BeNil())
					for k := range content {
						Expect(k).To(BeElementOf(originalKeyList))
					}
				})
			})

			When("source directory is empty", func() {
				BeforeEach(func() {
					Expect(helper.ListFilesInDir(commonVar.Context)).To(HaveLen(0))
				})

				It("name in devfile is personalized in non-interactive mode", func() {
					helper.Cmd("astra", "init", "--name", "aname", "--devfile-path",
						filepath.Join(helper.GetExamplePath(), "source", "devfiles", "nodejs",
							"devfile-with-starter-with-devfile.yaml")).ShouldPass()

					metadata := helper.GetMetadataFromDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
					Expect(metadata.Name).To(BeEquivalentTo("aname"))
					Expect(metadata.Language).To(BeEquivalentTo("nodejs"))
				})
			})

			Describe("telemetry", func() {

				for _, tt := range []struct {
					name string
					env  map[string]string
				}{
					{
						name: "astra_DISABLE_TELEMETRY=true and astra_TRACKING_CONSENT=yes",
						env: map[string]string{
							//lint:ignore SA1019 We deprecated this env var, but until it is removed, we still want to test it
							segment.DisableTelemetryEnv: "true",
							segment.TrackingConsentEnv:  "yes",
						},
					},
					{
						name: "astra_DISABLE_TELEMETRY=false and astra_TRACKING_CONSENT=no",
						env: map[string]string{
							//lint:ignore SA1019 We deprecated this env var, but until it is removed, we still want to test it
							segment.DisableTelemetryEnv: "false",
							segment.TrackingConsentEnv:  "no",
						},
					},
				} {
					tt := tt
					It("should error out if "+tt.name, func() {
						cmd := helper.Cmd("astra", "init", "--name", "aname", "--devfile", "go")
						for k, v := range tt.env {
							cmd = cmd.AddEnv(fmt.Sprintf("%s=%s", k, v))
						}
						stderr := cmd.ShouldFail().Err()

						//lint:ignore SA1019 We deprecated this env var, but until it is removed, we still want to test it
						Expect(stderr).To(ContainSubstring("%s and %s values are in conflict.", segment.DisableTelemetryEnv, segment.TrackingConsentEnv))
					})
				}

				type telemetryTest struct {
					title         string
					env           map[string]string
					setupFunc     func(cfg preference.Client)
					callerChecker func(stdout, stderr string, data segment.TelemetryData)
				}
				allowedTelemetryCallers := []string{segmentContext.VSCode, segmentContext.IntelliJ, segmentContext.JBoss}
				telemetryTests := []telemetryTest{
					{
						title: "no caller env var",
						callerChecker: func(_, _ string, td segment.TelemetryData) {
							cmdProperties := td.Properties.CmdProperties
							Expect(cmdProperties).Should(HaveKey(segmentContext.Caller))
							Expect(cmdProperties[segmentContext.Caller]).To(BeEmpty())
						},
					},
					{
						title: "empty caller env var",
						env: map[string]string{
							helper.TelemetryCaller: "",
						},
						callerChecker: func(_, _ string, td segment.TelemetryData) {
							cmdProperties := td.Properties.CmdProperties
							Expect(cmdProperties).Should(HaveKey(segmentContext.Caller))
							Expect(cmdProperties[segmentContext.Caller]).To(BeEmpty())
						},
					},
					{
						title: "invalid caller env var",
						env: map[string]string{
							helper.TelemetryCaller: "an-invalid-caller",
						},
						callerChecker: func(stdout, stderr string, td segment.TelemetryData) {
							By("not disclosing list of allowed values", func() {
								helper.DontMatchAllInOutput(stdout, allowedTelemetryCallers)
								helper.DontMatchAllInOutput(stderr, allowedTelemetryCallers)
							})

							By("setting the value as caller property in telemetry even if it is invalid", func() {
								Expect(td.Properties.CmdProperties[segmentContext.Caller]).To(Equal("an-invalid-caller"))
							})
						},
					},
					{
						title: "astra_TRACKING_CONSENT=yes env var should take precedence over ConsentTelemetry preference",
						env:   map[string]string{segment.TrackingConsentEnv: "yes"},
						callerChecker: func(_, _ string, td segment.TelemetryData) {
							cmdProperties := td.Properties.CmdProperties
							Expect(cmdProperties).Should(HaveKey(segmentContext.Caller))
							Expect(cmdProperties[segmentContext.Caller]).To(BeEmpty())
						},
						setupFunc: func(cfg preference.Client) {
							err := cfg.SetConfiguration(preference.ConsentTelemetrySetting, "false")
							Expect(err).ShouldNot(HaveOccurred())
						},
					},
				}
				for _, c := range allowedTelemetryCallers {
					c := c
					telemetryTests = append(telemetryTests, telemetryTest{
						title: fmt.Sprintf("valid caller env var: %s", c),
						env: map[string]string{
							helper.TelemetryCaller: c,
						},
						callerChecker: func(_, _ string, td segment.TelemetryData) {
							Expect(td.Properties.CmdProperties[segmentContext.Caller]).To(Equal(c))
						},
					})
				}
				for _, tt := range telemetryTests {
					tt := tt
					When("recording telemetry data with "+tt.title, func() {
						var stdout string
						var stderr string
						BeforeEach(func() {
							helper.EnableTelemetryDebug()

							ctx := context.Background()
							envConfig, err := config.GetConfiguration()
							Expect(err).To(BeNil())
							ctx = envcontext.WithEnvConfig(ctx, *envConfig)

							cfg, err := preference.NewClient(ctx)
							Expect(err).ShouldNot(HaveOccurred())
							if tt.setupFunc != nil {
								tt.setupFunc(cfg)
							}

							cmd := helper.Cmd("astra", "init", "--name", "aname", "--devfile", "go")
							for k, v := range tt.env {
								cmd = cmd.AddEnv(fmt.Sprintf("%s=%s", k, v))
							}
							stdout, stderr = cmd.ShouldPass().OutAndErr()
						})

						AfterEach(func() {
							helper.ResetTelemetry()
						})

						It("should record the telemetry data correctly", func() {
							td := helper.GetTelemetryDebugData()
							Expect(td.Event).To(ContainSubstring("astra init"))
							Expect(td.Properties.Success).To(BeTrue())
							Expect(td.Properties.Error == "").To(BeTrue())
							Expect(td.Properties.ErrorType == "").To(BeTrue())
							Expect(td.Properties.CmdProperties[segmentContext.DevfileName]).To(ContainSubstring("aname"))
							Expect(td.Properties.CmdProperties[segmentContext.ComponentType]).To(ContainSubstring("Go"))
							Expect(td.Properties.CmdProperties[segmentContext.Language]).To(ContainSubstring("Go"))
							Expect(td.Properties.CmdProperties[segmentContext.ProjectType]).To(ContainSubstring("Go"))
							Expect(td.Properties.CmdProperties[segmentContext.Flags]).To(ContainSubstring("devfile name"))
							Expect(td.Properties.CmdProperties[segmentContext.Platform]).To(BeNil())
							Expect(td.Properties.CmdProperties[segmentContext.PlatformVersion]).To(BeNil())
							tt.callerChecker(stdout, stderr, td)
						})

					})
				}
			})
		})
	}

	When("DevfileRegistriesList CRD is installed on cluster", func() {
		BeforeEach(func() {
			if !helper.IsKubernetesCluster() {
				Skip("skipped on non Kubernetes clusters")
			}
			devfileRegistriesLists := commonVar.CliRunner.Run("apply", "-f", helper.GetExamplePath("manifests", "devfileregistrieslists.yaml"))
			Expect(devfileRegistriesLists.ExitCode()).To(BeEquivalentTo(0))
		})

		When("CR for devfileregistrieslists is installed in namespace", func() {
			BeforeEach(func() {
				manifestFilePath := filepath.Join(commonVar.ConfigDir, "devfileRegistriesListCR.yaml")
				// NOTE: Use reachable URLs as we might be on a cluster with the registry operator installed, which would perform validations.
				err := helper.CreateFileWithContent(manifestFilePath, fmt.Sprintf(`
apiVersion: registry.devfile.io/v1alpha1
kind: DevfileRegistriesList
metadata:
  name: namespace-list
spec:
  devfileRegistries:
    - name: ns-devfile-reg
      url: %q
`, commonVar.GetDevfileRegistryURL()))
				Expect(err).ToNot(HaveOccurred())
				command := commonVar.CliRunner.Run("-n", commonVar.Project, "apply", "-f", manifestFilePath)
				Expect(command.ExitCode()).To(BeEquivalentTo(0))
			})

			It("should be able to download devfile from the in-cluster registry", func() {
				out := helper.Cmd("astra", "init", "--devfile-registry", "ns-devfile-reg", "--devfile", "go", "--name", "go-devfile").ShouldPass().Out()
				Expect(out).To(ContainSubstring("Downloading devfile \"go\" from registry \"ns-devfile-reg\""))
				helper.VerifyFileExists(filepath.Join(commonVar.Context, "devfile.yaml"))
			})
		})
	})

	Context("setting application ports", func() {
		When("running astra init --run-port with a Devfile with no commands", func() {
			BeforeEach(func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path",
					filepath.Join(helper.GetExamplePath(), "source", "devfiles", "nodejs", "devfile-without-commands.yaml"),
					"--run-port", "1234", "--run-port", "2345", "--run-port", "3456").ShouldPass().Out()
			})
			It("should ignore the run ports", func() {
				d := helper.ReadRawDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
				components, err := d.Data.GetComponents(common.DevfileOptions{
					ComponentOptions: common.ComponentOptions{
						ComponentType: v1alpha2.ContainerComponentType,
					},
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(components).To(HaveLen(1))
				Expect(components[0].Name).Should(Equal("runtime"))
				Expect(components[0].Container).ShouldNot(BeNil())
				Expect(components[0].Container.Endpoints).Should(HaveLen(2))
				Expect(components[0].Container.Endpoints[0].Name).Should(Equal("http-node"))
				Expect(components[0].Container.Endpoints[0].TargetPort).Should(Equal(3000))
				Expect(components[0].Container.Endpoints[1].Name).Should(Equal("debug"))
				Expect(components[0].Container.Endpoints[1].TargetPort).Should(Equal(5858))
				Expect(components[0].Container.Endpoints[1].Exposure).Should(Equal(v1alpha2.NoneEndpointExposure))
			})
		})

		When("running astra init --run-port with a Devfile with no commands", func() {
			BeforeEach(func() {
				helper.Cmd("astra", "init", "--name", "aname", "--devfile-path",
					filepath.Join(helper.GetExamplePath(), "source", "devfiles", "nodejs", "devfile-with-debugrun.yaml"),
					"--run-port", "1234", "--run-port", "2345", "--run-port", "3456").ShouldPass().Out()
			})
			It("should overwrite the ports into the container component referenced by the default run command", func() {
				d := helper.ReadRawDevfile(filepath.Join(commonVar.Context, "devfile.yaml"))
				components, err := d.Data.GetComponents(common.DevfileOptions{
					ComponentOptions: common.ComponentOptions{
						ComponentType: v1alpha2.ContainerComponentType,
					},
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(components).To(HaveLen(1))
				Expect(components[0].Name).Should(Equal("runtime"))
				Expect(components[0].Container).ShouldNot(BeNil())
				Expect(components[0].Container.Endpoints).Should(HaveLen(3))
				for i, p := range []int{1234, 2345, 3456} {
					Expect(components[0].Container.Endpoints[i].Name).Should(Equal(fmt.Sprintf("port-%d-tcp", p)))
					Expect(components[0].Container.Endpoints[i].TargetPort).Should(Equal(p))
					Expect(components[0].Container.Endpoints[i].Protocol).Should(Equal(v1alpha2.TCPEndpointProtocol))
				}
			})
		})
	})
})
