`Astra` - Fast, Iterative and Simplified container-based application development
---

[![GitHub release](https://img.shields.io/github/v/release/daniel-pickens/astra?style=for-the-badge)](https://github.com/daniel-pickens/astra/releases/latest)
![License](https://img.shields.io/github/license/daniel-pickens/astra?style=for-the-badge)
[![Gastrac](https://img.shields.io/badge/gastrac-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://astra.dev/gastrac)
[![Netlify Status](https://api.netlify.com/api/v1/badges/e07867b0-56a4-4905-92a9-a152ceab5f0d/deploy-status)](https://app.netlify.com/sites/astra-docusaurus-preview/deploys)

![logo](/docs/website/static/img/logo_small.png)

### Overview

`Astra` is a fast, and iterative CLI tool for container-based application development.
It is an implementation of the open [Devfile](https://devfile.io/) standard, supporting [Podman](https://podman.io/), [Kubernetes](https://kubernetes.io/) and [OpenShift](https://www.redhat.com/en/technologies/cloud-computing/openshift).
Key features:
* builds linux kernel from SUSE kernal, 
* builds a docker image with the kernel and run it in a container, 
* checks if the api's running in the pod are accessible from the host machine.
* checks speed of the api's running in the pod that are accessible from the host machine.
* writes checks to stdout log file for security hardening and diagnosis

**Why use `Astra`?**

* **Easy onboarding:** By auto-detecting the project source code, you can easily get started with `Astra`.
* **No cluster needed**: With Podman support, having a Kubernetes cluster is not required to get started with `Astra`. Using a common abstraction, `Astra` can run your application on Podman, Kubernetes or OpenShift.
* **Fast:** Spend less time maintaining your application deployment infrastructure and more time coding. Immediately have your application running each time you save.
* **Standalone:** `Astra` is a standalone tool that communicates directly with the container orchestrator API.
* **No configuration needed:** There is no need to dive into complex Kubernetes YAML configuration files. `Astra` abstracts those concepts away and lets you focus on what matters most: code.
* **Containers first:** We provide first class support for Podman, Kubernetes and OpenShift. Choose your favourite container orchestrator and develop your application.
* **Easy to learn:** Simple syntax and design centered around concepts familiar to developers, such as projects, applications, and components.

Learn more about the features provided by `Astra` on [astra.dev](https://astra.dev/docs/overview/features).


### Installing `Astra`

Please check the [installation guide on astra.dev](https://astra.dev/docs/overview/installation/).

### Official documentation

Visit [astra.dev](https://astra.dev/) to learn more about astra.

### Community, discussion, contribution, and support

#### Chat 
All of the developer and user discussions happen in the #Astra channel on the official Kubernetes Slack.

If you haven't already joined the Kubernetes Slack, you can invite yourself here.

Ask questions, inquire about Astra or even discuss a new feature:
https://app.slack.com/client/T09NY5SBT
#### Issues

If you find an issue with `Astra`, please [file it here](https://github.com/danielpickens/astra/issues).

#### Contributing

* Code: I am currently working on updating the code contribution guide.
* Documentation: To contribute to the documentation, please have a look at our [Documentation Guide](https://github.com/daniel-pickens/astra/wiki).

Astra houses an open community who welcomes any concerns, changes or ideas for `astra`! Come join the chat and hang out, ask or give feedback and just generally have a good time.

### Legal

#### License

Unless otherwise stated (ex. `/vendor` files), all code is licensed under the [Apache 2.0 License](LICENSE). 

#### Usage data
