---
title: Pushing Source Files
sidebar_position: 2
---

During the execution of `astra dev`, the files of the current directory and its sub-directories are pushed 
to the container. You can influence which files are pushed with the following methods.


## Pushing only specific files

`astra` uses the `dev.astra.push.path` related attribute from the devfile's run commands to push only the specified files and folders to the component.

The format of the attribute is `"dev.astra.push.path:<local_relative_path>": "<remote_relative_path>"`. We can mention multiple such attributes in the run command's `attributes` section.

```yaml
commands:
  - id: dev-run
    # highlight-start
    attributes:
      "dev.astra.push.path:target/quarkus-app": "remote-target/quarkus-app"
      "dev.astra.push.path:README.txt": "docs/README.txt"
    # highlight-end
    exec:
      component: tools
      commandLine: "java -jar remote-target/quarkus-app/quarkus-run.jar -Dquarkus.http.host=0.0.0.0"
      hotReloadCapable: true
      group:
        kind: run
        isDefault: true
      workingDir: $PROJECTS_ROOT
  - id: dev-debug
    # highlight-start
    attributes:
      "dev.astra.push.path:target/quarkus-app": "remote-target/quarkus-app"
      "dev.astra.push.path:README.txt": "docs/README.txt"
    # highlight-end
    exec:
      component: tools
      commandLine: "java -Xdebug -Xrunjdwp:server=y,transport=dt_socket,address=${DEBUG_PORT},suspend=n -jar remote-target/quarkus-app/quarkus-run.jar -Dquarkus.http.host=0.0.0.0"
      hotReloadCapable: true
      group:
        kind: debug
        isDefault: true
      workingDir: $PROJECTS_ROOT
```

In the above example the contents of the `quarkus-app` folder, which is inside the `target` folder, will be pushed to the remote location of `remote-target/quarkus-app` and the file `README.txt` will be pushed to `doc/README.txt`.
The local path is relative to the component's local folder. The remote location is relative to the folder containing the component's source code inside the container. 

## Ignoring files to push

`astra` excludes from the push the files present in the `.astraignore` file, or, if
this file does not exist, the files present in the `.gitignore` file.

By default, `astra` does not create a `.astraignore` file and relies on the `.gitignore` file.
Also, during each execution, `astra dev` adds the `.astra` entry to the `.gitignore` file if it is not already present in this file,
to avoid an infinite loop on the synchronization, this directory containing a file with the state of the sync.

If you want to use the `.astraignore` file instead, to have a different set of files ignored for sync and ignored for git, 
you will need to add the `.astra` directory to the `.astraignore` file.

