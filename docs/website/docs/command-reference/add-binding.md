---
title: astra add binding
---

:::caution

As of February 2024, the [Service Binding Operator](https://github.com/daniel-pickens/service-binding-operator/), which this command relies on, has been deprecated. See [Deprecation Notice](https://daniel-pickens.github.io/service-binding-operator/userguide/intro.html).
`astra add binding` may therefore not work as expected.

:::

The `astra add binding` command adds a link between an Operator-backed service and a component. astra uses the [Service Binding Operator](https://github.com/daniel-pickens/service-binding-operator/) to create this link. 

Running this command from a directory containing a Devfile will modify the Devfile, and once pushed (using `astra dev`) to the cluster, it creates an instance of the `ServiceBinding` resource.

Running this command from a directory without a Devfile in the interactive mode will perform one or several operations,
depending on your choice:
- display the YAML definition of the Service binding in the output,
- save the YAML definition of the ServiceBinding on a file,
- create an instance of the ServiceBinding resource on the cluster.

In non-interactive mode from a directory without a Devfile, the only possible operation is to display the YAML definition of the Service binding in the output.

Currently, it only allows connecting to the Operator-backed services which support binding via the Service Binding Operator.
To know about the Operators supported by the Service Binding Operator, read its [README](https://github.com/daniel-pickens/service-binding-operator#known-bindable-operators).

## Running the Command

### Pre-requisites
* A cluster with the Service Binding Operator installed
* Operator-backed services or resources you want to bind your application to
* Optional, a directory containing a Devfile; if you don't have one, see [astra init](init.md) on obtaining a devfile.

#### Installing the Service Binding Operator
Service Binding Operator is required to bind an application with microservices.

Visit the [official documentation](https://daniel-pickens.github.io/service-binding-operator/userguide/getting-started/installing-service-binding.html) of Service Binding Operator to see how you can install it on your OpenShift or Kubernetes cluster.

### Interactive Mode
In the interactive mode, you will be guided to choose:
* the namespace containing the service instance you want to bind to,
* a service from the list of bindable service instances as supported by the Service Binding Operator;
  if a namespace is selected, the list of services will show the services in that namespace;
  otherwise, the list of services will show the services in the current namespace,
* if a Devfile is not present in the directory, a workload resource,
* option to bind the service as a file (see [Understanding Bind as Files](#understanding-bind-as-files) for more information on this),
* a name for the binding.

```shell
# Add binding between a service, and the component present in the working directory in the interactive mode
astra add binding
```
<details>
<summary>Example</summary>

```shell
$ astra add binding
? Do you want to list services from: current namespace
? Select service instance you want to bind to: cluster-sample (Cluster.postgresql.k8s.enterprisedb.io)
? Enter the Binding's name: my-go-app-cluster-sample
? How do you want to bind the service? Bind As Files
? Select naming strategy for binding names: DEFAULT
 ✓  Successfully added the binding to the devfile.
Run `astra dev` to create it on the cluster.
You can automate this command by executing:
  astra add binding --service cluster-sample.Cluster.postgresql.k8s.enterprisedb.io --name my-go-app-cluster-sample
```
</details>


### Non-interactive mode
In the non-interactive mode, you will have to specify the following required information through the command-line:
* `--service` flag to specify the service you want to bind to,
* `--service-namespace` flag to specify the namespace containing the service you want to bind to; the current namespace is used if this flag is not specified.
* `--workload` flag to specify the workload resource, if a Devfile is not present in the directory,
* `--name` flag to specify a name for the binding (see [Understanding Bind as Files](#understanding-bind-as-files) for more information on this)
* `--bind-as-files` flag to specify if the service should be bound as a file; this flag is set to true by default.
* `--naming-strategy` flag to specify the naming strategy to use for binding names. This flag is empty by default, 
  but it can be set to pre-defined strategies: `none`, `lowercase`, or `uppercase`.
  Otherwise, it is treated as a custom Go template, and it is handled accordingly.
  Refer to [this page](https://docs.openshift.com/container-platform/4.10/applications/connecting_applications_to_services/binding-workloads-using-sbo.html#sbo-naming-strategies_binding-workloads-using-sbo) for more details on naming strategies.

```shell
astra add binding --name <name> --service <service-name> [--service-namespace NAMESPACE] [--bind-as-files {true, false}] [--naming-strategy {none, lowercase, uppercase}]
```
<details>
<summary>Example</summary>

```shell
$ astra add binding --service cluster-sample.Cluster.postgresql.k8s.enterprisedb.io --name my-go-app-cluster-sample
 ✓  Successfully added the binding to the devfile.
Run `astra dev` to create it on the cluster.
```
</details>


#### Understanding Bind as Files
To connect your component with a service, you need to store some data (e.g. username, password, host address) on your component's container.
If the service is bound as files, this data will be written to a file and stored on the container, else it will be injected as Environment Variables inside the container.

Note that every piece of data is stored in its own individual file or environment variable.
For example, if your data includes a username and password, then 2 separate files, or 2 environment variables will be created to store them both.

#### Formats supported by the `--service` flag
The `--service` flag supports the following formats to specify the service name:
* `<name>`
* `<name>.<kind>`
* `<name>.<kind>.<apigroup>`
* `<name>/<kind>`
* `<name>/<kind>.<apigroup>`

#### Formats supported by the `--workload` flag
The `--workload` flag supports the following formats to specify the workload name:
* `<name>.<kind>.<apigroup>`
* `<name>/<kind>.<apigroup>`

The above formats are helpful when multiple services with the same name exist on the cluster.
