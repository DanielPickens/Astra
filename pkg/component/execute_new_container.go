package component

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/generator"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"

	"github\.com/danielpickens/astra/pkg/configAutomount"
	"github\.com/danielpickens/astra/pkg/dev/kubedev/storage"
	"github\.com/danielpickens/astra/pkg/kclient"
	astralabels "github\.com/danielpickens/astra/pkg/labels"
	astragenerator "github\.com/danielpickens/astra/pkg/libdevfile/generator"
	"github\.com/danielpickens/astra/pkg/log"
	"github\.com/danielpickens/astra/pkg/util"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
)

func ExecuteInNewContainer(
	ctx context.Context,
	kubeClient kclient.ClientInterface,
	configAutomountClient configAutomount.Client,
	devfileObj parser.DevfileObj,
	componentName string,
	appName string,
	command v1alpha2.Command,
) error {
	policy, err := kubeClient.GetCurrentNamespacePolicy()
	if err != nil {
		return err
	}
	podTemplateSpec, err := generator.GetPodTemplateSpec(devfileObj, generator.PodTemplateParams{
		Options: common.DevfileOptions{
			FilterByName: command.Exec.Component,
		},
		PodSecurityAdmissionPolicy: policy,
	})
	if err != nil {
		return err
	}
	// Setting the restart policy to "never" so that pods are kept around after the job finishes execution; this is helpful in obtaining logs to debug.
	podTemplateSpec.Spec.RestartPolicy = "Never"

	if len(podTemplateSpec.Spec.Containers) != 1 {
		return fmt.Errorf("could not find the component")
	}

	podTemplateSpec.Spec.Containers[0].Command = []string{"/bin/sh"}
	podTemplateSpec.Spec.Containers[0].Args = getJobCmdline(command)

	volumes, err := storage.GetAutomountVolumes(configAutomountClient, podTemplateSpec.Spec.Containers, podTemplateSpec.Spec.InitContainers)
	if err != nil {
		return err
	}

	podTemplateSpec.Spec.Volumes = volumes

	// Create a Kubernetes Job and use the container image referenced by command.Exec.Component
	// Get the component for the command with command.Exec.Component
	getJobName := func() string {
		maxLen := kclient.JobNameastraMaxLength - len(command.Id)
		// We ignore the error here because our component name or app name will never be empty; which are the only cases when an error might be raised.
		name, _ := util.NamespaceKubernetesObjectWithTrim(componentName, appName, maxLen)
		name += "-" + command.Id
		return name
	}
	completionMode := batchv1.CompletionMode("Indexed")
	jobParams := astragenerator.JobParams{
		TypeMeta: generator.GetTypeMeta(kclient.JobsKind, kclient.JobsAPIVersion),
		ObjectMeta: metav1.ObjectMeta{
			Name: getJobName(),
		},
		PodTemplateSpec: *podTemplateSpec,
		SpecParams: astragenerator.JobSpecParams{
			CompletionMode:          &completionMode,
			TTLSecondsAfterFinished: pointer.Int32(60),
			BackOffLimit:            pointer.Int32(1),
		},
	}
	job := astragenerator.GetJob(jobParams)
	// Set labels and annotations
	job.SetLabels(astralabels.GetLabels(componentName, appName, GetComponentRuntimeFromDevfileMetadata(devfileObj.Data.GetMetadata()), astralabels.ComponentDeployMode, false))
	job.Annotations = map[string]string{}
	astralabels.AddCommonAnnotations(job.Annotations)
	astralabels.SetProjectType(job.Annotations, GetComponentTypeFromDevfileMetadata(devfileObj.Data.GetMetadata()))

	//	Make sure there are no existing jobs
	checkAndDeleteExistingJob := func() {
		items, dErr := kubeClient.ListJobs(astralabels.GetSelector(componentName, appName, astralabels.ComponentDeployMode, false))
		if dErr != nil {
			klog.V(4).Infof("failed to list jobs; cause: %s", dErr.Error())
			return
		}
		jobName := getJobName()
		for _, item := range items.Items {
			if strings.Contains(item.Name, jobName) {
				dErr = kubeClient.DeleteJob(item.Name)
				if dErr != nil {
					klog.V(4).Infof("failed to delete job %q; cause: %s", item.Name, dErr.Error())
				}
			}
		}
	}
	checkAndDeleteExistingJob()

	log.Sectionf("Executing command:")
	spinner := log.Spinnerf("Executing command in container (command: %s)", command.Id)
	defer spinner.End(false)

	var createdJob *batchv1.Job
	createdJob, err = kubeClient.CreateJob(job, "")
	if err != nil {
		return err
	}
	defer func() {
		err = kubeClient.DeleteJob(createdJob.Name)
		if err != nil {
			klog.V(4).Infof("failed to delete job %q; cause: %s", createdJob.Name, err)
		}
	}()

	var done = make(chan struct{}, 1)
	// Print the tip to use `astra logs` if the command is still running after 1 minute
	go func() {
		select {
		case <-time.After(1 * time.Minute):
			log.Info("\nTip: Run `astra logs --deploy --follow` to get the logs of the command output.")
		case <-done:
			return
		}
	}()

	// Wait for the command to complete execution
	_, err = kubeClient.WaitForJobToComplete(createdJob)
	done <- struct{}{}

	spinner.End(err == nil)

	if err != nil {
		err = fmt.Errorf("failed to execute (command: %s)", command.Id)
		// Print the job logs if the job failed
		jobLogs, logErr := kubeClient.GetJobLogs(createdJob, command.Exec.Component)
		if logErr != nil {
			log.Warningf("failed to fetch the logs of execution; cause: %s", logErr)
		}
		fmt.Println("Execution output:")
		_ = util.DisplayLog(false, jobLogs, log.GetStderr(), componentName, 100)
	}

	return err
}

func getJobCmdline(command v1alpha2.Command) []string {
	// deal with environment variables
	var cmdLine string
	setEnvVariable := util.GetCommandStringFromEnvs(command.Exec.Env)

	if setEnvVariable == "" {
		cmdLine = command.Exec.CommandLine
	} else {
		cmdLine = setEnvVariable + " && " + command.Exec.CommandLine
	}
	var args []string
	if command.Exec.WorkingDir != "" {
		// since we are using /bin/sh -c, the command needs to be within a single double quote instance, for example "cd /tmp && pwd"
		args = []string{"-c", "cd " + command.Exec.WorkingDir + " && " + cmdLine}
	} else {
		args = []string{"-c", cmdLine}
	}
	return args
}
