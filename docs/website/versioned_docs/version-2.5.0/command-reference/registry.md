---
title: astra registry
sidebar_position: 8
---

astra uses the portable *devfile* format to describe the components. astra can connect to various devfile registries to download devfiles for different languages and frameworks.

You can connect to publicly available devfile registries, or you can install your own [Secure Registry](../architecture/secure-registry).

You can use the `astra registry` command to manage the registries used by astra to retrieve devfile information.

## Listing the registries

You can use the following command to list the registries currently contacted by astra:

```
astra registry list
```

For example:

```
$ astra registry list
NAME                       URL                             SECURE
DefaultDevfileRegistry     https://registry.devfile.io     No
```

`DefaultDevfileRegistry` is the default registry used by astra; it is provided by the [devfile.io](https://devfile.io) project.

## Adding a registry

You can use the following command to add a registry:

```
astra registry add
```

For example:

```
$ astra registry add StageRegistry https://registry.stage.devfile.io
New registry successfully added
```

If you are deploying your own Secure Registry, you can specify the personal access token to authenticate to the secure registry with the `--token` flag:

```
$ astra registry add MyRegistry https://myregistry.example.com --token <access_token>
New registry successfully added
```

## Deleting a registry

You can delete a registry with the command:

```
astra registry delete
```

For example:

```
$ astra registry delete StageRegistry
? Are you sure you want to delete registry "StageRegistry" Yes
Successfully deleted registry
```

You can use the `--force` (or `-f`) flag to force the deletion of the registry without confirmation.

## Updating a registry

You can update the URL and/or the personal access token of a registry already registered with the command:

```
astra registry update
```

For example:

```
$ astra registry update MyRegistry https://otherregistry.example.com --token <other_access_token>
? Are you sure you want to update registry "MyRegistry" Yes
Successfully updated registry
```

You can use the `--force` (or `-f`) flag to force the update of the registry without confirmation.

