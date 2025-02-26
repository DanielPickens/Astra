---
title: astra list component
---

`astra list component` command is useful for getting information about components running on a specific namespace of a cluster or on Podman.

If the command is executed from a directory containing a Devfile, it also displays the component
defined in the Devfile as part of the list, prefixed with a star(*).

For each component, the command displays:
- its name,
- its project type,
- on which mode it is running (None, Dev, Deploy, or both), note that None is only applicable to the component 
defined in the local Devfile,
- by which application the component has been deployed,
- the platform on which the component is running (cluster or podman).

### Running the command
```shell
astra list component
```
<details>
<summary>Example</summary>

```shell
$ astra list component
 ✓  Listing components from namespace 'my-percona-server-mongodb-operator' [292ms]
 NAME              PROJECT TYPE  RUNNING IN  MANAGED                          PLATFORM
 * my-nodejs       nodejs        Deploy      astra (v3.7)                       cluster
 my-go-app         go            Dev         astra (v3.7)                       podman
 mongodb-instance  Unknown       None        percona-server-mongodb-operator  cluster
```
</details>

### Targeting a specific platform

By default, `astra list component` will search components in both the current namespace of the cluster and podman. You can restrict the search to one of the platforms only, using the `--platform` flag, giving a value `cluster` or `podman`.

:::tip use of cache

`astra list component` makes use of cache for performance reasons. This is the same cache that is referred by `kubectl` command 
when you do `kubectl api-resources --cached=true`. As a result, if you were to install an Operator/CRD on the 
Kubernetes cluster, and create a resource from it using `astra`, you might not see it in the `astra list component` output. This 
would be the case for 10 minutes timeframe for which the cache is considered valid. Beyond this 10 minutes, the 
cache is updated anyway.

If you would like to invalidate the cache before the 10 minutes timeframe, you could manually delete it by doing:
```shell
rm -rf ~/.kube/cache/discovery/api.crc.testing_6443/
```
Above example shows how to invalidate the cache for a CRC cluster. Note that you will have to modify the `api.crc.
testing_6443` part based on the cluster you are working against.