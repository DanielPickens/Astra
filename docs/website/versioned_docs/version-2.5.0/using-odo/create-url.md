---
title: Create URLs using astra
sidebar_position: 2
sidebar_label: Create URL
---

In the [previous section](./create-component) we created two components — a Spring Boot application (`backend`) listening on port 8080 and a Nodejs application (`frontend`) listening on port 3000 — and pushed them to the Kubernetes cluster. These are also the respective default ports (8080 for Spring Boot and 3000 for Nodejs) for Spring Boot and Nodejs component types. In this guide, we will create URLs to access these components from the host system.

Note that the URLs we create in this section will only help you access the components in web browser; the application itself won't be usable till we create some services and links which we will cover in the next section.

## OpenShift

If you are using [Code Ready Containers (CRC)](https://github.com/code-ready/crc) or another form of OpenShift cluster, astra has already created URLs for you by using the [OpenShift Routes](https://docs.openshift.com/container-platform/latest/networking/routes/route-configuration.html) feature. Execute `astra url list` from the component directory of the `backend` and `frontend` components to get the URLs astra created for these components. If you observe the `astra push` output closely, astra prints the URL in it as well.

Below are example `astra url list` outputs for the backend and frontend components. Note that URLs would be different in your case:

```shell
# backend component
$ astra url list
Found the following URLs for component backend
NAME         STATE      URL                                            PORT     SECURE     KIND
8080-tcp     Pushed     http://8080-tcp-app-myproject.hostname.com     8080     false      route

# frontend component
$ astra url list
Found the following URLs for component frontend
NAME          STATE      URL                                             PORT     SECURE     KIND
http-3000     Pushed     http://http-3000-app-myproject.hostname.com     3000     false      route

```

## Kubernetes

If you are using a Kubernetes cluster, you will have to create a URL using `astra url` command. This is because astra can not assume the host information to be used to create a URL. To be able to create URLs on a Kubernetes cluster, please make sure that you have [Ingress Controller](../getting-started/cluster-setup/kubernetes/#enabling-ingress) installed.

If you are working on a [minikube](../getting-started/cluster-setup/kubernetes), Ingress can be enabled using:
```shell
minikube addons enable ingress
```

If you are working on any other kind of Kubernetes cluster, please check with your cluster administrator to enable the Ingress Controller. In this guide, we cover URL creation for minikube setup. For any other Kubernetes cluster, please replace `$(minikube ip).nip.io` in below commands with the host information for your specific cluster.

### Backend component

Our backend component, which is based on Spring Boot, listens on port 8080. `cd` into the directory for this component and execute below command:

```shell
astra url create --port 8080 --host $(minikube ip).nip.io
astra push
```
astra follows a "create & push" workflow for most commands. But in this case, adding `--now` flag to `astra url create` could reduce two commands into a single command:
```shell
astra url create --port 8080 --host $(minikube ip).nip.io --now
````
### Frontend component

Our frontend component, which is based on Nodejs, listens on port 3000. `cd` into the directory for this component and execute below command:

```shell
astra url create --port 3000 --host $(minikube ip).nip.io
astra push
```
Again, if you would prefer to get this done in a single command:
```shell
astra url create --port 3000 --host $(minikube ip).nip.io --now
```
