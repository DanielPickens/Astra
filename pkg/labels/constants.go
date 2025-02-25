package labels

const (

	// kubernetesInstanceLabel identifies the component name
	kubernetesInstanceLabel = "app.kubernetes.io/instance"

	// kubernetesNameLabel identifies the type of the component
	// const kubernetesNameLabel = "app.kubernetes.io/name"

	// kubernetesManagedByLabel identifies the manager of the component
	kubernetesManagedByLabel = "app.kubernetes.io/managed-by"

	// kubernetesManagedByVersionLabel identifies the version of manager used to deploy the resource
	kubernetesManagedByVersionLabel = "app.kubernetes.io/managed-by-version"

	// kubernetesPartOfLabel identifies the application to which the component belongs
	kubernetesPartOfLabel = "app.kubernetes.io/part-of"

	// kubernetesStorageNameLabel is applied to all storage resources that are created
	kubernetesStorageNameLabel = "app.kubernetes.io/storage-name"

	openshiftRunTimeLabel = "app.openshift.io/runtime"

	// astraModeLabel indicates which command were used to create the component, either dev or deploy
	astraModeLabel = "astra.dev/mode"

	// astraProjectTypeAnnotation indicates the project type of the component
	astraProjectTypeAnnotation = "astra.dev/project-type"

	appLabel = "app"

	componentLabel = "component"

	// devfileStorageLabel is applied to all storage resources for devfile components that are created
	devfileStorageLabel = "storage-name"

	sourcePVCLabel = "astra-source-pvc"
)

const (
	// ComponentDevMode indicates the resource is deployed using dev command
	ComponentDevMode = "Dev"

	// ComponentDeployMode indicates the resource is deployed using deploy command
	ComponentDeployMode = "Deploy"

	// ComponentAnyMode is used to search resources deployed using either dev or deploy command
	ComponentAnyMode = ""

	// astraManager is the value of the manager when a component is managed by astra
	astraManager = "astra"
)
