package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github\.com/danielpickens/astra/tests/helper"
)

var _ = Describe("astra generic", Label(helper.LabelSkipOnOpenShift), func() {
	// Tastra: A neater way to provide astra path. Currently we assume \
	// astra and oc in $PATH already
	var oc helper.OcRunner
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		oc = helper.NewOcRunner("oc")
		commonVar = helper.CommonBeforeEach()
	})

	// Clean up after the test
	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	for _, label := range []string{
		helper.LabelNoCluster, helper.LabelUnauth,
	} {
		label := label
		Context("label "+label, Label(label), func() {
			When("running astra --help", func() {
				var output string
				BeforeEach(func() {
					output = helper.Cmd("astra", "--help").ShouldPass().Out()
				})
				It("returns full help contents including usage, examples, commands, utility commands, component shortcuts, and flags sections", func() {
					helper.MatchAllInOutput(output, []string{"Usage:", "Examples:", "Main Commands:", "OpenShift Commands:", "Utility Commands:", "Flags:"})
					helper.DontMatchAllInOutput(output, []string{"--kubeconfig"})
				})
				It("does not support the --kubeconfig flag", func() {
					helper.DontMatchAllInOutput(output, []string{"--kubeconfig"})
				})
			})

			When("running astra without subcommand and flags", func() {
				var output string
				BeforeEach(func() {
					output = helper.Cmd("astra").ShouldPass().Out()
				})
				It("a short vesion of help contents is returned, an error is not expected", func() {
					Expect(output).To(ContainSubstring("To see a full list of commands, run 'astra --help'"))
				})
			})

			It("returns error when using an invalid command", func() {
				output := helper.Cmd("astra", "hello").ShouldFail().Err()
				Expect(output).To(ContainSubstring("Invalid command - see available commands/subcommands above"))
			})

			It("returns JSON error", func() {
				By("using an invalid command with JSON output", func() {
					res := helper.Cmd("astra", "unknown-command", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
				})

				By("using an invalid describe sub-command with JSON output", func() {
					res := helper.Cmd("astra", "describe", "unknown-sub-command", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
				})

				By("using an invalid list sub-command with JSON output", func() {
					res := helper.Cmd("astra", "list", "unknown-sub-command", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
				})

				By("omitting required subcommand with JSON output", func() {
					res := helper.Cmd("astra", "describe", "-o", "json").ShouldFail()
					stdout, stderr := res.Out(), res.Err()
					Expect(stdout).To(BeEmpty())
					Expect(helper.IsJSON(stderr)).To(BeTrue())
				})
			})

			It("returns error when using an invalid command with --help", func() {
				output := helper.Cmd("astra", "hello", "--help").ShouldFail().Err()
				Expect(output).To(ContainSubstring("unknown command 'hello', type --help for a list of all commands"))
			})
		})
	}

	Context("When deleting two project one after the other", func() {
		It("should be able to delete sequentially", func() {
			project1 := helper.CreateRandProject()
			project2 := helper.CreateRandProject()

			helper.DeleteProject(project1)
			helper.DeleteProject(project2)
		})
		It("should be able to delete them in any order", func() {
			project1 := helper.CreateRandProject()
			project2 := helper.CreateRandProject()
			project3 := helper.CreateRandProject()

			helper.DeleteProject(project2)
			helper.DeleteProject(project1)
			helper.DeleteProject(project3)
		})
	})

	Context("executing astra version command", func() {
		const (
			reastraVersion        = `^astra\s*v[0-9]+.[0-9]+.[0-9]+(?:-\w+)?\s*\(\w+(-\w+)?\)`
			reKubernetesVersion = `Kubernetes:\s*v[0-9]+.[0-9]+.[0-9]+((-\w+\.[0-9]+)?\+\w+)?`
			rePodmanVersion     = `Podman Client:\s*[0-9]+.[0-9]+.[0-9]+((-\w+\.[0-9]+)?\+\w+)?`
			reJSONVersion       = `^v{0,1}[0-9]+.[0-9]+.[0-9]+((-\w+\.[0-9]+)?\+\w+)?`
		)
		When("executing the complete command with server info", func() {
			var astraVersion string
			BeforeEach(func() {
				astraVersion = helper.Cmd("astra", "version").ShouldPass().Out()
			})
			for _, podman := range []bool{true, false} {
				podman := podman
				It("should show the version of astra major components including server login URL", helper.LabelPodmanIf(podman, func() {
					By("checking the human readable output", func() {
						Expect(astraVersion).Should(MatchRegexp(reastraVersion))

						// astra tests setup (CommonBeforeEach) is designed in a way that if a test is labelled with 'podman', it will not have cluster configuration
						// so we only test podman info on podman labelled test, and clsuter info otherwise
						// Tastra (pvala): Change this behavior when we write tests that should be tested on both podman and cluster simultaneously
						// Ref: https://github\.com/danielpickens/astra/issues/6719
						if podman {
							Expect(astraVersion).Should(MatchRegexp(rePodmanVersion))
							Expect(astraVersion).To(ContainSubstring(helper.GetPodmanVersion()))
						} else {
							Expect(astraVersion).Should(MatchRegexp(reKubernetesVersion))
							if !helper.IsKubernetesCluster() {
								serverURL := oc.GetCurrentServerURL()
								Expect(astraVersion).Should(ContainSubstring("Server: " + serverURL))
								ocpMatcher := ContainSubstring("OpenShift: ")
								if serverVersion := commonVar.CliRunner.GetVersion(); serverVersion == "" {
									// Might indicate a user permission error on certain clusters (observed with a developer account on Prow nightly jobs)
									ocpMatcher = Not(ocpMatcher)
								}
								Expect(astraVersion).Should(ocpMatcher)
							}
						}
					})

					By("checking the JSON output", func() {
						astraVersion = helper.Cmd("astra", "version", "-o", "json").ShouldPass().Out()
						Expect(helper.IsJSON(astraVersion)).To(BeTrue())
						helper.JsonPathSatisfiesAll(astraVersion, "version", MatchRegexp(reJSONVersion))
						helper.JsonPathExist(astraVersion, "gitCommit")
						if podman {
							helper.JsonPathSatisfiesAll(astraVersion, "podman.client.version", MatchRegexp(reJSONVersion), Equal(helper.GetPodmanVersion()))
						} else {
							helper.JsonPathSatisfiesAll(astraVersion, "cluster.kubernetes.version", MatchRegexp(reJSONVersion))
							if !helper.IsKubernetesCluster() {
								serverURL := oc.GetCurrentServerURL()
								helper.JsonPathContentIs(astraVersion, "cluster.serverURL", serverURL)
								m := BeEmpty()
								if serverVersion := commonVar.CliRunner.GetVersion(); serverVersion != "" {
									// A blank serverVersion might indicate a user permission error on certain clusters (observed with a developer account on Prow nightly jobs)
									m = Not(m)
								}
								helper.JsonPathSatisfiesAll(astraVersion, "cluster.openshift", m)
							}
						}
					})
				}))
			}

			for _, label := range []string{helper.LabelNoCluster, helper.LabelUnauth} {
				label := label
				It("should show the version of astra major components", Label(label), func() {
					Expect(astraVersion).Should(MatchRegexp(reastraVersion))
				})
			}
		})

		When("podman client is bound to delay and astra version is run", Label(helper.LabelPodman), func() {
			var astraVersion string
			BeforeEach(func() {
				delayer := helper.GenerateDelayedPodman(commonVar.Context, 2)
				astraVersion = helper.Cmd("astra", "version").WithEnv("PODMAN_CMD="+delayer, "PODMAN_CMD_INIT_TIMEOUT=1s").ShouldPass().Out()
			})
			It("should not print podman version if podman cmd timeout has been reached", func() {
				Expect(astraVersion).Should(MatchRegexp(reastraVersion))
				Expect(astraVersion).ToNot(ContainSubstring("Podman Client:"))
			})
		})
		It("should only print client info when using --client flag", func() {
			By("checking human readable output", func() {
				astraVersion := helper.Cmd("astra", "version", "--client").ShouldPass().Out()
				Expect(astraVersion).Should(MatchRegexp(reastraVersion))
				Expect(astraVersion).ToNot(SatisfyAll(ContainSubstring("Server"), ContainSubstring("Kubernetes"), ContainSubstring("Podman Client")))
			})

			By("checking JSON output", func() {
				astraVersion := helper.Cmd("astra", "version", "--client", "-o", "json").ShouldPass().Out()
				Expect(helper.IsJSON(astraVersion)).To(BeTrue())
				helper.JsonPathSatisfiesAll(astraVersion, "version", MatchRegexp(reJSONVersion))
				helper.JsonPathExist(astraVersion, "gitCommit")
				helper.JsonPathSatisfiesAll(astraVersion, "cluster", BeEmpty())
				helper.JsonPathSatisfiesAll(astraVersion, "podman", BeEmpty())
			})
		})
	})

	Describe("Experimental Mode", Label(helper.LabelNoCluster), func() {
		AfterEach(func() {
			helper.ResetExperimentalMode()
		})

		When("experimental mode is enabled", func() {
			BeforeEach(func() {
				helper.EnableExperimentalMode()
			})

			AfterEach(func() {
				helper.ResetExperimentalMode()
			})

			It("should display warning message", func() {
				out := helper.Cmd("astra", "version", "--client").ShouldPass().Out()
				Expect(out).Should(ContainSubstring("Experimental mode enabled. Use at your own risk."))
			})
		})
	})

})
