---
title: Installation
sidebar_position: 4
toc_min_heading_level: 2
toc_max_heading_level: 4
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

`astra` can be used as either a [CLI tool](#cli-installation) or an [IDE plugin](#ide-installation) on [Mac](#macos), [Windows](#windows) or [Linux](#linux).

Each release is *signed*, *checksummed*, *verified*, and then pushed to our [binary mirror](https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/).

For more information on the changes of each release, they can be viewed either on [GitHub](https://github\.com/danielpickens/astra/releases) or the [blog](/blog).

## CLI Installation

### Linux

<Tabs
defaultValue="amd64"
values={[
{label: 'Intel / AMD 64', value: 'amd64'},
{label: 'ARM 64', value: 'arm64'},
{label: 'PowerPC', value: 'ppc64le'},
{label: 'IBM Z', value: 's390x'},
]}>

<TabItem value="amd64">

Installing `astra` on `amd64` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-amd64 -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-amd64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="arm64">

Installing `astra` on `arm64` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-arm64 -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-arm64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="ppc64le">

Installing `astra` on `ppc64le` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-ppc64le -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-ppc64le.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="s390x">

Installing `astra` on `s390x` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-s390x -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-linux-s390x.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

</Tabs>

---

### MacOS

#### Homebrew

**NOTE:** This will install from the *main* branch on GitHub

Installing `astra` using [Homebrew](https://brew.sh/):

1. Install astra:

```shell
brew install astra-dev
```

2. Verify the version you installed is up-to-date:

```shell
astra version
```

#### Binary

<Tabs
defaultValue="intel"
values={[
{label: 'Intel', value: 'intel'},
{label: 'Apple Silicon', value: 'arm'},
]}>

<TabItem value="intel">

Installing `astra` on `amd64` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-darwin-amd64 -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-darwin-amd64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
chmod +x ./astra
sudo mv ./astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="arm">

Installing `astra` on `arm64` architecture:

1. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-darwin-arm64 -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-darwin-arm64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
chmod +x ./astra
sudo mv ./astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

</Tabs>

---

### Windows

1. Open a PowerShell terminal

2. Download the latest release from the mirror:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-windows-amd64.exe -o astra.exe
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://developers.redhat.com/content-gateway/rest/mirror/pub/openshift-v4/clients/astra/v3.16.1/astra-windows-amd64.exe.sha256 -o astra.exe.sha256
# Visually compare the output of both files
Get-FileHash astra.exe
type astra.exe.sha256
```

4. Add the binary to your `PATH`

## IDE Installation

### Visual Studio Code (VSCode)

The [OpenShift Toolkit](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-openshift-connector) VSCode extension uses both `astra` and `oc` binaries
to interact with Kubernetes or OpenShift cluster.

To install it:

1. Open VS Code.
2. Launch VS Code **Quick Open** (`Ctrl+P`).
3. Paste the following command and press `Enter`:

```
ext install redhat.vscode-openshift-connector
```

### JetBrains IDEs
The [OpenShift Toolkit by Red Hat](https://plugins.jetbrains.com/plugin/12030-openshift-toolkit-by-red-hat/) plugin can be installed
to interact with OpenShift or Kubernetes clusters right from your JetBrains IDEs like IntelliJ IDEA, WebStorm or PyCharm.
It uses `astra` and `oc` binaries for fast iterative application development on those clusters.

To install it:

1. Press `Ctrl+Alt+S` to open the IDE settings and select **Plugins**.
2. Find the "**OpenShift Toolkit by Red Hat**" plugin in the **Marketplace** and click **Install**.

## Alternative installation methods

#### Source code
1. Clone the repository and cd into it.
   ```shell
   git clone https://github\.com/danielpickens/astra.git
   cd astra
   ```
2. Install tools used by the build and test system.
   ```shell
   make goget-tools
   ```
3. Build the executable from the sources in `cmd/astra`.
   ```shell
   make bin
   ```
4. Check the build version to verify that it was built properly.
   ```shell
   ./astra version
   ```
5. Install the executable in the system's GOPATH.
   ```shell
   make install
   ```
6. Check the binary version to verify that it was installed properly; verify that it is same as the build version.
   ```shell
   astra version
   ```

### Maven plugin
It is possible to integrate the `astra` binary download in a Maven project using [astra Downloader Plugin](https://github.com/tnb-software/astra-downloader).
The download can be executed using the `download` goal which automatically retrieves the version for the current architecture:
```shell
mvn software.tnb:astra-downloader-maven-plugin:0.1.3:download \
  -Dastra.target.file=$HOME/bin/astra \
  -Dastra.version=v3.16.1
```

### asdf
The [asdf version manager](https://asdf-vm.com/) is a tool for managing multiple runtime versions using a common CLI.
With `asdf` installed, the [asdf plugin for astra](https://github.com/rm3l/asdf-astra) can be used to install any released version of `astra`:
```
asdf plugin add astra
asdf install astra 3.16.1
asdf global astra 3.16.1
```

## Nightly builds

Nightly builds of `astra` are also available. Note that these builds are provided as is and can be highly unstable.

### Linux

<Tabs
defaultValue="amd64"
values={[
{label: 'Intel / AMD 64', value: 'amd64'},
{label: 'ARM 64', value: 'arm64'},
{label: 'PowerPC', value: 'ppc64le'},
{label: 'IBM Z', value: 's390x'},
]}>

<TabItem value="amd64">

Installing `astra` on `amd64` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-amd64 -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-amd64-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-amd64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-amd64-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="arm64">

Installing `astra` on `arm64` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-arm64 -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-arm64-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-arm64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-arm64-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="ppc64le">

Installing `astra` on `ppc64le` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-ppc64le -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-ppc64le-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-ppc64le.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-ppc64le-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="s390x">

Installing `astra` on `s390x` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-s390x -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-s390x-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-s390x.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-linux-s390x-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
sudo install -o root -g root -m 0755 astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

</Tabs>

---

### MacOS

<Tabs
defaultValue="intel"
values={[
{label: 'Intel', value: 'intel'},
{label: 'Apple Silicon', value: 'arm'},
]}>

<TabItem value="intel">

Installing `astra` on `amd64` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-amd64 -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-amd64-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-amd64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-amd64-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
chmod +x ./astra
sudo mv ./astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

<TabItem value="arm">

Installing `astra` on `arm64` architecture:

1. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-arm64 -o astra
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-arm64-${gitCommitId} -o astra
```

2. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-arm64.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-darwin-arm64-${gitCommitId}.sha256 -o astra.sha256
echo "$(<astra.sha256)  astra" | shasum -a 256 --check
```

3. Install astra:
```shell
chmod +x ./astra
sudo mv ./astra /usr/local/bin/astra
```

4. (Optional) If you do not have root access, you can install `astra` to the local directory and add it to your `$PATH`:

```shell
mkdir -p $HOME/bin 
cp ./astra $HOME/bin/astra
export PATH=$PATH:$HOME/bin
# (Optional) Add the $HOME/bin to your shell initialization file
echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
```
</TabItem>

</Tabs>

---

### Windows

1. Open a PowerShell terminal

2. Download the latest nightly build:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-windows-amd64.exe -o astra.exe
```

To download a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-windows-amd64-${gitCommitId}.exe -o astra.exe
```


3. (Recommended) Verify the downloaded binary with the SHA-256 sum:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-windows-amd64.exe.sha256 -o astra.exe.sha256
# Visually compare the output of both files
Get-FileHash astra.exe
type astra.exe.sha256
```

To verify a download for a specific commit instead, run:
```shell
curl -L https://s3.us-east.cloud-object-storage.appdomain.cloud/astra-nightly-builds/astra-windows-amd64-${gitCommitId}.exe.sha256 -o astra.exe.sha256
# Visually compare the output of both files
Get-FileHash astra.exe
type astra.exe.sha256
```

4. Add the binary to your `PATH`
