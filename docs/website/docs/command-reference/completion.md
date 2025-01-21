---
title: astra completion
---

## Description

`astra completion` is used to generate shell completion code. The generated code provides interactive shell completion code for `astra`.

There is support for the following terminal shells:
- [Bash](https://www.gnu.org/software/bash/)
- [Zsh](https://zsh.sourceforge.io/)
- [Fish](https://fishshell.com/)
- [Powershell](https://docs.microsoft.com/en-us/powershell/)

## Running the Command

To generate the shell completion code, the command can be ran as follows:

```sh
astra completion [SHELL]
```

### Bash

Load into your current shell environment:

```sh
source <(astra completion bash)
```

Load persistently:

```sh
# Save the completion to a file
astra completion bash > ~/.astra/completion.bash.inc

# Load the completion from within your $HOME/.bash_profile
source ~/.astra/completion.bash.inc
```

### Zsh

Load into your current shell environment:

```sh
source <(astra completion zsh)
```

Load persistently:

```sh
astra completion zsh > "${fpath[1]}/_astra"
```

### Fish

Load into your current shell environment:

```sh
source <(astra completion fish)
```

Load persistently:

```sh
astra completion fish > ~/.config/fish/completions/astra.fish
```

### Powershell

Load into your current shell environment:

```sh
astra completion powershell | Out-String | Invoke-Expression
```

Load persistently:

```sh
astra completion powershell >> $PROFILE
```
