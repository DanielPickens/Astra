package integration

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github\.com/danielpickens/astra/pkg/labels"
	"github\.com/danielpickens/astra/pkg/util"
	"github\.com/danielpickens/astra/tests/helper"
)

var _ = Describe("astra delete command tests", Label(helper.LabelSkipOnOpenShift), func() {
	var commonVar helper.CommonVar
	var cmpName, deploymentName, serviceName string
	var getDeployArgs, getSVCArgs []string

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		cmpName = helper.RandString(6)
		helper.Chdir(commonVar.Context)
		getDeployArgs = []string{"get", "deployment", "-n", commonVar.Project}
		getSVCArgs = []string{"get", "svc", "-n", commonVar.Project}
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	for _, ctx := range []struct {
		title       string
		devfileName string
		setupFunc   func()
		// Tastra(pvala): Find a better solution to renaming a resource when the data is in a different location
		renameServiceFunc func(newName string)
	}{
		{
			title:       "a component is bootstrapped",
			devfileName: "devfile-deploy-with-multiple-resources.yaml",
			renameServiceFunc: func(newName string) {
				helper.ReplaceString(filepath.Join(commonVar.Context, "devfile.yaml"), fmt.Sprintf("name: %s", serviceName), fmt.Sprintf("name: %s", newName))
			},
		},
		{
			title:       "a component is bootstrapped using a devfile.yaml with URI-referenced Kubernetes components",
			devfileName: "devfile-deploy-with-multiple-resources-and-k8s-uri.yaml",
			setupFunc: func() {
				helper.CopyExample(
					filepath.Join("source", "devfiles", "nodejs", "kubernetes", "devfile-deploy-with-multiple-resources-and-k8s-uri"),
					filepath.Join(commonVar.Context, "kubernetes", "devfile-deploy-with-multiple-resources-and-k8s-uri"))
			},
			renameServiceFunc: func(newName string) {
				helper.ReplaceString(filepath.Join(commonVar.Context, "kubernetes", "devfile-deploy-with-multiple-resources-and-k8s-uri", "outerloop-deploy-2.yaml"), fmt.Sprintf("name: %s", serviceName), fmt.Sprintf("name: %s", newName))
			},
		},
	} {
		// this is a workaround to ensure that the for loop works with `It` blocks
		ctx := ctx
		When(ctx.title, func() {
			BeforeEach(func() {
				deploymentName = "my-component"
				serviceName = "my-cs"
				helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
				if ctx.setupFunc != nil {
					ctx.setupFunc()
				}
				helper.Cmd("astra", "init", "--name", cmpName, "--devfile-path",
					helper.GetExamplePath("source", "devfiles", "nodejs", ctx.devfileName)).ShouldPass()
			})

			for _, podman := range []bool{true, false} {
				podman := podman
				When("the component is deployed in DEV mode and dev mode stopped", helper.LabelPodmanIf(podman, func() {
					var devSession helper.DevSession
					BeforeEach(func() {
						var err error
						devSession, err = helper.StartDevMode(helper.DevSessionOpts{
							RunOnPodman: podman,
						})
						Expect(err).ToNot(HaveOccurred())

						devSession.Kill()
						devSession.WaitEnd()

						component := helper.NewComponent(cmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
						component.ExpectIsDeployed()
					})

					AfterEach(func() {
						args := []string{"delete", "component", "--name", cmpName}
						if !podman {
							args = append(args, "--namespace", commonVar.Project)
						}
						args = append(args, "-f", "--wait")
						helper.Cmd("astra", args...).ShouldPass()
					})

					for _, runningIn := range []string{"", "dev", "deploy"} {
						runningIn := runningIn
						When(fmt.Sprintf("the component is deleted using its name (and namespace) from another directory (running-in=%q)", runningIn), func() {
							var out string
							BeforeEach(func() {
								otherDir := filepath.Join(commonVar.Context, "tmp")
								helper.MakeDir(otherDir)
								helper.Chdir(otherDir)
								args := []string{"delete", "component", "--name", cmpName}
								if !podman {
									args = append(args, "--namespace", commonVar.Project)
								}
								if runningIn != "" {
									args = append(args, "--running-in", runningIn)
								}
								out = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
							})

							if runningIn == "deploy" {
								It("should output that there are no resources to be deleted", func() {
									Expect(out).To(ContainSubstring("No resource found for component %q", cmpName))
								})
							} else {
								It("should have deleted the component", func() {
									By("listing the resource to delete", func() {
										if podman {
											Expect(out).To(ContainSubstring("- " + cmpName))
										} else {
											Expect(out).To(ContainSubstring("Deployment: " + cmpName))
										}
									})
									By("deleting the deployment", func() {
										component := helper.NewComponent(cmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
										component.ExpectIsNotDeployed()
									})
								})
							}

							if !podman {
								When("astra delete command is run again with nothing deployed on the cluster", func() {
									var stdOut string
									BeforeEach(func() {
										// wait until the resources are deleted from the first delete
										Eventually(string(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents()), 60, 3).ShouldNot(ContainSubstring(deploymentName))
										Eventually(string(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents()), 60, 3).ShouldNot(ContainSubstring(serviceName))
									})
									It("should output that there are no resources to be deleted", func() {
										helper.CreateInvalidDevfile(commonVar.Context)
										helper.Chdir(commonVar.Context)
										args := []string{"delete", "component", "--name", cmpName, "--namespace", commonVar.Project}
										if runningIn != "" {
											args = append(args, "--running-in", runningIn)
										}
										Eventually(func() string {
											stdOut = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
											return stdOut
										}, 60, 3).Should(ContainSubstring("No resource found for component %q in namespace %q", cmpName, commonVar.Project))
									})
								})
							}
						})

						Context("the component is deleted while having access to the devfile.yaml", func() {
							When("the component is deleted with --files", func() {
								var stdOut string
								BeforeEach(func() {
									args := []string{"delete", "component", "--files"}
									if runningIn != "" {
										args = append(args, "--running-in", runningIn)
									}
									stdOut = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
								})

								It("should have deleted the component", func() {
									if runningIn == "deploy" {
										By("outputting that there are no resources to be deleted", func() {
											Expect(stdOut).To(ContainSubstring("No resource found for component %q", cmpName))
										})
									} else {
										By("listing the resources to delete", func() {
											Expect(stdOut).To(ContainSubstring(cmpName))
										})
										By("deleting the deployment", func() {
											component := helper.NewComponent(cmpName, "app", labels.ComponentDevMode, commonVar.Project, commonVar.CliRunner)
											component.ExpectIsNotDeployed()
										})
									}

									deletableFileNames := []string{util.DotastraDirectory, "devfile.yaml"}
									var deletableFilesPaths []string
									By("listing the files to delete", func() {
										for _, f := range deletableFileNames {
											deletableFilesPaths = append(deletableFilesPaths, filepath.Join(commonVar.Context, f))
										}
										helper.MatchAllInOutput(stdOut, deletableFilesPaths)
									})
									By("ensuring that appropriate files have been removed", func() {
										files := helper.ListFilesInDir(commonVar.Context)
										for _, f := range deletableFileNames {
											Expect(files).ShouldNot(ContainElement(f))
										}
									})
								})
							})
						})
					}
				}))
			}

			When("the component is deployed in DEPLOY mode", func() {
				BeforeEach(func() {
					helper.Cmd("astra", "deploy").AddEnv("PODMAN_CMD=echo").ShouldPass()
					Expect(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents()).To(ContainSubstring(deploymentName))
					Expect(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents()).To(ContainSubstring(serviceName))
				})

				for _, runningIn := range []string{"", "dev", "deploy"} {
					runningIn := runningIn
					When("the component is deleted using its name and namespace from another directory", func() {
						var out string
						BeforeEach(func() {
							otherDir := filepath.Join(commonVar.Context, "tmp")
							helper.MakeDir(otherDir)
							helper.Chdir(otherDir)
							args := []string{"delete", "component", "--name", cmpName, "--namespace", commonVar.Project}
							if runningIn != "" {
								args = append(args, "--running-in", runningIn)
							}
							out = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
						})

						if runningIn == "dev" {
							It("should output that there are no resources to be deleted", func() {
								Expect(out).To(ContainSubstring("No resource found for component %q", cmpName))
							})
						} else {
							It("should have deleted the component", func() {
								By("listing the resource to delete", func() {
									Expect(out).To(ContainSubstring("Deployment: " + deploymentName))
									Expect(out).To(ContainSubstring("Service: " + serviceName))
								})
								By("deleting the deployment", func() {
									Eventually(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(deploymentName))
								})
								By("deleting the service", func() {
									Eventually(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(serviceName))
								})
							})
						}
					})

					for _, withFiles := range []bool{true, false} {
						withFiles := withFiles
						When(fmt.Sprintf("a resource is changed in the devfile and the component is deleted while having access to the devfile.yaml with --files=%v --running-in=%v",
							withFiles, runningIn), func() {
							var changedServiceName, stdout string
							BeforeEach(func() {
								changedServiceName = "my-changed-cs"
								ctx.renameServiceFunc(changedServiceName)

								args := []string{"delete", "component"}
								if withFiles {
									args = append(args, "--files")
								}
								if runningIn != "" {
									args = append(args, "--running-in", runningIn)
								}
								stdout = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
							})
							It("should delete the component", func() {
								if runningIn == "dev" {
									By("outputting that there are no resources to be deleted", func() {
										Expect(stdout).To(ContainSubstring("No resource found for component %q", cmpName))
									})
								} else {
									By("showing warning about undeleted service belonging to the component", func() {
										Expect(stdout).To(SatisfyAll(
											ContainSubstring("There are still resources left in the cluster that might be belonging to the deleted component"),
											Not(ContainSubstring(changedServiceName)),
											ContainSubstring(fmt.Sprintf("Service: %s", serviceName)),
											ContainSubstring("astra delete component --name %s --namespace %s", cmpName, commonVar.Project),
										))
									})
								}

								files := helper.ListFilesInDir(commonVar.Context)
								if withFiles {
									By("ensuring that devfile.yaml has been removed because it was created with astra init", func() {
										Expect(files).ShouldNot(ContainElement("devfile.yaml"))
									})
								} else {
									By("ensuring that devfile.yaml still exists", func() {
										Expect(files).To(ContainElement("devfile.yaml"))
									})
								}
							})

						})

						When("the component is deleted while having access to the devfile.yaml", func() {
							var stdOut string
							BeforeEach(func() {
								args := []string{"delete", "component"}
								if withFiles {
									args = append(args, "--files")
								}
								if runningIn != "" {
									args = append(args, "--running-in", runningIn)
								}
								stdOut = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
							})
							It("should have deleted the component", func() {
								if runningIn == "dev" {
									By("outputting that there are no resources to be deleted", func() {
										Expect(stdOut).To(ContainSubstring("No resource found for component %q", cmpName))
									})
								} else {
									By("listing the resources to delete", func() {
										Expect(stdOut).To(ContainSubstring(cmpName))
										Expect(stdOut).To(ContainSubstring("Deployment: " + deploymentName))
										Expect(stdOut).To(ContainSubstring("Service: " + serviceName))
									})
									By("deleting the deployment", func() {
										Eventually(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(deploymentName))
									})
									By("deleting the service", func() {
										Eventually(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(serviceName))
									})
								}
								files := helper.ListFilesInDir(commonVar.Context)
								if withFiles {
									By("ensuring that devfile.yaml has been removed because it was created with astra init", func() {
										Expect(files).ShouldNot(ContainElement("devfile.yaml"))
									})
								} else {
									By("ensuring that devfile.yaml still exists", func() {
										Expect(files).To(ContainElement("devfile.yaml"))
									})
								}
							})
						})
					}
				}
			})

			When("the component is running in both DEV and DEPLOY mode and dev mode is killed", func() {
				var devSession helper.DevSession
				BeforeEach(func() {
					helper.Cmd("astra", "deploy").AddEnv("PODMAN_CMD=echo").ShouldPass()
					Expect(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents()).To(ContainSubstring(deploymentName))
					Expect(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents()).To(ContainSubstring(serviceName))

					var err error
					devSession, err = helper.StartDevMode(helper.DevSessionOpts{})
					Expect(err).ToNot(HaveOccurred())

					devSession.Kill()
					devSession.WaitEnd()

					component := helper.NewComponent(cmpName, "app", "runtime", commonVar.Project, commonVar.CliRunner)
					component.ExpectIsDeployed()
				})

				checkDeployResourcesDeletion := func(stdout string, checkDevNotDeleted bool) {
					By("listing the resources to delete", func() {
						Expect(stdout).To(ContainSubstring(cmpName))
						Expect(stdout).To(ContainSubstring("Deployment: " + deploymentName))
						Expect(stdout).To(ContainSubstring("Service: " + serviceName))
					})
					By("deleting the deployment", func() {
						Eventually(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(deploymentName))
					})
					By("deleting the service", func() {
						Eventually(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents(), 60, 3).ShouldNot(ContainSubstring(serviceName))
					})

					if checkDevNotDeleted {
						By("not deleting Dev resources", func() {
							Consistently(commonVar.CliRunner.Run("get", "deployments", "-n", commonVar.Project).Out.Contents()).
								WithTimeout(30 * time.Second).
								WithPolling(3 * time.Second).
								Should(ContainSubstring(cmpName + "-app"))
						})
					}
				}

				checkDevResourcesDeletion := func(stdout string, checkDeployNotDeleted bool) {
					By("listing the resources to delete", func() {
						Expect(stdout).To(ContainSubstring(cmpName))
					})
					By("deleting the deployment", func() {
						component := helper.NewComponent(cmpName, "app", "runtime", commonVar.Project, commonVar.CliRunner)
						component.ExpectIsNotDeployed()
					})

					if checkDeployNotDeleted {
						By("not deleting Deploy resources", func() {
							Consistently(func(g Gomega) {
								g.Expect(commonVar.CliRunner.Run(getDeployArgs...).Out.Contents()).Should(ContainSubstring(deploymentName))
								g.Expect(commonVar.CliRunner.Run(getSVCArgs...).Out.Contents()).Should(ContainSubstring(serviceName))
							}).WithTimeout(30 * time.Second).WithPolling(3 * time.Second).Should(Succeed())
						})
					}
				}

				for _, runningIn := range []string{"", "dev", "deploy"} {
					runningIn := runningIn
					for _, withFiles := range []bool{true, false} {
						withFiles := withFiles
						When("the component is deleted while having access to the devfile.yaml", func() {
							var stdout string

							BeforeEach(func() {
								args := []string{"delete", "component"}
								if runningIn != "" {
									args = append(args, "--running-in", runningIn)
								}
								if withFiles {
									args = append(args, "--files")
								}
								stdout = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
							})

							It("should delete the appropriate resources", func() {
								if runningIn == "" {
									By("deleting all resources if no running mode specified", func() {
										checkDevResourcesDeletion(stdout, false)
										checkDeployResourcesDeletion(stdout, false)
									})
								} else {
									By("deleting only resources running in the specified mode: "+runningIn, func() {
										switch runningIn {
										case "dev":
											checkDevResourcesDeletion(stdout, true)
										case "deploy":
											checkDeployResourcesDeletion(stdout, true)
										}
									})
								}

								files := helper.ListFilesInDir(commonVar.Context)
								if withFiles {
									By("ensuring that devfile.yaml has been removed because it was created with astra init", func() {
										Expect(files).ShouldNot(ContainElement("devfile.yaml"))
									})
								} else {
									By("ensuring that devfile.yaml still exists", func() {
										Expect(files).To(ContainElement("devfile.yaml"))
									})
								}
							})
						})
					}

					When("the component is deleted using its name and namespace from another directory", func() {
						var stdout string

						BeforeEach(func() {
							otherDir := filepath.Join(commonVar.Context, "tmp")
							helper.MakeDir(otherDir)
							helper.Chdir(otherDir)
							args := []string{"delete", "component", "--name", cmpName, "--namespace", commonVar.Project}
							if runningIn != "" {
								args = append(args, "--running-in", runningIn)
							}
							stdout = helper.Cmd("astra", append(args, "-f")...).ShouldPass().Out()
						})

						It("should delete the appropriate resources", func() {
							if runningIn == "" {
								By("deleting all resources if no running mode specified", func() {
									checkDevResourcesDeletion(stdout, false)
									checkDeployResourcesDeletion(stdout, false)
								})
							} else {
								By("deleting only resources running in the specified mode: "+runningIn, func() {
									switch runningIn {
									case "dev":
										checkDevResourcesDeletion(stdout, true)
									case "deploy":
										checkDeployResourcesDeletion(stdout, true)
									}
								})
							}
						})
					})
				}
			})
		})
	}

	for _, withFiles := range []bool{true, false} {
		withFiles := withFiles
		When("deleting a component containing preStop event that is deployed with DEV and --files="+strconv.FormatBool(withFiles), func() {
			var out string
			BeforeEach(func() {
				// Hardcoded names from devfile-with-valid-events.yaml
				cmpName = "nodejs"
				helper.CopyExample(filepath.Join("source", "devfiles", "nodejs", "project"), commonVar.Context)
				helper.Cmd("astra", "init", "--name", cmpName, "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-with-valid-events.yaml")).ShouldPass()
				session := helper.CmdRunner("astra", "dev", "--random-ports")
				defer session.Kill()
				helper.WaitForOutputToContain("[Ctrl+c] - Exit", 180, 10, session)
				// Ensure that the pod is in running state
				Eventually(string(commonVar.CliRunner.Run("get", "pods", "-n", commonVar.Project).Out.Contents()), 60, 3).Should(ContainSubstring(cmpName))
				// running in verbosity since the preStop events information is only printed in v4
				args := []string{"delete", "component", "-v", "4", "-f"}
				if withFiles {
					args = append(args, "--files")
				}
				out = helper.Cmd("astra", args...).ShouldPass().Out()
			})
			It("should delete the component", func() {
				By("listing preStop events", func() {
					for _, cmdName := range []string{"myprestop", "secondprestop", "thirdprestop"} {
						Expect(out).To(ContainSubstring("Executing pre-stop command in container (command: %s)", cmdName))
					}
				})
				files := helper.ListFilesInDir(commonVar.Context)
				if withFiles {
					By("ensuring that appropriate files have been removed", func() {
						Expect(files).ShouldNot(ContainElement("devfile.yaml"))
					})
				} else {
					By("ensuring that devfile.yaml still exists", func() {
						Expect(files).To(ContainElement("devfile.yaml"))
					})
				}
			})
		})
	}
	When("running astra deploy for an exec command bound to fail", func() {
		BeforeEach(func() {
			helper.CopyExampleDevFile(
				filepath.Join("source", "devfiles", "nodejs", "devfile-deploy-exec.yaml"),
				path.Join(commonVar.Context, "devfile.yaml"),
				cmpName)
			helper.ReplaceString(filepath.Join(commonVar.Context, "devfile.yaml"), `image: registry.access.redhat.com/ubi8/nodejs-14:latest`, `image: registry.access.redhat.com/ubi8/nodejs-does-not-exist-14:latest`)
			// We terminate after 5 seconds because the job should have been created by then and is bound to fail.
			helper.Cmd("astra", "deploy").WithTerminate(5, nil).ShouldRun()
		})
		It("should print the job in the list of resources to be deleted with named delete command", func() {
			out := helper.Cmd("astra", "delete", "component", "-f").ShouldPass().Out()
			Expect(out).To(SatisfyAll(
				ContainSubstring("There are still resources left in the cluster that might be belonging to the deleted component."),
				ContainSubstring(fmt.Sprintf("Job: %s-app-deploy-exec", cmpName))))
		})
	})
})
