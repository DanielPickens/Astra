---
title: astra list
---

`astra list` command combines the [`astra list binding`](./list-binding.md) and [`astra list component`](./list-component.md) commands.

## Running the command

```shell
astra list
```

<details>
<summary>Example</summary>

```shell
$ astra list
 ✓  Listing components from namespace 'my-percona-server-mongodb-operator' [292ms]
 NAME              PROJECT TYPE  RUNNING IN  MANAGED                          PLATFORM
 * my-nodejs       nodejs        Deploy      astra (v3.7)                       cluster
 my-go-app         go            Dev         astra (v3.7)                       podman
 mongodb-instance  Unknown       None        percona-server-mongodb-operator  cluster

Bindings:
 NAME                        APPLICATION                 SERVICES                                                   RUNNING IN 
 my-go-app-mongodb-instance  my-go-app-app (Deployment)  mongodb-instance (PerconaServerMongoDB.psmdb.percona.com)  Dev
```
</details>
