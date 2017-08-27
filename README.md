# powerline-go

<img align="center" padding="5px" alt=":)" width="834px" src="/screenshot.png"/>

Attempted fork of [Otann/powerline-go](https://github.com/Otann/powerline-go)

This application does not cover all features of powerline-shell, only those that I currently use.
 
For now it is only configurable trough the source recompilation, which is quite fast with Go.

## Usage

If you haven't install go, please fallow [this](https://golang.org/doc/install).
    
Then install the binary with

    go get github.com/sanyatuning/powerline-go
    go install github.com/sanyatuning/powerline-go

### Bash

Install powerline-go and add the following to your `~/.bashrc`

    function _update_ps1() {
       export PS1="$(powerline-go bash $? ${COLUMNS} 2> /dev/null)"
    }

    export PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"

### Zsh

Install powerline-go and add the following to your `~/.zshrc`

    function powerline_precmd() {
      export PS1="$(powerline-go zsh $? ${COLUMNS} 2> /dev/null)"
    }

    function install_powerline_precmd() {
      for s in "${precmd_functions[@]}"; do
        if [ "$s" = "powerline_precmd" ]; then
          return
        fi
      done
      precmd_functions+=(powerline_precmd)
    }

    install_powerline_precmd

## Performance

TODO: update

```
$ time ~/git/milkbikis/powerline-shell/powerline-shell.py > /dev/null
real    0m0.092s
user    0m0.027s
sys     0m0.046s
```

```
$ time ~/go/src/github.com/sivel/powerline-shell/powerline-shell > /dev/null
real    0m0.007s
user    0m0.002s
sys     0m0.004s
```
