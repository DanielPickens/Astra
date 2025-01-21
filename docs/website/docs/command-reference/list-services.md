---
title: astra list services
---

## Description

You can use `astra list services` to list all the bindable Operator backed services on the cluster.

## Running the command

To list bindable services in the current project/namespace:
```shell
astra list services
```
<details>
<summary>Example</summary>

```shell
$ astra list services
 ✓  Listing bindable services from namespace "myproject" [82ms]

 NAME                                                  NAMESPACE 
 redis-standalone/Redis.redis.redis.opstreelabs.in/v1  myproject 
```
</details>

To list bindable services in all projects/namespaces accessible to the user:
```shell
astra list services -A 
```
<details>
<summary>Example</summary>

```shell
astra list services -A
 ✓  Listing bindable services from all namespaces [182ms]

 NAME                                                  NAMESPACE  
 redis-standalone/Redis.redis.redis.opstreelabs.in/v1  myproject  
 hello-world/RabbitmqCluster.rabbitmq.com/v1           newproject 
```
</details>

To list bindable services in a particular project/namespace that is accessible to the user:
```shell
astra list services -n <project-name>
```
<details>
<summary>Example</summary>

```shell
$ astra list services -n newproject
 ✓  Listing bindable services from namespace "newproject" [45ms]

 NAME                                         NAMESPACE  
 hello-world/RabbitmqCluster.rabbitmq.com/v1  newproject 
```
</details>

To get the JSON formatted output for any of the above commands, add `-o json` to the commands shown above. That 
would be:
* `astra list services -o json`
* `astra list services -A -o json`
* `astra list services -n <project-name> -o json`