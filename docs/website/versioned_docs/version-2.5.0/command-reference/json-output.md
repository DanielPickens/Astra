---
title: JSON Output
sidebar_position: 100
---

The `astra` commands that output some content generally accept a `-o json` flag to output this content in a JSON format, suitable for other programs to parse this output more easily.

The output structure is similar to Kubernetes resources, with `kind`, `apiVersion`, `metadata` ,`spec` and `status` fields.

List commands return a `List` resource, containing an `items` (or similar) field listing the items of the list, each item being also similar to Kubernetes resources.

Delete commands return a `Status` resource; see the [Status Kubernetes resource](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/status/).

Other commands return a resource associated with the command (`Application`, `Storage`', `URL`, etc).

The exhaustive list of commands accepting the `-o json` flag is currently:

| commands                       | Kind (version)                          | Kind (version) of list items                                 | Complete content?         | 
|--------------------------------|-----------------------------------------|--------------------------------------------------------------|---------------------------|
| astra application describe       | Application (astra.dev/v1alpha1)          | *n/a*                                                        |no                         |
| astra application list           | List (astra.dev/v1alpha1)                 | Application (astra.dev/v1alpha1)                               | ?                         |
| astra catalog list components    | List (astra.dev/v1alpha1)                 | *missing*                                                    | yes                       |
| astra catalog list services      | List (astra.dev/v1alpha1)                 | ClusterServiceVersion (operators.coreos.com/v1alpha1)        | ?                         |
| astra catalog describe component | *missing*                               | *n/a*                                                        | yes                       |
| astra catalog describe service   | CRDDescription (astra.dev/v1alpha1)       | *n/a*                                                        | yes                       |
| astra component create           | Component (astra.dev/v1alpha1)            | *n/a*                                                        | yes                       |
| astra component describe         | Component (astra.dev/v1alpha1)            | *n/a*                                                        | yes                       |
| astra component list             | List (astra.dev/v1alpha1)                 | Component (astra.dev/v1alpha1)                                 | yes                       |
| astra config view                | DevfileConfiguration (astra.dev/v1alpha1) | *n/a*                                                        | yes                       |
| astra debug info                 | astraDebugInfo (astra.dev/v1alpha1)         | *n/a*                                                        | yes                       |
| astra env view                   | EnvInfo (astra.dev/v1alpha1)              | *n/a*                                                        | yes                       |
| astra preference view            | PreferenceList (astra.dev/v1alpha1)       | *n/a*                                                        | yes                       |
| astra project create             | Project (astra.dev/v1alpha1)              | *n/a*                                                        | yes                       |
| astra project delete             | Status (v1)                             | *n/a*                                                        | yes                       |
| astra project get                | Project (astra.dev/v1alpha1)              | *n/a*                                                        | yes                       |
| astra project list               | List (astra.dev/v1alpha1)                 | Project (astra.dev/v1alpha1)                                   | yes                       |
| astra registry list              | List (astra.dev/v1alpha1)                 | *missing*                                                    | yes                       |
| astra service create             | Service                                 | *n/a*                                                        | yes                       |
| astra service describe           | Service                                 | *n/a*                                                        | yes                       |
| astra service list               | List (astra.dev/v1alpha1)                 | Service                                                      | yes                       |
| astra storage create             | Storage (astra.dev/v1alpha1)              | *n/a*                                                        | yes                       |
| astra storage delete             | Status (v1)                             | *n/a*                                                        | yes                       |
| astra storage list               | List (astra.dev/v1alpha1)                 | Storage (astra.dev/v1alpha1)                                   | yes                       |
| astra url list                   | List (astra.dev/v1alpha1)                 | URL (astra.dev/v1alpha1)                                       | yes                       |
