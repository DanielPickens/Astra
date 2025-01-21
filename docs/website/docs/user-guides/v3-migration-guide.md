---
title: Migrate from v2 to v3
sidebar_position: 4
---

### Migrate an existing `astra` component from v2 to v3
If you have created an `astra` component using `astra` v2, this section will help you move it to use `astra` v3.
#### Step 1 
`cd` into the component directory, and delete the component from the Kubernetes cluster:
```shell
astra delete
```
#### Step 2
Download and [install astra v3](../overview/installation.md).

#### Step 3
Run [`astra dev`](../command-reference/dev.md) to start developing your application using `astra` v3.

#### Step 4
Run `astra list` to see a list of components that are running on the cluster and what version of `astra` they are running.

### Where are `astra push` and `astra watch` commands?
`astra push` and `astra watch` have been replaced in v3 by a single command - `astra dev`. 

In v2, if you wanted to automatically sync the code on developer system with the application running on a Kubernetes
cluster, you had to perform two steps - first `astra push` to start the application on the cluster, and second `astra watch`
to automatically sync the code. In v3, `astra dev` performs both actions with a single command.

`astra dev` is not _just_ a replacement for these two commands. It's also different in behaviour in that, it's a 
long-running process that's going to block the terminal. Hitting `Ctrl+c` will stop the process and cleanup the 
component from the cluster. In v2, you had to use `astra delete`/`astra component delete` to delete inner loop resources 
of the component from the cluster.

### What happened to Ingress/Route?
If you have used `astra` v2, you must have used Ingress (on Kubernetes) or Route (on OpenShift) to access the 
application that was pushed to the cluster using `astra push`. `astra` v3 no longer creates an Ingress or a Route. Instead,
it uses port-forwarding.

When running `astra dev`, `astra` forwards a port on the development system to the port on the container cluster allowing 
you remote access to your deployed application. It also prints the information when the application has started on the
cluster:
```shell
$ astra dev
...
...
-  Forwarding from 127.0.0.1:40001 -> 8080
```
This indicates that the port 40001 on the development system has been forwarded to port 8080 of the application 
represented by the current `astra` component.

:::info NOTE
`astra` no longer supports creation of Ingress / Route out of the box. The `astra url` set of commands no longer exist 
in v3.
:::

### Changes to the way component debugging works
In `astra` v2, `astra push --debug` was used to run a component in debug mode. To setup port forwarding to the component's
debug port, you had to run `astra debug port-forward`.

In `astra` v3, you need to specify the debug port in the `devfile.yaml` as an endpoint, and run `astra dev --debug` to 
start the component in debug mode. For example, a `container` component in the devfile should look like below, where 
port 3000 is the application port and 5858 is the debug port:

```yaml
- name: runtime
  container:
    image: registry.access.redhat.com/ubi8/nodejs-12:1-36
    memoryLimit: 1024Mi
    endpoints:
    - name: "3000-tcp"
      targetPort: 3000
    - name: debug
      exposure: none
      targetPort: 5858
```
### Changes to default configurations

#### Ephemeral storage

