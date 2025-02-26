package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// Tastras are ignored from the parser (e.g. Tastra(andronat):... || Tastra:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_CustomDeploymentStrategyParams = map[string]string{
	"":            "CustomDeploymentStrategyParams are the input to the Custom deployment strategy.",
	"image":       "Image specifies a container image which can carry out a deployment.",
	"environment": "Environment holds the environment which will be given to the container for Image.",
	"command":     "Command is optional and overrides CMD in the container Image.",
}

func (CustomDeploymentStrategyParams) SwaggerDoc() map[string]string {
	return map_CustomDeploymentStrategyParams
}

var map_DeploymentCause = map[string]string{
	"":             "DeploymentCause captures information about a particular cause of a deployment.",
	"type":         "Type of the trigger that resulted in the creation of a new deployment",
	"imageTrigger": "ImageTrigger contains the image trigger details, if this trigger was fired based on an image change",
}

func (DeploymentCause) SwaggerDoc() map[string]string {
	return map_DeploymentCause
}

var map_DeploymentCauseImageTrigger = map[string]string{
	"":     "DeploymentCauseImageTrigger represents details about the cause of a deployment originating from an image change trigger",
	"from": "From is a reference to the changed object which triggered a deployment. The field may have the kinds DockerImage, ImageStreamTag, or ImageStreamImage.",
}

func (DeploymentCauseImageTrigger) SwaggerDoc() map[string]string {
	return map_DeploymentCauseImageTrigger
}

var map_DeploymentCondition = map[string]string{
	"":                   "DeploymentCondition describes the state of a deployment config at a certain point.",
	"type":               "Type of deployment condition.",
	"status":             "Status of the condition, one of True, False, Unknown.",
	"lastUpdateTime":     "The last time this condition was updated.",
	"lastTransitionTime": "The last time the condition transitioned from one status to another.",
	"reason":             "The reason for the condition's last transition.",
	"message":            "A human readable message indicating details about the transition.",
}

func (DeploymentCondition) SwaggerDoc() map[string]string {
	return map_DeploymentCondition
}

var map_DeploymentConfig = map[string]string{
	"":       "Deployment Configs define the template for a pod and manages deploying new images or configuration changes. A single deployment configuration is usually analogous to a single micro-service. Can support many different deployment patterns, including full restart, customizable rolling updates, and  fully custom behaviors, as well as pre- and post- deployment hooks. Each individual deployment is represented as a replication controller.\n\nA deployment is \"triggered\" when its configuration is changed or a tag in an Image Stream is changed. Triggers can be disabled to allow manual control over a deployment. The \"strategy\" determines how the deployment is carried out and may be changed at any time. The `latestVersion` field is updated when a new deployment is triggered by any means.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"spec":   "Spec represents a desired deployment state and how to deploy to it.",
	"status": "Status represents the current deployment state.",
}

func (DeploymentConfig) SwaggerDoc() map[string]string {
	return map_DeploymentConfig
}

var map_DeploymentConfigList = map[string]string{
	"":      "DeploymentConfigList is a collection of deployment configs.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is a list of deployment configs",
}

func (DeploymentConfigList) SwaggerDoc() map[string]string {
	return map_DeploymentConfigList
}

var map_DeploymentConfigRollback = map[string]string{
	"":                   "DeploymentConfigRollback provides the input to rollback generation.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"name":               "Name of the deployment config that will be rolled back.",
	"updatedAnnotations": "UpdatedAnnotations is a set of new annotations that will be added in the deployment config.",
	"spec":               "Spec defines the options to rollback generation.",
}

func (DeploymentConfigRollback) SwaggerDoc() map[string]string {
	return map_DeploymentConfigRollback
}

var map_DeploymentConfigRollbackSpec = map[string]string{
	"":                       "DeploymentConfigRollbackSpec represents the options for rollback generation.",
	"from":                   "From points to a ReplicationController which is a deployment.",
	"revision":               "Revision to rollback to. If set to 0, rollback to the last revision.",
	"includeTriggers":        "IncludeTriggers specifies whether to include config Triggers.",
	"includeTemplate":        "IncludeTemplate specifies whether to include the PodTemplateSpec.",
	"includeReplicationMeta": "IncludeReplicationMeta specifies whether to include the replica count and selector.",
	"includeStrategy":        "IncludeStrategy specifies whether to include the deployment Strategy.",
}

