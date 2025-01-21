---
sidebar_position: 1
title: Introduction
---

### What is astra?

`astra` is a fast, iterative and straightforward CLI tool for developers who write, build, and deploy applications on Kubernetes.

We abstract the complex concepts of Kubernetes so you can focus on one thing: `code`.

Choose your favourite framework and `astra` will deploy it *fast* and *often* to your container orchestrator cluster.

`astra` is focused on [inner loop](/docs/introduction#what-is-inner-loop-and-outer-loop) development as well as tooling that would help users transition to the [outer loop](/docs/introduction#what-is-inner-loop-and-outer-loop).

Brendan Burns, one of the co-founders of Kubernetes, said in the [book Kubernetes Patterns](https://www.redhat.com/cms/managed-files/cm-oreilly-kubernetes-patterns-ebook-f19824-201910-en.pdf):

> It (Kubernetes) is the foundation on which applications will be built, and it provides a large library of APIs and tools for building these applications, but it does little to provide the application or container developer with any hints or guidance for how these various pieces can be combined into a complete, reliable system that satisfies their business needs and goals.

`astra` satisfies that need by making Kubernetes development *super easy* for application developers and cloud engineers.

### What is "inner loop" and "outer loop"?

The **inner loop** consists of local coding, building, running, and testing the application -- all activities that you, as a developer, can control. 

The **outer loop** consists of the larger team processes that your code flows through on its way to the cluster: code reviews, integration tests, security and compliance, and so on. 

The inner loop could happen mostly on your laptop. The outer loop happens on shared servers and runs in containers, and is often automated with continuous integration/continuous delivery (CI/CD) pipelines. 

Usually, a code commit to source control is the transition point between the inner and outer loops.

*([Source](https://developers.redhat.com/blog/2020/06/16/enterprise-kubernetes-development-with-astra-the-cli-tool-for-developers#improving_the_developer_workflow))*

### Why should I use `astra`?

You should use `astra` if:
* You love frameworks such as Node.js, Spring Boot or dotNet
* Your application is intended to run in a Kubernetes-like infrastructure
* You don't want to spend time fighting with DevOps and learning Kubernetes in order to deploy to your enterprise infrastructure

If you are an application developer wishing to deploy to Kubernetes easily, then `astra` is for you.

### How is astra different from `kubectl` and `oc`?

Both [`kubectl`](https://github.com/kubernetes/kubectl) and [`oc`](https://github.com/openshift/oc/) require deep understanding of Kubernetes and OpenShift concepts.

`astra` is different as it focuses on application developers and cloud engineers. Both `kubectl` and `oc` are DevOps oriented tools and help in deploying applications to and maintaining a Kubernetes cluster provided you know Kubernetes well.

`astra` is not meant to:
* Maintain a production Kubernetes cluster
* Perform sysadmin tasks against a Kubernetes cluster
