---
title: astra link
sidebar_position: 7
---

`astra link` command helps link an astra component to an Operator backed service or another astra component. It does this by using [Service Binding Operator](https://github.com/daniel-pickens/service-binding-operator). At the time of writing this, astra makes use of the Service Binding library and not the Operator itself to achieve the desired functionality.

In this document we will cover various options to create link between a component & a service, and a component & another component. The steps in this document are going to be based on the [astra quickstart project](https://github.com/dharmit/astra-quickstart/) that we covered in [Quickstart guide](../getting-started/quickstart). The outputs mentioned in this document are based on commands executed on [minikube cluster](../getting-started/cluster-setup/kubernetes).

This document assumes that you know how to [create components](../command-reference/create) and [services](../command-reference/service). It also assumes that you have cloned the [astra quickstart project](https://github.com/dharmit/astra-quickstart/). Terminology used in this document:

- *quickstart project*: git clone of the astra quickstart project having below directory structure:
    ```shell
    $ tree -L 1
    .
    ├── backend
    ├── frontend
    ├── postgrescluster.yaml
    ├── quickstart.code-workspace
    └── README.md
    
    2 directories, 3 files
    ```
- *backend component*: `backend` directory in above tree structure
- *frontend component*: `frontend` directory in above tree structure
- *Postgres service*: Operator backed service created from *backend component* using the `astra service create --from-file ../postgrescluster.yaml` command.

## Various linking options

astra provides various options to link a component with an Operator backed service or another astra component. All these options (or flags) can be used irrespective of whether you are linking a component to a service or another component.

### Default behaviour

By default, `astra link` creates a directory named `kubernetes/` in your component directory and stores the information (YAML manifests) about services and links in it. When you do `astra push`, astra compares these manifests with the state of the things on the Kubernetes cluster and decides whether it needs to create, modify or destroy resources to match what is specified by the user.

### The `--inlined` flag

If you specified `--inlined` flag to the `astra link` command, astra will store the link information inline in the `devfile.yaml` in the component directory instead of creating a file under `kubernetes/` directory. The behaviour of `--inlined` flag is similar in both the `astra link` and `astra service create` commands. This flag is helpful if you would like everything to be stored in a single `devfile.yaml`. You will have to remember to use `--inlined` flag with each `astra link` and `astra service create` commands that you execute for the component.

### The `--map` flag

At times, you might want to add more binding information to the component than what is available by default. For example, if you are linking the component with a service and would like to bind some information from the service's spec (short for specification), you could use the `--map` flag. Note that astra doesn't do any validation against the spec of the service/component being linked. Using this flag is recommended only if you are comfortable with reading the Kubernetes YAML manifests.

### The `--bind-as-files` flag

For all the linking options discussed so far, astra injects the binding information into the component as environment variables. If you would like to instead mount this information as files, you could use the `--bind-as-files` flag. This will make astra inject the binding information as files into the `/bindings` location within your component's Pod. Comparing with the environment variables paradigm, when you use `--bind-as-files`, the files are named after the keys and the value of these keys is stored as the contents of these files.

## Examples

### Default `astra link`

We will link the backend component with the Postgres service using default `astra link` command. For the backend component, make sure that your component and service are pushed to the cluster:

```shell
$ astra list
APP     NAME        PROJECT       TYPE       STATE      MANAGED BY astra
app     backend     myproject     spring     Pushed     Yes


$ astra service list
NAME                      MANAGED BY astra     STATE      AGE
PostgresCluster/hippo     Yes (backend)      Pushed     59m41s

```

Now, run `astra link` to link the backend component with the Postgres service:
```shell
astra link PostgresCluster/hippo
```
Example output:
```shell
$ astra link PostgresCluster/hippo
 ✓  Successfully created link between component "backend" and service "PostgresCluster/hippo"

To apply the link, please use `astra push`
```
And then run `astra push` for the link to actually get created on the Kubernetes cluster.

Upon successful `astra push`, you can notice a few things:
1. When you open the URL for the application deployed by backend component, it shows you a list of tastra items in the database. For example, for below `astra url list` output, we will append the path where tastras are listed:
  ```shell
  $ astra url list
  Found the following URLs for component backend
  NAME         STATE      URL                                       PORT     SECURE     KIND
  8080-tcp     Pushed     http://8080-tcp.192.168.39.112.nip.io     8080     false      ingress
  
  ```
  The correct path for such URL would be - http://8080-tcp.192.168.39.112.nip.io/api/v1/tastras. Note that exact URL would be different for your setup. Also note that there are no tastras in the database unless you add some, so the URL might just show an empty JSON object.
2. You can see binding information related to Postgres service injected into the backend component. This binding information is injected, by default, as environment variables. You can check it out using the `astra describe` command from backend component's directory:
  ```shell
  astra describe 
  ```
  Example output:
    ```shell
    $ astra describe
    Component Name: backend
    Type: spring
    Environment Variables:
     · PROJECTS_ROOT=/projects
     · PROJECT_SOURCE=/projects
     · DEBUG_PORT=5858
    Storage:
     · m2 of size 3Gi mounted to /home/user/.m2
    URLs:
     · http://8080-tcp.192.168.39.112.nip.io exposed via 8080
    Linked Services:
     · PostgresCluster/hippo
       Environment Variables:
        · POSTGRESCLUSTER_PGBOUNCER-EMPTY
        · POSTGRESCLUSTER_PGBOUNCER.INI
        · POSTGRESCLUSTER_ROOT.CRT
        · POSTGRESCLUSTER_VERIFIER
        · POSTGRESCLUSTER_ID_ECDSA
        · POSTGRESCLUSTER_PGBOUNCER-VERIFIER
        · POSTGRESCLUSTER_TLS.CRT
        · POSTGRESCLUSTER_PGBOUNCER-URI
        · POSTGRESCLUSTER_PATRONI.CRT-COMBINED
        · POSTGRESCLUSTER_USER
        · pgImage
        · pgVersion
        · POSTGRESCLUSTER_CLUSTERIP
        · POSTGRESCLUSTER_HOST
        · POSTGRESCLUSTER_PGBACKREST_REPO.CONF
        · POSTGRESCLUSTER_PGBOUNCER-USERS.TXT
        · POSTGRESCLUSTER_SSH_CONFIG
        · POSTGRESCLUSTER_TLS.KEY
        · POSTGRESCLUSTER_CONFIG-HASH
        · POSTGRESCLUSTER_PASSWORD
        · POSTGRESCLUSTER_PATRONI.CA-ROOTS
        · POSTGRESCLUSTER_DBNAME
        · POSTGRESCLUSTER_PGBOUNCER-PASSWORD
        · POSTGRESCLUSTER_SSHD_CONFIG
        · POSTGRESCLUSTER_PGBOUNCER-FRONTEND.KEY
        · POSTGRESCLUSTER_PGBACKREST_INSTANCE.CONF
        · POSTGRESCLUSTER_PGBOUNCER-FRONTEND.CA-ROOTS
        · POSTGRESCLUSTER_PGBOUNCER-HOST
        · POSTGRESCLUSTER_PORT
        · POSTGRESCLUSTER_ROOT.KEY
        · POSTGRESCLUSTER_SSH_KNOWN_HOSTS
        · POSTGRESCLUSTER_URI
        · POSTGRESCLUSTER_PATRONI.YAML
        · POSTGRESCLUSTER_DNS.CRT
        · POSTGRESCLUSTER_DNS.KEY
        · POSTGRESCLUSTER_ID_ECDSA.PUB
        · POSTGRESCLUSTER_PGBOUNCER-FRONTEND.CRT
        · POSTGRESCLUSTER_PGBOUNCER-PORT
        · POSTGRESCLUSTER_CA.CRT
    ```
  Few of these variables are used in the backend component's [`src/main/resources/application.properties` file](https://github.com/dharmit/astra-quickstart/blob/main/backend/src/main/resources/application.properties) so that the Java Springboot application can connect to the Postgres database service.
3. Lastly, astra has created a directory called `kubernetes/` in your backend component's directory which contains below files. 
  ```shell
  $ ls kubernetes 
  astra-service-backend-postgrescluster-hippo.yaml  astra-service-hippo.yaml
  ```
  This files contains the information (YAML manifests) about two things:
    1. `astra-service-hippo.yaml` - the Postgres service we created using `astra service create --from-file ../postgrescluster.yaml` command.
    2. `astra-service-backend-postgrescluster-hippo.yaml` - the link we created using `astra link` command.
  
### `astra link` with `--inlined`

Using `--inlined` flag with `astra link` command does the exact same thing to our application (that is, injects binding information) as an `astra link` command without the flag does. However, the subtle difference is that in above case we saw two manifest files under `kubernetes/` directory — one for the Postgres service and other for the link between the backend component and this service — but when we pass `--inlined` flag, astra does not create a file under `kubernetes/` directory to store the YAML manifest, but stores it inline in the `devfile.yaml` file.

To see this, let's unlink our component from the Postgres service first:

```shell
astra unlink PostgresCluster/hippo
```
Example output:
```shell
$ astra unlink PostgresCluster/hippo
 ✓  Successfully unlinked component "backend" from service "PostgresCluster/hippo"

To apply the changes, please use `astra push`
```
To unlink them on the cluster, run `astra push`. Now if you take a look at the `kubernetes/` directory, you'll see only one file in it:
```shell
$ ls kubernetes 
astra-service-hippo.yaml
```
Next, let's use the `--inlined` flag to create a link:
```shell
astra link PostgresCluster/hippo --inlined
```
Example output:
```shell
$ astra link PostgresCluster/hippo --inlined
 ✓  Successfully created link between component "backend" and service "PostgresCluster/hippo"

To apply the link, please use `astra push`
```
Just like the time without `--inlined` flag, you need to do `astra push` for the link to get created on the cluster. But where did astra store the configuration/manifest required to create this link? astra stores this in `devfile.yaml`. You can see an entry like below in this file:
```yaml
 kubernetes:
    inlined: |
      apiVersion: binding.operators.coreos.com/v1alpha1
      kind: ServiceBinding
      metadata:
        creationTimestamp: null
        name: backend-postgrescluster-hippo
      spec:
        application:
          group: apps
          name: backend-app
          resource: deployments
          version: v1
        bindAsFiles: false
        detectBindingResources: true
        services:
        - group: postgres-operator.crunchydata.com
          id: hippo
          kind: PostgresCluster
          name: hippo
          version: v1beta1
      status:
        secret: ""
  name: backend-postgrescluster-hippo
```
Now if you were to do `astra unlink PostgresCluster/hippo`, astra would first remove the link information from the `devfile.yaml` and then a subsequent `astra push` would delete the link from the cluster.

### Custom bindings

`astra link` accepts the flag `--map` which can inject custom binding information into the component. Such binding information will be fetched from the manifest of the resource we are linking to our component. For example, speaking in context of the backend component and Postgres service, we can inject information from the Postgres service's manifest ([`postgrescluster.yaml` file](https://github.com/dharmit/astra-quickstart/blob/main/postgrescluster.yaml)) into the backend component.

Considering the name of your `PostgresCluster` service is `hippo` (check the output of `astra service list` if your PostgresCluster service is named differently), if we wanted to inject the value of `postgresVersion` from that YAML definition into our backend component:
```shell
astra link PostgresCluster/hippo --map pgVersion='{{ .hippo.spec.postgresVersion }}'
```
Note that, if the name of your Postgres service is different from `hippo`, you will have to specify that in the above command in place `.hippo`. For example, if your `PostgresCluster` service is named as `database`, you would change the link command to as shown below:

```shell
$ astra service list
NAME                      MANAGED BY astra     STATE      AGE
PostgresCluster/database     Yes (backend)      Pushed     2h5m43s

$ astra link PostgresCluster/hippo --map pgVersion='{{ .database.spec.postgresVersion }}'
```

After a link operation, do `astra push` as usual. Upon successful completion of push operation, you can run below command from your backend component directory to validate if custom mapping got injected properly:

```shell
astra exec -- env | grep pgVersion
```
Example output:
```shell
$ astra exec -- env | grep pgVersion
pgVersion=13
```

Since a user might want to inject more than just one piece of custom binding information, `astra link` accepts multiple key-value pairs of mappings. The only constraint being that these should be specified as `--map <key>=<value>`. For example, if you want to also inject Postgres image information along with the version, you could do:

```shell
astra link PostgresCluster/hippo --map pgVersion='{{ .hippo.spec.postgresVersion }}' --map pgImage='{{ .hippo.spec.image }}'
```
and do `astra push`. The way to validate if both the mappings got injected correctly would be to do:
```shell
astra exec -- env | grep -e "pgVersion\|pgImage"
```
Example output:
```shell
$ astra exec -- env | grep -e "pgVersion\|pgImage"
pgVersion=13
pgImage=registry.developers.crunchydata.com/crunchydata/crunchy-postgres-ha:centos8-13.4-0

```

#### To inline or not?

You can stick to the default behaviour wherein `astra link` will generate a manifest file for the link under `kubernetes/` directory, or you could use `--inlined` flag if you prefer to store everything in a single `devfile.yaml` file. It doesn't matter what you use for this functionality of adding custom mappings.

## Binding as files

Another helpful flag that `astra link` provides is called `--bind-as-files`. When this flag is passed, the binding information is not injected into the component's Pod as environment variables but is mounted as a filesystem. We will see a few examples that will make things clearer.

Ensure that there are no existing links between the backend component and the Postgres service. You could do this by running `astra describe` in the backend component's directory and check if you see something like below in the output:
```shell
Linked Services:
 · PostgresCluster/hippo
```
Unlink the service from the component using:
```shell
astra unlink PostgresCluster/hippo
astra push
```

## `--bind-as-files` examples

### With default `astra link`

Default behaviour means astra creating the manifest file under `kubernetes/` directory to store the link information. Link the backend component and Postgres service using:
```shell
astra link PostgresCluster/hippo --bind-as-files
astra push
```

Example `astra describe` output:
```shell
$ astra describe
Component Name: backend
Type: spring
Environment Variables:
 · PROJECTS_ROOT=/projects
 · PROJECT_SOURCE=/projects
 · DEBUG_PORT=5858
 · SERVICE_BINDING_ROOT=/bindings
 · SERVICE_BINDING_ROOT=/bindings
Storage:
 · m2 of size 3Gi mounted to /home/user/.m2
URLs:
 · http://8080-tcp.192.168.39.112.nip.io exposed via 8080
Linked Services:
 · PostgresCluster/hippo
   Files:
    · /bindings/backend-postgrescluster-hippo/pgbackrest_instance.conf
    · /bindings/backend-postgrescluster-hippo/user
    · /bindings/backend-postgrescluster-hippo/ssh_known_hosts
    · /bindings/backend-postgrescluster-hippo/clusterIP
    · /bindings/backend-postgrescluster-hippo/password
    · /bindings/backend-postgrescluster-hippo/patroni.yaml
    · /bindings/backend-postgrescluster-hippo/pgbouncer-frontend.crt
    · /bindings/backend-postgrescluster-hippo/pgbouncer-host
    · /bindings/backend-postgrescluster-hippo/root.key
    · /bindings/backend-postgrescluster-hippo/pgbouncer-frontend.key
    · /bindings/backend-postgrescluster-hippo/pgbouncer.ini
    · /bindings/backend-postgrescluster-hippo/uri
    · /bindings/backend-postgrescluster-hippo/config-hash
    · /bindings/backend-postgrescluster-hippo/pgbouncer-empty
    · /bindings/backend-postgrescluster-hippo/port
    · /bindings/backend-postgrescluster-hippo/dns.crt
    · /bindings/backend-postgrescluster-hippo/pgbouncer-uri
    · /bindings/backend-postgrescluster-hippo/root.crt
    · /bindings/backend-postgrescluster-hippo/ssh_config
    · /bindings/backend-postgrescluster-hippo/dns.key
    · /bindings/backend-postgrescluster-hippo/host
    · /bindings/backend-postgrescluster-hippo/patroni.crt-combined
    · /bindings/backend-postgrescluster-hippo/pgbouncer-frontend.ca-roots
    · /bindings/backend-postgrescluster-hippo/tls.key
    · /bindings/backend-postgrescluster-hippo/verifier
    · /bindings/backend-postgrescluster-hippo/ca.crt
    · /bindings/backend-postgrescluster-hippo/dbname
    · /bindings/backend-postgrescluster-hippo/patroni.ca-roots
    · /bindings/backend-postgrescluster-hippo/pgbackrest_repo.conf
    · /bindings/backend-postgrescluster-hippo/pgbouncer-port
    · /bindings/backend-postgrescluster-hippo/pgbouncer-verifier
    · /bindings/backend-postgrescluster-hippo/id_ecdsa
    · /bindings/backend-postgrescluster-hippo/id_ecdsa.pub
    · /bindings/backend-postgrescluster-hippo/pgbouncer-password
    · /bindings/backend-postgrescluster-hippo/pgbouncer-users.txt
    · /bindings/backend-postgrescluster-hippo/sshd_config
    · /bindings/backend-postgrescluster-hippo/tls.crt
```
Everything that was an environment variable in the `key=value` format in the earlier `astra describe` output is now mounted as file. Let's we `cat` the contents of few of these files:
```shell
$ astra exec -- cat /bindings/backend-postgrescluster-hippo/password
q({JC:jn^mm/Bw}eu+j.GX{k

$ astra exec -- cat /bindings/backend-postgrescluster-hippo/user    
hippo

$ astra exec -- cat /bindings/backend-postgrescluster-hippo/clusterIP
10.101.78.56
```

### With `--inlined`

The result of using `--bind-as-files` and `--inlined` together is similar to `astra link --inlined`, in that, the manifest of the link gets stored in the `devfile.yaml` instead of being stored in a separate file under `kubernetes/` directory. Other than that, the `astra describe` output would like same as saw in the [above section](#with-default-astra-link).

### Custom bindings

When you pass custom bindings while linking the backend component with the Postgres service, these custom bindings are injected not as environment variables but mounted as files. Consider below example:

```shell
astra link PostgresCluster/hippo --map pgVersion='{{ .hippo.spec.postgresVersion }}' --map pgImage='{{ .hippo.spec.image }}' --bind-as-files
astra push
```

These custom bindings got mounted as files instead of being injected as environment variables. The way to validate if that worked would be:
```shell
$ astra exec -- cat /bindings/backend-postgrescluster-hippo/pgVersion
13

$ astra exec -- cat /bindings/backend-postgrescluster-hippo/pgImage  
registry.developers.crunchydata.com/crunchydata/crunchy-postgres-ha:centos8-13.4-0
```
