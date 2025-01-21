---
title: Create Component
sidebar_position: 1
sidebar_label: Creating components
---

# Creating components using astra

[Component](../getting-started/basics#component) is the most basic unit of operation for astra. And the way to create one is using `astra create` (short for `astra component create`) command.

In simplest terms, when you "create" an astra component, you populate your current working directory with the file `devfile.yaml`. A Devfile is a manifest file that contains information about various resources (URL, Storage, Services, etc.) that correspond to your component, and will be created on the Kubernetes cluster when you execute `astra push` command. Most astra commands will first modify (add or remove configuration from) this file, and then subsequent `astra push` will create or delete the resources from the Kubernetes cluster.

However, astra users are not expected to know how the `devfile.yaml` is organized; it is the astra commands that would create, update, or delete it.

One final thing to keep in mind - there can be only one astra component in a directory. Nesting astra components is not expected to work well. In other terms, if you have multiple parts (components), say frontend and backend, of your microservices application that you want to create astra components for, you should put them in separate directories and not try to nest them. Take a look at example structure below:
```shell
$ tree my-awesome-microservices-app 
my-awesome-microservices-app
├── backend
│   └── devfile.yaml
└── frontend
    └── devfile.yaml
```
In this guide, we are going to create a Spring Boot component and a Nodejs component to deploy parts of the [astra quickstart](https://github.com/dharmit/astra-quickstart) project to a Kubernetes cluster.

Let's clone the project first:
```shell
git clone https://github.com/dharmit/astra-quickstart
cd astra-quickstart
```

Next, create a project <!-- add link to project command reference here --> on the Kubernetes cluster in which we will be creating our component. This is to keep our Kubernetes cluster clean of the tasks we perform (this step is optional):
```shell
astra project create myproject
```
Alternatively, you could also use one of the existing projects on the cluster:
```shell
astra project list
```
Now, set the project in which you want to create the component:
```shell
# replace <project-name> with a valid value from the list
astra project set <project-name>
```

astra supports interactive and non-interactive ways of creating a component. We will create the Spring Boot component interactively and Nodejs component non-interactively. The Spring Boot component is in `backend` directory. It contains code for the REST API that our Nodejs component will be interacting with. This Nodejs component is in `frontend` directory.

## Creating a component interactively

To interactively create the Spring Boot component, `cd` into the cloned project (already done if you copy-pasted the command above), then `cd` into `backend` directory, and execute:
```shell
cd backend
astra create
```
You will be prompted with a few questions one after the another. Go through each one of them to create a component.

1. First question is about selecting the component type:
    ```shell
    $ astra create
    ? Which devfile component type do you wish to create  [Use arrows to move, enter to select, type to filter]
    > java-maven
    java-maven
    java-openliberty
    java-openliberty
    java-quarkus
    java-quarkus
    java-springboot
    ```
   By default, `java-maven` is selected for us. Since this is a Spring Boot application, we should be selecting `java-springboot`. 

    We can either scroll down to `java-springboot` using the arrow key, or start typing `spring` on the prompt. Typing `spring` will lead to astra filtering the component type based on your input.

2. Next, astra asks you to name the component:
    ```shell
    $ astra create                
    ? Which devfile component type do you wish to create java-springboot
    ? What do you wish to name the new devfile component (java-springboot) backend
    ```
    Name it `backend`.

3. Next, astra asks you for the project in which you would like to create the component. Use the project `myproject` that we created earlier or the one you had set using `astra project set` command
   ```shell
   $ astra create
   ? Which devfile component type do you wish to create java-springboot
   ? What do you wish to name the new devfile component java-springboot
   ? What project do you want the devfile component to be created in myproject
   ```
   Now you will have a `devfile.yaml` in your current working directory. But astra is not done asking you questions yet.
4. Lastly, astra asks you if you would like to download a "starter project". Since we already cloned the astra-quickstart project, we answer in No by typing `n` and hitting the return key. We discuss starter projects later in [this document](#starter-projects):
   ```shell
   $ astra create
   ? Which devfile component type do you wish to create java-springboot
   ? What do you wish to name the new devfile component java-springboot
   ? What project do you want the devfile component to be created in myproject
   Devfile Object Validation
   ✓  Checking devfile existence [66186ns]
   ✓  Creating a devfile component from registry: stage [92202ns]
   Validation
   ✓  Validating if devfile name is correct [99609ns]
   ? Do you want to download a starter project (y/N) n
   ```
   
Your Spring Boot component is now ready for use.

## Creating a component non-interactively

To non-interactively create the Nodejs component to deploy our frontend code, `cd` into the cloned `frontend` directory and execute:
```shell
# assuming you are in the astra-quickstart/backend directory
cd ../frontend 
astra create nodejs frontend -n myproject
```
Here `nodejs` is the type of the component, `frontend` is the name of the component, and `-n myproject` tells astra to use the project `myproject` for the mentioned `astra create` operation.

## Starter projects

Besides creating a component for an existing code, you could also use "starter project" when creating a component.

Starter projects are example projects developed by the community to showcase the usability of devfiles. An astra user can use these starter projects by running `astra create` command in an empty directory.

### Starter projects in interactive mode

To interactively create a Java Spring Boot component using the starter project, you can follow the below steps:
```shell
mkdir myastraComponent && cd myastraComponent
astra create
```
In the questions that astra asks you next, provide answers like below:
```shell
$ astra create
? Which devfile component type do you wish to create java-springboot
? What do you wish to name the new devfile component myFirstComponent
? What project do you want the devfile component to be created in myproject
Devfile Object Validation
 ✓  Checking devfile existence [60122ns]
 ✓  Creating a devfile component from registry: stage [91411ns]
Validation
 ✓  Validating if devfile name is correct [35749ns]
? Do you want to download a starter project Yes

Starter Project
 ✓  Downloading starter project springbootproject from https://github.com/astra-devfiles/springboot-ex.git [716ms]

Please use `astra push` command to create the component with source deployed
```

### Starter projects in non-interactive mode

To non-interactively create a Java Spring Boot component using the starter project, you can follow the below steps:
```shell
mkdir myastraComponent && cd myastraComponent
astra create java-springboot myFirstComponent --starter springbootproject
```

## Push the component to Kubernetes

astra follows a "create & push" workflow for almost all the commands. Meaning, most astra commands won't create resources on Kubernetes cluster unless you run `astra push` command.

Among the various ways described above, irrespective of how you created the component, the next step to create the resources for our component on the cluster would be to run `astra push`.

Execute below command from the component directory of both the `frontend` and `backend` components:
```shell
astra push
```