func (DeploymentConfigRollbackSpec) SwaggerDoc() map[string]string {
	return map_DeploymentConfigRollbackSpec
}

var map_DeploymentConfigSpec = map[string]string{
	"":                     "DeploymentConfigSpec represents the desired state of the deployment.",
	"strategy":             "Strategy describes how a deployment is executed.",
	"minReadySeconds":      "MinReadySeconds is the minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)",
	"triggers":             "Triggers determine how updates to a DeploymentConfig result in new deployments. If no triggers are defined, a new deployment can only occur as a result of an explicit client update to the DeploymentConfig with a new LatestVersion. If null, defaults to having a config change trigger.",
	"replicas":             "Replicas is the number of desired replicas.",
	"revisionHistoryLimit": "RevisionHistoryLimit is the number of old ReplicationControllers to retain to allow for rollbacks. This field is a pointer to allow for differentiation between an explicit zero and not specified. Defaults to 10. (This only applies to DeploymentConfigs created via the new group API resource, not the legacy resource.)",
	"test":                 "Test ensures that this deployment config will have zero replicas except while a deployment is running. This allows the deployment config to be used as a continuous deployment test - triggering on images, running the deployment, and then succeeding or failing. Post strategy hooks and After actions can be used to integrate successful deployment with an action.",
	"paused":               "Paused indicates that the deployment config is paused resulting in no new deployments on template changes or changes in the template caused by other triggers.",
	"selector":             "Selector is a label query over pods that should match the Replicas count.",
	"template":             "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
}

func (DeploymentConfigSpec) SwaggerDoc() map[string]string {
	return map_DeploymentConfigSpec
}

var map_DeploymentConfigStatus = map[string]string{
	"":                    "DeploymentConfigStatus represents the current deployment state.",
	"latestVersion":       "LatestVersion is used to determine whether the current deployment associated with a deployment config is out of sync.",
	"observedGeneration":  "ObservedGeneration is the most recent generation observed by the deployment config controller.",
	"replicas":            "Replicas is the total number of pods targeted by this deployment config.",
	"updatedReplicas":     "UpdatedReplicas is the total number of non-terminated pods targeted by this deployment config that have the desired template spec.",
	"availableReplicas":   "AvailableReplicas is the total number of available pods targeted by this deployment config.",
	"unavailableReplicas": "UnavailableReplicas is the total number of unavailable pods targeted by this deployment config.",
	"details":             "Details are the reasons for the update to this deployment config. This could be based on a change made by the user or caused by an automatic trigger",
	"conditions":          "Conditions represents the latest available observations of a deployment config's current state.",
	"readyReplicas":       "Total number of ready pods targeted by this deployment.",
}

func (DeploymentConfigStatus) SwaggerDoc() map[string]string {
	return map_DeploymentConfigStatus
}

var map_DeploymentDetails = map[string]string{
	"":        "DeploymentDetails captures information about the causes of a deployment.",
	"message": "Message is the user specified change message, if this deployment was triggered manually by the user",
	"causes":  "Causes are extended data associated with all the causes for creating a new deployment",
}

func (DeploymentDetails) SwaggerDoc() map[string]string {
	return map_DeploymentDetails
}

var map_DeploymentLog = map[string]string{
	"": "DeploymentLog represents the logs for a deployment\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
}

func (DeploymentLog) SwaggerDoc() map[string]string {
	return map_DeploymentLog
}

var map_DeploymentLogOptions = map[string]string{
	"":             "DeploymentLogOptions is the REST options for a deployment log\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"container":    "The container for which to stream logs. Defaults to only container if there is one container in the pod.",
	"follow":       "Follow if true indicates that the build log should be streamed until the build terminates.",
	"previous":     "Return previous deployment logs. Defaults to false.",
	"sinceSeconds": "A relative time in seconds before the current time from which to show logs. If this value precedes the time a pod was started, only logs since the pod start will be returned. If this value is in the future, no logs will be returned. Only one of sinceSeconds or sinceTime may be specified.",
	"sinceTime":    "An RFC3339 timestamp from which to show logs. If this value precedes the time a pod was started, only logs since the pod start will be returned. If this value is in the future, no logs will be returned. Only one of sinceSeconds or sinceTime may be specified.",
	"timestamps":   "If true, add an RFC3339 or RFC3339Nano timestamp at the beginning of every line of log output. Defaults to false.",
	"tailLines":    "If set, the number of lines from the end of the logs to show. If not specified, logs are shown from the creation of the container or sinceSeconds or sinceTime",
	"limitBytes":   "If set, the number of bytes to read from the server before terminating the log output. This may not display a complete final line of logging, and may return slightly more or slightly less than the specified limit.",
	"nowait":       "NoWait if true causes the call to return immediately even if the deployment is not available yet. Otherwise the server will wait until the deployment has started.",
	"version":      "Version of the deployment for which to view logs.",
}

