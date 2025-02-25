package integration

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/utils/pointer"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"

	"github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/tests/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("astra dev debug command tests", func() {
	var cmpName string
	var commonVar helper.CommonVar

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		cmpName = helper.RandString(6)
		helper.Chdir(commonVar.Context)
		Expect(helper.VerifyFileExists(".astra/env/env.yaml")).To(BeFalse())
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	for _, podman := range []bool{false, true} {
		podman := podman
		When("a component is bootstrapped", Label(helper.LabelSkipOnOpenShift), func() {
			BeforeEach(func() {
				helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
				helper.Cmd("astra", "init", "--name", cmpName, "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-debugrun.yaml")).ShouldPass()
				Expect(helper.VerifyFileExists(".astra/env/env.yaml")).To(BeFalse())
			})
			When("running astra dev with debug flag and custom port mapping for port forwarding", helper.LabelPodmanIf(podman, func() {
				var (
					devSession                helper.DevSession
					LocalPort, LocalDebugPort int
				)
				const (
					ContainerPort      = "3000"
					ContainerDebugPort = "5858"
				)

				BeforeEach(func() {
					LocalPort = helper.GetCustomStartPort()
					LocalDebugPort = LocalPort + 1
					opts := []string{"--debug", fmt.Sprintf("--port-forward=%d:%s", LocalPort, ContainerPort), fmt.Sprintf("--port-forward=%d:%s", LocalDebugPort, ContainerDebugPort)}
					if podman {
						opts = append(opts, "--forward-localhost")
					}
					var err error
					devSession, err = helper.StartDevMode(helper.DevSessionOpts{
						CmdlineArgs:   opts,
						NoRandomPorts: true,
						RunOnPodman:   podman,
					})
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})

				It("should connect to relevant custom ports forwarded", func() {
					By("connecting to the application port", func() {
						helper.HttpWaitForWithStatus(fmt.Sprintf("http://%s", devSession.Endpoints[ContainerPort]), "Hello from Node.js Starter Application!", 12, 5, 200)
					})
					By("expecting a ws connection when tried to connect on default debug port locally", func() {
						// 400 response expected because the endpoint expects a websocket request and we are doing a HTTP GET
						// We are just using this to validate if nodejs agent is listening on the other side
						url := fmt.Sprintf("http://%s", devSession.Endpoints[ContainerDebugPort])
						Expect(url).To(ContainSubstring(strconv.Itoa(LocalDebugPort)))

						helper.HttpWaitForWithStatus(url, "WebSockets request was expected", 12, 5, 400)
					})
				})
			}))

			When("running astra dev with debug flag", helper.LabelPodmanIf(podman, func() {
				var devSession helper.DevSession

				BeforeEach(func() {
					var err error
					opts := helper.DevSessionOpts{
						CmdlineArgs: []string{"--debug"},
						RunOnPodman: podman,
					}
					if podman {
						opts.CmdlineArgs = append(opts.CmdlineArgs, "--forward-localhost")
					}
					devSession, err = helper.StartDevMode(opts)
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})

				It("should connect to relevant ports forwarded", func() {
					By("connecting to the application port", func() {
						helper.HttpWaitForWithStatus("http://"+devSession.Endpoints["3000"], "Hello from Node.js Starter Application!", 12, 5, 200)
					})
					By("expecting a ws connection when tried to connect on default debug port locally", func() {
						// 400 response expected because the endpoint expects a websocket request and we are doing a HTTP GET
						// We are just using this to validate if nodejs agent is listening on the other side
						helper.HttpWaitForWithStatus("http://"+devSession.Endpoints["5858"], "WebSockets request was expected", 12, 5, 400)
					})
				})

				// #6056
				It("should not add a DEBUG_PORT variable to the container", func() {
					cmp := helper.NewComponent(cmpName, "app", "runtime", commonVar.Project, commonVar.CliRunner)
					stdout, _ := cmp.Exec("runtime", []string{"sh", "-c", "echo -n ${DEBUG_PORT}"}, pointer.Bool(true))
					Expect(stdout).To(BeEmpty())
				})
			}))
		})
	}

	for _, podman := range []bool{false, true} {
		podman := podman
		When("creating nodejs component, doing astra dev and run command has dev.astra.push.path attribute", Label(helper.LabelSkipOnOpenShift), helper.LabelPodmanIf(podman, func() {
			var devSession helper.DevSession
			var devStarted bool
			BeforeEach(func() {
				helper.Cmd("astra", "init", "--name", cmpName, "--devfile-path",
					helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-remote-attributes.yaml")).ShouldPass()
				helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)

				// create a folder and file which shouldn't be pushed
				helper.MakeDir(filepath.Join(commonVar.Context, "views"))
				_, _ = helper.CreateSimpleFile(filepath.Join(commonVar.Context, "views"), "view", ".html")

				helper.ReplaceString("package.json", "node server.js", "node server-debug/server.js")
				var err error
				devSession, err = helper.StartDevMode(helper.DevSessionOpts{
					RunOnPodman: podman,
					CmdlineArgs: []string{"--debug"},
				})
				Expect(err).ToNot(HaveOccurred())
				devStarted = true
			})
			AfterEach(func() {
				if devStarted {
					devSession.Stop()
					devSession.WaitEnd()
				}
			})

			It("should sync only the mentioned files at the appropriate remote destination", func() {
				component := helper.NewComponent(cmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
				stdOut, _ := component.Exec("runtime", []string{"ls", "-lai", "/projects"}, pointer.Bool(true))

				helper.MatchAllInOutput(stdOut, []string{"package.json", "server-debug"})
				helper.DontMatchAllInOutput(stdOut, []string{"test", "views", "devfile.yaml"})

				stdOut, _ = component.Exec("runtime", []string{"ls", "-lai", "/projects/server-debug"}, pointer.Bool(true))
				helper.MatchAllInOutput(stdOut, []string{"server.js", "test"})
			})
		}))
	}

	for _, devfileHandlerCtx := range []struct {
		name          string
		sourceHandler func(path string, originalCmpName string)
	}{
		{
			name: "with metadata.name",
		},
		{
			name: "without metadata.name",
			sourceHandler: func(path string, originalCmpName string) {
				helper.UpdateDevfileContent(filepath.Join(path, "devfile.yaml"), []helper.DevfileUpdater{helper.DevfileMetadataNameRemover})
				helper.ReplaceString(filepath.Join(path, "package.json"), "nodejs-starter", originalCmpName)
			},
		},
	} {
		devfileHandlerCtx := devfileHandlerCtx
		for _, podman := range []bool{false, true} {
			podman := podman
			When("a composite command is used as debug command - "+devfileHandlerCtx.name, Label(helper.LabelSkipOnOpenShift), helper.LabelPodmanIf(podman, func() {
				var devfileCmpName string
				var devSession helper.DevSession

				BeforeEach(func() {
					devfileCmpName = helper.RandString(6)
					helper.CopyExampleDevFile(
						filepath.Join("source", "devfiles", "nodejs", "devfileCompositeRunAndDebug.yaml"),
						filepath.Join(commonVar.Context, "devfile.yaml"),
						devfileCmpName)
					helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
					if devfileHandlerCtx.sourceHandler != nil {
						devfileHandlerCtx.sourceHandler(commonVar.Context, devfileCmpName)
					}
					var err error
					opts := helper.DevSessionOpts{
						RunOnPodman: podman,
						CmdlineArgs: []string{"--debug"},
					}
					if podman {
						opts.CmdlineArgs = append(opts.CmdlineArgs, "--forward-localhost")
					}
					devSession, err = helper.StartDevMode(opts)
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})

				It("should run successfully", func() {
					By("verifying from the output that all commands have been executed", func() {
						helper.MatchAllInOutput(devSession.StdOut, []string{
							"Building your application in container",
							"Executing the application (command: mkdir)",
							"Executing the application (command: echo)",
							"Executing the application (command: install)",
							"Executing the application (command: start-debug)",
						})
					})

					By("verifying that any command that did not succeed in the middle has logged such information correctly", func() {
						helper.MatchAllInOutput(devSession.ErrOut, []string{
							"Devfile command \"echo\" exited with an error status",
							"intentional-error-message",
						})
					})

					By("building the application only once", func() {
						// Because of the Spinner, the "Building your application in container" is printed twice in the captured stdout.
						// The bracket allows to match the last occurrence with the command execution timing information.
						Expect(strings.Count(devSession.StdOut, "Building your application in container (command: install) [")).
							To(BeNumerically("==", 1), "\nOUTPUT: "+devSession.StdOut+"\n")
					})

					By("verifying that the command did run successfully", func() {
						// Verify the command executed successfully
						cmp := helper.NewComponent(devfileCmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
						out, _ := cmp.Exec("runtime", []string{"stat", "/projects/testfolder"}, pointer.Bool(true))
						Expect(out).To(ContainSubstring("/projects/testfolder"))
					})

					By("expecting a ws connection when tried to connect on default debug port locally", func() {
						// 400 response expected because the endpoint expects a websocket request and we are doing a HTTP GET
						// We are just using this to validate if nodejs agent is listening on the other side
						helper.HttpWaitForWithStatus("http://"+devSession.Endpoints["5858"], "WebSockets request was expected", 12, 5, 400)
					})
				})
			}))
		}
	}

	When("a composite apply command is used as debug command", Label(helper.LabelSkipOnOpenShift), func() {
		deploymentNames := []string{"my-openshift-component", "my-k8s-component"}
		var devSession helper.DevSession

		const (
			DEVFILE_DEBUG_PORT = "5858"
		)

		BeforeEach(func() {
			helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
			helper.CopyExampleDevFile(
				filepath.Join("source", "devfiles", "nodejs", "devfile-composite-apply-commands.yaml"),
				filepath.Join(commonVar.Context, "devfile.yaml"),
				cmpName)
			var err error
			devSession, err = helper.StartDevMode(helper.DevSessionOpts{
				EnvVars:     []string{"PODMAN_CMD=echo"},
				CmdlineArgs: []string{"--debug"},
			})
			Expect(err).ToNot(HaveOccurred())
		})
		It("should execute the composite apply commands successfully", func() {
			checkDeploymentExists := func() {
				out := commonVar.CliRunner.Run("get", "deployments").Out.Contents()
				helper.MatchAllInOutput(string(out), deploymentNames)
			}
			checkImageBuilt := func() {
				Expect(devSession.StdOut).To(ContainSubstring("Building & Pushing Image"))
				Expect(devSession.StdOut).To(ContainSubstring("build -t quay.io/unknown-account/myimage -f " + filepath.Join(commonVar.Context, "Dockerfile ") + commonVar.Context))
				Expect(devSession.StdOut).To(ContainSubstring("push quay.io/unknown-account/myimage"))
			}

			checkWSConnection := func() {
				// 400 response expected because the endpoint expects a websocket request and we are doing a HTTP GET
				// We are just using this to validate if nodejs agent is listening on the other side
				helper.HttpWaitForWithStatus("http://"+devSession.Endpoints[DEVFILE_DEBUG_PORT], "WebSockets request was expected", 12, 5, 400)
			}
			By("expecting a ws connection when tried to connect on default debug port locally", func() {
				checkWSConnection()
			})

			By("checking is the image was successfully built", func() {
				checkImageBuilt()
			})

			By("checking the deployment was created successfully", func() {
				checkDeploymentExists()
			})

			By("checking astra dev watches correctly", func() {
				// making changes to the project again
				helper.ReplaceString(filepath.Join(commonVar.Context, "server.js"), "from Node.js Starter Application", "from the new Node.js Starter Application")
				err := devSession.WaitSync()
				Expect(err).ToNot(HaveOccurred())
				checkDeploymentExists()
				checkImageBuilt()
				checkWSConnection()
			})

			By("cleaning up the resources on ending the session", func() {
				devSession.Stop()
				devSession.WaitEnd()
				out := commonVar.CliRunner.Run("get", "deployments").Out.Contents()
				helper.DontMatchAllInOutput(string(out), deploymentNames)
			})
		})
	})

	for _, devfileHandlerCtx := range []struct {
		name          string
		sourceHandler func(path string, originalCmpName string)
	}{
		{
			name: "with metadata.name",
		},
		{
			name: "without metadata.name",
			sourceHandler: func(path string, originalCmpName string) {
				helper.UpdateDevfileContent(filepath.Join(path, "devfile.yaml"), []helper.DevfileUpdater{helper.DevfileMetadataNameRemover})
				helper.ReplaceString(filepath.Join(path, "package.json"), "nodejs-starter", originalCmpName)
			},
		},
	} {
		devfileHandlerCtx := devfileHandlerCtx
		for _, podman := range []bool{false, true} {
			podman := podman
			When("running build and debug commands as composite in different containers and a shared volume - "+devfileHandlerCtx.name, Label(helper.LabelSkipOnOpenShift), helper.LabelPodmanIf(podman, func() {
				var devfileCmpName string
				var devSession helper.DevSession

				BeforeEach(func() {
					// Tastra(rm3l): For some reason, this does not work on Podman
					if podman {
						Skip("Does not work on Podman due to permission issues related in the volume mount path: /bin/sh: /artifacts/build-result: Permission denied")
					}
					devfileCmpName = helper.RandString(6)
					helper.CopyExampleDevFile(
						filepath.Join("source", "devfiles", "nodejs", "devfileCompositeBuildRunDebugInMultiContainersAndSharedVolume.yaml"),
						filepath.Join(commonVar.Context, "devfile.yaml"),
						devfileCmpName)
					helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
					if devfileHandlerCtx.sourceHandler != nil {
						devfileHandlerCtx.sourceHandler(commonVar.Context, devfileCmpName)
					}
					var err error
					opts := helper.DevSessionOpts{
						RunOnPodman: podman,
						CmdlineArgs: []string{"--debug"},
					}
					if podman {
						opts.CmdlineArgs = append(opts.CmdlineArgs, "--forward-localhost")
					}
					devSession, err = helper.StartDevMode(opts)
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					devSession.Stop()
					devSession.WaitEnd()
				})

				It("should run successfully", func() {
					By("verifying from the output that all commands have been executed", func() {
						helper.MatchAllInOutput(devSession.StdOut, []string{
							"Building your application in container (command: mkdir)",
							"Building your application in container (command: sleep-cmd-build)",
							"Building your application in container (command: build-cmd)",
							"Executing the application (command: sleep-cmd-run)",
							"Executing the application (command: echo-with-error)",
							"Executing the application (command: check-build-result)",
							"Executing the application (command: start-debug)",
						})
					})

					By("verifying that any command that did not succeed in the middle has logged such information correctly", func() {
						helper.MatchAllInOutput(devSession.ErrOut, []string{
							"Devfile command \"echo-with-error\" exited with an error status",
							"intentional-error-message",
						})
					})

					By("building the application only once per exec command in the build command", func() {
						// Because of the Spinner, the "Building your application in container" is printed twice in the captured stdout.
						// The bracket allows to match the last occurrence with the command execution timing information.
						out := devSession.StdOut
						for _, cmd := range []string{"mkdir", "sleep-cmd-build", "build-cmd"} {
							Expect(strings.Count(out, fmt.Sprintf("Building your application in container (command: %s) [", cmd))).
								To(BeNumerically("==", 1), "\nOUTPUT: "+devSession.StdOut+"\n")
						}
					})

					By("verifying that the command did run successfully", func() {
						// Verify the command executed successfully
						cmp := helper.NewComponent(devfileCmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
						out, _ := cmp.Exec("runtime", []string{"stat", "/projects/testfolder"}, pointer.Bool(true))
						Expect(out).To(ContainSubstring("/projects/testfolder"))
					})

					By("expecting a ws connection when tried to connect on default debug port locally", func() {
						// 400 response expected because the endpoint expects a websocket request and we are doing a HTTP GET
						// We are just using this to validate if nodejs agent is listening on the other side
						helper.HttpWaitForWithStatus("http://"+devSession.Endpoints["5858"], "WebSockets request was expected", 12, 5, 400)
					})
				})
			}))
		}
	}

	for _, podman := range []bool{false, true} {
		podman := podman
		When("a component without debug command is bootstrapped", Label(helper.LabelSkipOnOpenShift), helper.LabelPodmanIf(podman, func() {
			BeforeEach(func() {
				helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
				helper.Cmd("astra", "init", "--name", cmpName, "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-without-debugrun.yaml")).ShouldPass()
				Expect(helper.VerifyFileExists(".astra/env/env.yaml")).To(BeFalse())
			})

			It("should log error about missing debug command when running astra dev --debug", func() {
				devSession, err := helper.StartDevMode(helper.DevSessionOpts{
					RunOnPodman: podman,
					CmdlineArgs: []string{"--debug"},
				})
				Expect(err).ShouldNot(HaveOccurred())
				defer func() {
					devSession.Stop()
					devSession.WaitEnd()
				}()
				Expect(devSession.ErrOut).To(ContainSubstring("Missing default debug command"))
			})
		}))
	}

	// More details on https://github.com/devfile/api/issues/852#issuecomment-1211928487
	for _, podman := range []bool{false, true} {
		podman := podman
		When("starting with Devfile with autoBuild or deployByDefault components",
			Label(helper.LabelSkipOnOpenShift), // No need to repeat this test on OCP, as it is already covered on vanilla K8s
			helper.LabelPodmanIf(podman, func() {
				BeforeEach(func() {
					helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
					helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", "devfile-autobuild-deploybydefault.yaml"),
						filepath.Join(commonVar.Context, "devfile.yaml"),
						cmpName)
				})

				When("running astra dev with some components not referenced in the Devfile", func() {
					var devSession helper.DevSession

					BeforeEach(func() {
						var err error
						var envvars []string
						if podman {
							envvars = append(envvars, "astra_PUSH_IMAGES=false")
						} else {
							envvars = append(envvars, "PODMAN_CMD=echo")
						}
						args := []string{"--debug"}
						if podman {
							args = append(args, "--forward-localhost")
						}
						devSession, err = helper.StartDevMode(helper.DevSessionOpts{
							CmdlineArgs: args,
							EnvVars:     envvars,
							RunOnPodman: podman,
						})
						Expect(err).ShouldNot(HaveOccurred())
					})

					AfterEach(func() {
						devSession.Stop()
						if podman {
							devSession.WaitEnd()
						}
					})

					It("should create the appropriate resources", func() {
						if podman {
							k8sOcComponents := helper.ExtractK8sAndOcComponentsFromOutputOnPodman(devSession.ErrOut)
							By("handling Kubernetes/OpenShift components that would have been created automatically", func() {
								Expect(k8sOcComponents).Should(ContainElements(
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-true-and-not-referenced",
									"k8s-deploybydefault-not-set-and-not-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-true-and-not-referenced",
									"ocp-deploybydefault-not-set-and-not-referenced",
								))
							})
							By("not handling Kubernetes/OpenShift components with deployByDefault=false", func() {
								Expect(k8sOcComponents).ShouldNot(ContainElements(
									"k8s-deploybydefault-false-and-referenced",
									"k8s-deploybydefault-false-and-not-referenced",
									"ocp-deploybydefault-false-and-referenced",
									"ocp-deploybydefault-false-and-not-referenced",
								))
							})
							By("not handling referenced Kubernetes/OpenShift components with deployByDefault unset", func() {
								Expect(k8sOcComponents).ShouldNot(ContainElement("k8s-deploybydefault-not-set-and-referenced"))
							})
						} else {
							By("automatically applying Kubernetes/OpenShift components with deployByDefault=true", func() {
								for _, l := range []string{
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-true-and-not-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-true-and-not-referenced",
								} {
									Expect(devSession.StdOut).Should(ContainSubstring("Creating resource Pod/%s", l))
								}
							})
							By("automatically applying non-referenced Kubernetes/OpenShift components with deployByDefault not set", func() {
								for _, l := range []string{
									"k8s-deploybydefault-not-set-and-not-referenced",
									"ocp-deploybydefault-not-set-and-not-referenced",
								} {
									Expect(devSession.StdOut).Should(ContainSubstring("Creating resource Pod/%s", l))
								}
							})
							By("not applying Kubernetes/OpenShift components with deployByDefault=false", func() {
								for _, l := range []string{
									"k8s-deploybydefault-false-and-referenced",
									"k8s-deploybydefault-false-and-not-referenced",
									"ocp-deploybydefault-false-and-referenced",
									"ocp-deploybydefault-false-and-not-referenced",
								} {
									Expect(devSession.StdOut).ShouldNot(ContainSubstring("Creating resource Pod/%s", l))
								}
							})
							By("not applying referenced Kubernetes/OpenShift components with deployByDefault unset", func() {
								Expect(devSession.StdOut).ShouldNot(ContainSubstring("Creating resource Pod/k8s-deploybydefault-not-set-and-referenced"))
							})
						}

						imageMessagePrefix := "Building & Pushing Image"
						if podman {
							imageMessagePrefix = "Building Image"
						}

						By("automatically applying image components with autoBuild=true", func() {
							for _, tag := range []string{
								"autobuild-true-and-referenced",
								"autobuild-true-and-not-referenced",
							} {
								Expect(devSession.StdOut).Should(ContainSubstring("%s: localhost:5000/astra-dev/node:%s", imageMessagePrefix, tag))
							}
						})
						By("automatically applying non-referenced Image components with autoBuild not set", func() {
							Expect(devSession.StdOut).Should(ContainSubstring("%s: localhost:5000/astra-dev/node:autobuild-not-set-and-not-referenced", imageMessagePrefix))
						})
						By("not applying image components with autoBuild=false", func() {
							for _, tag := range []string{
								"autobuild-false-and-referenced",
								"autobuild-false-and-not-referenced",
							} {
								Expect(devSession.StdOut).ShouldNot(ContainSubstring("localhost:5000/astra-dev/node:%s", tag))
							}
						})
						By("not applying referenced Image components with deployByDefault unset", func() {
							Expect(devSession.StdOut).ShouldNot(ContainSubstring("localhost:5000/astra-dev/node:autobuild-not-set-and-referenced"))
						})
					})
				})

				When("running astra dev with some components referenced in the Devfile", func() {
					var devSession helper.DevSession

					BeforeEach(func() {
						var err error
						//Tastra (rm3l): we do not support passing a custom debug command yet. That's why we are manually updating the Devfile to change the default debug command.
						helper.UpdateDevfileContent(filepath.Join(commonVar.Context, "devfile.yaml"), []helper.DevfileUpdater{
							helper.DevfileCommandGroupUpdater("start-app-debug", v1alpha2.ExecCommandType, &v1alpha2.CommandGroup{
								Kind:      v1alpha2.DebugCommandGroupKind,
								IsDefault: pointer.Bool(false),
							}),
							helper.DevfileCommandGroupUpdater("debug-with-referenced-components", v1alpha2.CompositeCommandType, &v1alpha2.CommandGroup{
								Kind:      v1alpha2.DebugCommandGroupKind,
								IsDefault: pointer.Bool(true),
							}),
						})

						var envvars []string
						if podman {
							envvars = append(envvars, "astra_PUSH_IMAGES=false")
						} else {
							envvars = append(envvars, "PODMAN_CMD=echo")
						}
						args := []string{"--debug"}
						if podman {
							args = append(args, "--forward-localhost")
						}
						devSession, err = helper.StartDevMode(helper.DevSessionOpts{
							CmdlineArgs: args,
							EnvVars:     envvars,
							RunOnPodman: podman,
						})
						Expect(err).ShouldNot(HaveOccurred())
					})

					AfterEach(func() {
						devSession.Stop()
						if podman {
							devSession.WaitEnd()
						}
					})

					It("should create the appropriate resources", func() {
						if podman {
							k8sOcComponents := helper.ExtractK8sAndOcComponentsFromOutputOnPodman(devSession.ErrOut)
							By("handling Kubernetes/OpenShift components to create automatically", func() {
								Expect(k8sOcComponents).Should(ContainElements(
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-true-and-not-referenced",
									"k8s-deploybydefault-not-set-and-not-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-true-and-not-referenced",
									"ocp-deploybydefault-not-set-and-not-referenced",
								))
							})

							By("handling referenced Kubernetes/OpenShift components", func() {
								Expect(k8sOcComponents).Should(ContainElements(
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-false-and-referenced",
									"k8s-deploybydefault-not-set-and-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-false-and-referenced",
									"ocp-deploybydefault-not-set-and-referenced",
								))
							})

							By("not handling non-referenced Kubernetes/OpenShift components with deployByDefault=false", func() {
								Expect(k8sOcComponents).ShouldNot(ContainElements(
									"k8s-deploybydefault-false-and-not-referenced",
									"ocp-deploybydefault-false-and-not-referenced",
								))
							})
						} else {
							By("applying referenced Kubernetes/OpenShift components", func() {
								for _, l := range []string{
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-false-and-referenced",
									"k8s-deploybydefault-not-set-and-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-false-and-referenced",
									"ocp-deploybydefault-not-set-and-referenced",
								} {
									Expect(devSession.StdOut).Should(ContainSubstring("Creating resource Pod/%s", l))
								}
							})

							By("automatically applying Kubernetes/OpenShift components with deployByDefault=true", func() {
								for _, l := range []string{
									"k8s-deploybydefault-true-and-referenced",
									"k8s-deploybydefault-true-and-not-referenced",
									"ocp-deploybydefault-true-and-referenced",
									"ocp-deploybydefault-true-and-not-referenced",
								} {
									Expect(devSession.StdOut).Should(ContainSubstring("Creating resource Pod/%s", l))
								}
							})
							By("automatically applying non-referenced Kubernetes/OpenShift components with deployByDefault not set", func() {
								for _, l := range []string{
									"k8s-deploybydefault-not-set-and-not-referenced",
									"ocp-deploybydefault-not-set-and-not-referenced",
								} {
									Expect(devSession.StdOut).Should(ContainSubstring("Creating resource Pod/%s", l))
								}
							})

							By("not applying non-referenced Kubernetes/OpenShift components with deployByDefault=false", func() {
								for _, l := range []string{
									"k8s-deploybydefault-false-and-not-referenced",
									"ocp-deploybydefault-false-and-not-referenced",
								} {
									Expect(devSession.StdOut).ShouldNot(ContainSubstring("Creating resource Pod/%s", l))
								}
							})
						}

						imageMessagePrefix := "Building & Pushing Image"
						if podman {
							imageMessagePrefix = "Building Image"
						}

						By("applying referenced image components", func() {
							for _, tag := range []string{
								"autobuild-true-and-referenced",
								"autobuild-false-and-referenced",
								"autobuild-not-set-and-referenced",
							} {
								Expect(devSession.StdOut).Should(ContainSubstring("%s: localhost:5000/astra-dev/node:%s", imageMessagePrefix, tag))
							}
						})
						By("automatically applying image components with autoBuild=true", func() {
							for _, tag := range []string{
								"autobuild-true-and-referenced",
								"autobuild-true-and-not-referenced",
							} {
								Expect(devSession.StdOut).Should(ContainSubstring("%s: localhost:5000/astra-dev/node:%s", imageMessagePrefix, tag))
							}
						})
						By("automatically applying non-referenced Image components with autoBuild not set", func() {
							Expect(devSession.StdOut).Should(ContainSubstring("%s: localhost:5000/astra-dev/node:autobuild-not-set-and-not-referenced", imageMessagePrefix))
						})
						By("not applying non-referenced image components with autoBuild=false", func() {
							Expect(devSession.StdOut).ShouldNot(ContainSubstring("localhost:5000/astra-dev/node:autobuild-false-and-not-referenced"))
						})
					})
				})

			}))
	}
})
