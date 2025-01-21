Usage Data
---

You can help improve `astra` by allowing it to collect usage data.


If the user has consented to `astra` collecting usage data, the following data will be collected when a command is executed -

* Command Name
* Command Duration
* Command Success
* Pseudonymized error message and error type (in case of failure)
* Whether the command was run from a terminal
* Whether the command was run in experimental mode
* `astra` version in use

In addition to this, the following data about user's identity is also noted - 
* OS type
* Timezone
* Locale

The following tables describe the additional information collected by `astra` commands.

**astra v3**

| Command                           | Data                                                                                                            |
|-----------------------------------|-----------------------------------------------------------------------------------------------------------------|
| astra init                          | Component Type, Devfile Name, Language, Project Type, Interactive Mode (bool)                                   |
| astra dev                           | Component Type, Devfile Name, Language, Project Type, Platform (podman, kubernetes, openshift), Platform version|
| astra deploy                        | Component Type, Devfile Name, Language, Project Type, Platform (kubernetes, openshift), Platform version        |
| astra <create/set/delete> namespace | Cluster Type (Possible values: OpenShift 3, OpenShift 4, Kubernetes)                                            |

**astra v3 GUI**

The astra v3 GUI is accessible (by default at http://localhost:20000) when the command `astra dev` is running.

| Page                 | Data
|----------------------|-------------------------
| YAML (main page)     | Page accessed, UI started, Devfile saved to disk, Devfile cleared, Devfile applied |
| Metadata             | Page accessed, Metadata applied |
| Commands             | Page accessed, Start create command, Create command |
| Events               | Page accessed, Add event |
| Containers           | Page accessed, Create container |
| Images               | Page accessed, Create Image |
| Resources            | Page accessed, Create Resource |

**astra v2**

| Command                  | Data                                                                 |
|--------------------------|----------------------------------------------------------------------|
| astra create               | Component Type, Devfile name                                         |
| astra push                 | Component Type, Cluster Type, Language, Project Type                 |
| astra project <create/set> | Cluster Type (Possible values: OpenShift 3, OpenShift 4, Kubernetes) |


All the data collected above is pseudonymized to keep the user information anonymous.

Note: Telemetry data is not collected when you run `--help` for commands.

###  Enable/Disable preference

#### Enable
`astra preference set ConsentTelemetry true`

#### Disable
`astra preference set ConsentTelemetry false`

Alternatively you can _disable_ telemetry by setting the `astra_TRACKING_CONSENT` environment variable to `no`.
This environment variable will override the `ConsentTelemetry` value set by `astra preference`.