func (DeploymentLogOptions) SwaggerDoc() map[string]string {
	return map_DeploymentLogOptions
}

var map_DeploymentRequest = map[string]string{
	"":                "DeploymentRequest is a request to a deployment config for a new deployment.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"name":            "Name of the deployment config for requesting a new deployment.",
	"latest":          "Latest will update the deployment config with the latest state from all triggers.",
	"force":           "Force will try to force a new deployment to run. If the deployment config is paused, then setting this to true will return an Invalid error.",
	"excludeTriggers": "ExcludeTriggers instructs the instantiator to avoid processing the specified triggers. This field overrides the triggers from latest and allows clients to control specific logic. This field is ignored if not specified.",
}

func (DeploymentRequest) SwaggerDoc() map[string]string {
	return map_DeploymentRequest
}

var map_DeploymentStrategy = map[string]string{
	"":                      "DeploymentStrategy describes how to perform a deployment.",
	"type":                  "Type is the name of a deployment strategy.",
	"customParams":          "CustomParams are the input to the Custom deployment strategy, and may also be specified for the Recreate and Rolling strategies to customize the execution process that runs the deployment.",
	"recreateParams":        "RecreateParams are the input to the Recreate deployment strategy.",
	"rollingParams":         "RollingParams are the input to the Rolling deployment strategy.",
	"resources":             "Resources contains resource requirements to execute the deployment and any hooks.",
	"labels":                "Labels is a set of key, value pairs added to custom deployer and lifecycle pre/post hook pods.",
	"annotations":           "Annotations is a set of key, value pairs added to custom deployer and lifecycle pre/post hook pods.",
	"activeDeadlineSeconds": "ActiveDeadlineSeconds is the duration in seconds that the deployer pods for this deployment config may be active on a node before the system actively tries to terminate them.",
}

func (DeploymentStrategy) SwaggerDoc() map[string]string {
	return map_DeploymentStrategy
}

var map_DeploymentTriggerImageChangeParams = map[string]string{
	"":                   "DeploymentTriggerImageChangeParams represents the parameters to the ImageChange trigger.",
	"automatic":          "Automatic means that the detection of a new tag value should result in an image update inside the pod template.",
	"containerNames":     "ContainerNames is used to restrict tag updates to the specified set of container names in a pod. If multiple triggers point to the same containers, the resulting behavior is undefined. Future API versions will make this a validation error. If ContainerNames does not point to a valid container, the trigger will be ignored. Future API versions will make this a validation error.",
	"from":               "From is a reference to an image stream tag to watch for changes. From.Name is the only required subfield - if From.Namespace is blank, the namespace of the current deployment trigger will be used.",
	"lastTriggeredImage": "LastTriggeredImage is the last image to be triggered.",
}

func (DeploymentTriggerImageChangeParams) SwaggerDoc() map[string]string {
	return map_DeploymentTriggerImageChangeParams
}

var map_DeploymentTriggerPolicy = map[string]string{
	"":                  "DeploymentTriggerPolicy describes a policy for a single trigger that results in a new deployment.",
	"type":              "Type of the trigger",
	"imageChangeParams": "ImageChangeParams represents the parameters for the ImageChange trigger.",
}

func (DeploymentTriggerPolicy) SwaggerDoc() map[string]string {
	return map_DeploymentTriggerPolicy
}

var map_ExecNewPodHook = map[string]string{
	"":              "ExecNewPodHook is a hook implementation which runs a command in a new pod based on the specified container which is assumed to be part of the deployment template.",
	"command":       "Command is the action command and its arguments.",
	"env":           "Env is a set of environment variables to supply to the hook pod's container.",
	"containerName": "ContainerName is the name of a container in the deployment pod template whose container image will be used for the hook pod's container.",
	"volumes":       "Volumes is a list of named volumes from the pod template which should be copied to the hook pod. Volumes names not found in pod spec are ignored. An empty list means no volumes will be copied.",
}

