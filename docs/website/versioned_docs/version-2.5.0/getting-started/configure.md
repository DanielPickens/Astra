---
title: Configuration
sidebar_position: 6
---
# Configuring astra global settings

The global settings for astra can be found in `preference.yaml` file; which is located by default in the `.astra` directory of the user's HOME directory.

Example:

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs
defaultValue="linux"
values={[
{label: 'Linux', value: 'linux'},
{label: 'Windows', value: 'windows'},
{label: 'Mac', value: 'mac'},
]}>

<TabItem value="linux">

```sh
/home/userName/.astra/preference.yaml
```

</TabItem>

<TabItem value="windows">

```sh
C:\\Users\userName\.astra\preference.yaml
```

</TabItem>

<TabItem value="mac">

```sh
/Users/userName/.astra/preference.yaml
```

</TabItem>

</Tabs>

---
A  different location can be set for the `preference.yaml` by exporting `GLOBALastraCONFIG` in the user environment.

### View the configuration
To view the current configuration, run `astra preference view`.

```shell
astra preference view
```
Example:
```shell
$ astra preference view
PARAMETER             CURRENT_VALUE
UpdateNotification
Timeout
PushTimeout
RegistryCacheTime
Ephemeral
ConsentTelemetry
```
### Set a configuration
To set a value for a preference key, run `astra preference set <key> <value>`.
```shell
astra preference set updatenotification false
```
Example:
```shell
$ astra preference set updatenotification false
Global preference was successfully updated
```
Note that the preference key is case-insensitive.

### Unset a configuration
To unset a value of a preference key, run `astra preference unset <key>`; use `-f` flag to skip the confirmation.
```shell
astra preference unset updatednotification
```
Example:
```shell
$ astra preference unset updatednotification
? Do you want to unset updatenotification in the preference (y/N) y
Global preference was successfully updated
```

Unsetting a preference key sets it to an empty value in the preference file. astra will use the [default value](./configure#preference-key-table) for such configuration.

### Preference Key Table

| Preference         | Description                                                                    | Default                |
|--------------------|--------------------------------------------------------------------------------|------------------------|
| UpdateNotification | Control whether a notification to update astra is shown                          | True                   |
| NamePrefix         | Set a default name prefix for an astra resource (component, storage, etc)        | Current directory name |
| Timeout            | Timeout for Kubernetes server connection check                                 | 1 second               |
| PushTimeout        | Timeout for waiting for a component to start                                   | 240 seconds            |
| RegistryCacheTime  | For how long (in minutes) astra will cache information from the Devfile registry | 4 Minutes              |
| Ephemeral          | Control whether astra should create a emptyDir volume to store source code       | True                   |
| ConsentTelemetry   | Control whether astra can collect telemetry for the user's astra usage             | False                  |
