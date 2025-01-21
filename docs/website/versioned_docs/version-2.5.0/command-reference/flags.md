---
title: Common Flags
sidebar_position: 50
---

### Available Flags
Following are the flags commonly available with almost every astra command.
* `--context` - Use this flag to set the context directory where the component is defined.
* `--project` - Use this flag to set the project for the component; defaults to the project defined in the  local configuration; if none is available, then current project on the cluster
* `--app` - Use this flag to set the application of the component; defaults to the application defined in the local configuration; if none is available, then _app_
* `--kubeconfig` - Use this flag to set path to the kubeconfig if not using the default configuration
* `--show-log` - Use this flag to see the logs
*  `-f`, `--force` - Use this flag to tell the command not to prompt user for confirmation
* `-v`, `--v` - Use this flag to set the verbosity level. See [Logging in astra](https://github\.com/danielpickens/astra/wiki/Dev:-Logging-in-astra) for more information.
* `-h`, `--help` - Use this flag to get help on a command

**Note:** Some flags might not be available in some commands, run the command with `--help` to get a list of all the available flags.