func (ExecNewPodHook) SwaggerDoc() map[string]string {
	return map_ExecNewPodHook
}

var map_LifecycleHook = map[string]string{
	"":              "LifecycleHook defines a specific deployment lifecycle action. Only one type of action may be specified at any time.",
	"failurePolicy": "FailurePolicy specifies what action to take if the hook fails.",
	"execNewPod":    "ExecNewPod specifies the options for a lifecycle hook backed by a pod.",
	"tagImages":     "TagImages instructs the deployer to tag the current image referenced under a container onto an image stream tag.",
}

func (LifecycleHook) SwaggerDoc() map[string]string {
	return map_LifecycleHook
}

var map_RecreateDeploymentStrategyParams = map[string]string{
	"":               "RecreateDeploymentStrategyParams are the input to the Recreate deployment strategy.",
	"timeoutSeconds": "TimeoutSeconds is the time to wait for updates before giving up. If the value is nil, a default will be used.",
	"pre":            "Pre is a lifecycle hook which is executed before the strategy manipulates the deployment. All LifecycleHookFailurePolicy values are supported.",
	"mid":            "Mid is a lifecycle hook which is executed while the deployment is scaled down to zero before the first new pod is created. All LifecycleHookFailurePolicy values are supported.",
	"post":           "Post is a lifecycle hook which is executed after the strategy has finished all deployment logic. All LifecycleHookFailurePolicy values are supported.",
}

func (RecreateDeploymentStrategyParams) SwaggerDoc() map[string]string {
	return map_RecreateDeploymentStrategyParams
}

var map_RollingDeploymentStrategyParams = map[string]string{
	"":                    "RollingDeploymentStrategyParams are the input to the Rolling deployment strategy.",
	"updatePeriodSeconds": "UpdatePeriodSeconds is the time to wait between individual pod updates. If the value is nil, a default will be used.",
	"intervalSeconds":     "IntervalSeconds is the time to wait between polling deployment status after update. If the value is nil, a default will be used.",
	"timeoutSeconds":      "TimeoutSeconds is the time to wait for updates before giving up. If the value is nil, a default will be used.",
	"maxUnavailable":      "MaxUnavailable is the maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of total pods at the start of update (ex: 10%). Absolute number is calculated from percentage by rounding down.\n\nThis cannot be 0 if MaxSurge is 0. By default, 25% is used.\n\nExample: when this is set to 30%, the old RC can be scaled down by 30% immediately when the rolling update starts. Once new pods are ready, old RC can be scaled down further, followed by scaling up the new RC, ensuring that at least 70% of original number of pods are available at all times during the update.",
	"maxSurge":            "MaxSurge is the maximum number of pods that can be scheduled above the original number of pods. Value can be an absolute number (ex: 5) or a percentage of total pods at the start of the update (ex: 10%). Absolute number is calculated from percentage by rounding up.\n\nThis cannot be 0 if MaxUnavailable is 0. By default, 25% is used.\n\nExample: when this is set to 30%, the new RC can be scaled up by 30% immediately when the rolling update starts. Once old pods have been killed, new RC can be scaled up further, ensuring that total number of pods running at any time during the update is atmost 130% of original pods.",
	"pre":                 "Pre is a lifecycle hook which is executed before the deployment process begins. All LifecycleHookFailurePolicy values are supported.",
	"post":                "Post is a lifecycle hook which is executed after the strategy has finished all deployment logic. All LifecycleHookFailurePolicy values are supported.",
}

func (RollingDeploymentStrategyParams) SwaggerDoc() map[string]string {
	return map_RollingDeploymentStrategyParams
}

var map_TagImageHook = map[string]string{
	"":              "TagImageHook is a request to tag the image in a particular container onto an ImageStreamTag.",
	"containerName": "ContainerName is the name of a container in the deployment config whose image value will be used as the source of the tag. If there is only a single container this value will be defaulted to the name of that container.",
	"to":            "To is the target ImageStreamTag to set the container's image onto.",
}

func (TagImageHook) SwaggerDoc() map[string]string {
	return map_TagImageHook
}

// AUTO-GENERATED FUNCTIONS END HERE
