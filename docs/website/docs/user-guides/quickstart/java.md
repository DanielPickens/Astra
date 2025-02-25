---
title: Developing with Java (Spring Boot)
sidebar_position: 2
---

## Step 0. Creating the initial source code (optional)

import InitialSourceCodeInfo from './docs-mdx/initial_source_code_description.mdx';

<InitialSourceCodeInfo/>

For Java, we will use the [Spring Initializr](https://start.spring.io/) to generate the example source code:

1. Navigate to [start.spring.io](https://start.spring.io/).
2. Select **Maven** under **Project**.
3. Click on "Add" under "Dependencies".
4. Select "Spring Web".
5. Click "Generate" to generate and download the source code.

Finally, extract the downloaded source code archive in the 'quickstart-demo' directory.

Your source code has now been generated and created in the directory.

## Step 1. Preparing the target platform

import PreparingTargetPlatform from './docs-mdx/preparing_the_target_platform.mdx';

<PreparingTargetPlatform/>

## Step 2. Initializing your application (`astra init`)


import InitSampleOutput from './docs-mdx/java/java_astra_init_output.mdx';
import InitDescription from './docs-mdx/astra_init_description.mdx';

<InitDescription framework="Java (Spring Boot)" initout=<InitSampleOutput/> />

## Step 3. Developing your application continuously (`astra dev`)

import DevSampleOutput from './docs-mdx/java/java_astra_dev_output.mdx';
import DevPodmanSampleOutput from './docs-mdx/java/java_astra_dev_podman_output.mdx';

import DevDescription from './docs-mdx/astra_dev_description.mdx';

<DevDescription framework="Java (Spring Boot)" devout=<DevSampleOutput/> devpodmanout=<DevPodmanSampleOutput/> />

_You can now follow the [advanced guide](../advanced/deploy/java.md) to deploy the application to production._