By default, `astra` v2 used [ephemeral storage](https://docs.openshift.com/container-platform/4.11/storage/understanding-ephemeral-storage.html) 
for the components created using it. However, this has changed in `astra` v3, and it now uses the underlying storage 
(Persistent Volumes) configured for use by the users. If you would like to continue using ephemeral storage for `astra` 
components, you could change the configuration by doing:
```shell
astra preference set Ephemeral true
```

### Commands added, modified or removed in v3

The following table contains a list of `astra` commands that have either been modified or removed. In case of a 
modification, the modified command is mentioned in the `v3` column. Please refer below legend beforehand:
* 👷 currently not implemented, but might get implemented in future
* ❌ not implemented, no plans for implementation

| v2                                  | v3                                                                                                                                                                                                                                            |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| app delete                          | ❌                                                                                                                                                                                                                                             |
| app describe                        | ❌                                                                                                                                                                                                                                             |
| app list                            | ❌                                                                                                                                                                                                                                             |
| catalog describe service            | ❌                                                                                                                                                                                                                                             |
| catalog describe component          | registry --details                                                                                                                                                                                                                            |
| catalog list service                | 👷astra list services                                                                                                                                                                                                                           |
| catalog list component              | registry                                                                                                                                                                                                                                      |
| catalog search service              | ❌                                                                                                                                                                                                                                             |
| catalog search component            | registry --filter                                                                                                                                                                                                                             |
| config set                          | ❌                                                                                                                                                                                                                                             |
| config unset                        | ❌                                                                                                                                                                                                                                             |
| config view                         | ❌                                                                                                                                                                                                                                             |  |
| debug info                          | ❌ (not needed as debug mode is start with `astra dev --debug` command that blocks terminal, the application can be running in inner-loop mode only if `astra dev` is running in terminal)                                                             |
| debug port-forward                  | ❌ (port forwarding is automatic when users execute `astra dev --debug` as long as [the endpoint is defined in the devfile](#changes-to-the-way-component-debugging-works))                                                                      |
| env set                             | ❌                                                                                                                                                                                                                                             |
| env uset                            | ❌                                                                                                                                                                                                                                             |
| env view                            | ❌                                                                                                                                                                                                                                             |
| preference set                      | preference set                                                                                                                                                                                                                                |
| preference unset                    | preference unset                                                                                                                                                                                                                              |
| preference view                     | preference view                                                                                                                                                                                                                               |
| project create                      | create namespace                                                                                                                                                                                                                              |
| project delete                      | delete namespace                                                                                                                                                                                                                              |
| project get                         | ❌                                                                                                                                                                                                                                             |
| project list                        | list namespace                                                                                                                                                                                                                                |
| project set                         | set namespace                                                                                                                                                                                                                                 |
| registry add                        | preference add registry                                                                                                                                                                                                                       |
| registry delete                     | preference remove registry                                                                                                                                                                                                                    |
| registry list                       | preference view                                                                                                                                                                                                                               |
| registry update                     | no command for update. If needed, it can be done using preference remove registry and preference add registry                                                                                                                                 |
| service create/delete/describe/list | ❌                                                                                                                                                                                                                                             |
| storage create/delete/list          | ❌                                                                                                                                                                                                                                             |
| test                                | 👷(will be implemented after v3-GA) #6070                                                                                                                                                                                                     |
| url create/delete/list              | ❌ (`astra dev` automatically sets port forwarding between container and localhost. If users for some reason require Ingress or Route for inner-loop development they will have to explicitly define them in the devfile as kubernetes components) |
| build-images                        | build-images                                                                                                                                                                                                                                  |  |
| deploy                              | deploy                                                                                                                                                                                                                                        |
| login                               | login                                                                                                                                                                                                                                         |
| logout                              | logout                                                                                                                                                                                                                                        |
| create / component create           | init                                                                                                                                                                                                                                          |
| delete / component delete           | delete component                                                                                                                                                                                                                              |
| describe / component describe       | describe component                                                                                                                                                                                                                            |
| exec / component exec               | ❌                                                                                                                                                                                                                                             |
| link / component link               | add binding                                                                                                                                                                                                                                   | list / component list |
| list                                | log / component log                                                                                                                                                                                                                           |
| logs                                | push / component push                                                                                                                                                                                                                         | in v3 there is only dev command that behaves more like watch in v2. There is an option to disable automatic reloads when a file is changed (--no-watch) in future there will be an option to trigger “sync” explicitly when `--no-watch` option is used) |
| status / component status           | ❌                                                                                                                                                                                                                                             |
| unlink / component unlink           | remove binding                                                                                                                                                                                                                                |
| watch / component watch             | dev                                                                                                                                                                                                                                           |
| ❌                                   | describe binding                                                                                                                                                                                                                              |
| ❌                                   | list binding                                                                                                                                                                                                                                  |
| ❌                                   | analyze                                                                                                                                                                                                                                       |