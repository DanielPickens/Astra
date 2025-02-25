---
title: Creating Kubernetes resources
sidebar_position: 2
---
# Creating Kubernetes resources using astra

While astra is mainly focused on application developers who would like to care less about Kubernetes and more about getting their application running on top of it, it also tries to make things simple for application architects or devfile stack authors who are comfortable with Kubernetes. One such feature of astra that we will discuss in this guide is creation of Kubernetes resources like Pods, Deployments, Services (Kubernetes Services, not the Operator backed ones) and such using astra. Using this, if an advanced user would like to create some Kubernetes resources, they could edit the `devfile.yaml` and add it there. An `astra push` after the edit would create the resource on the cluster. A resource thus created would co-exist with an astra component. 

In this guide, we will create an nginx Deployment using its Kubernetes manifest. We will write this manifest in the `devfile.yaml`. Upon doing `astra push`, you will be able to see a Deployment and its Pods along with the component on the Kubernetes cluster using `kubectl` commands.

## Create an astra component

As with other resources like URL, Storage and Services, to create a Kubernetes resource, we first need to have an astra component. We will keep it simple by using a nodejs starter project. In an empty directory, create a component using:
```shell
astra create nodejs --starter nodejs-starter
```
Example:
```shell
$ astra create nodejs --starter nodejs-starter
Devfile Object Validation
 ✓  Checking devfile existence [52416ns]
 ✓  Creating a devfile component from registry: stage [95517ns]
Validation
 ✓  Validating if devfile name is correct [97488ns]

Starter Project
 ✓  Downloading starter project nodejs-starter from https://github.com/astra-devfiles/nodejs-ex.git [593ms]

Please use `astra push` command to create the component with source deployed
```

## Reference the Deployment manifest in `devfile.yaml`

astra supports referencing a URI in the `devfile.yaml` so that you don't have to copy the entire manifest into it. Assuming you have a Deployment manifest like below stored in a file called `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
    replicas: 3
    selector:
      matchLabels:
        app: nginx
    template:
      metadata:
        labels:
          app: nginx
      spec:
        containers:
          - name: nginx
            image: quay.io/bitnami/nginx
            ports:
              - containerPort: 80
```

Add below to the `components` section in `devfile.yaml`. Note that there is already a `runtime` component in this section:

```yaml
- name: nginx-deploy
  kubernetes:
    uri: deployment.yaml
```
## Push to the cluster

Now you need to do `astra push`. astra will create the component and also the nginx deployment for you on the cluster. Note that, unlike for Operator backed service, astra won't show you any message indicating that a service or some resource was created on the cluster. This is by design because the feature is meant for advanced users who can play with resources created in this way through `kubectl` CLI. However, if astra fails to create a resource, it will error out and let you know about it.

See if the Deployment and its Pods were created on the cluster using:
```shell
kubectl get deploy
kubectl get pods
```

## Good to know

astra adds a Kubernetes label to the resources created in this way. It is `app.kubernetes.io/managed-by: astra`. astra also sets the `ownerReferences` for such objects to the underlying astra component so that when you do `astra delete`, such resources are deleted from the cluster. Other than that, astra doesn't help in managing such resources and users are expected to know how to do so.