# Go Boilerplate

This is my current boilerplate for starting new golang projects. This Boilerplate is for building a small golang API server and contains a bunch of nice additions for deployment, as well as providing a nice project structure. This boilerplate uses the [`cobra`](https://github.com/spf13/cobra) CLI framework to provide a hooked-up-out-of-the-box configuration file. There is also a ready to go `Dockerfile` that builds a lightweight `alpine` image to snugly wrap your binary.

You will want to install `cobra` on your system if you don't already have it in order to add more commands. The boilerplate only has `serve` and `version`. Add commands with `cobra add {{ .Command.Name }}`.

To adopt this boilerplate you will need to change names in the following places:

- Once you clone this repo you will need to change the name of the folder to the name of your project.
- The `Makefile` controls most of the configuration. Just change the variables at the top to match your needs:

```Makefile
BINARY            = go-boilerplate
GITHUB_USERNAME   = jackzampolin
DOCKER_REPO       = quay.io/jackzampolin
VERSION           = v0.1.0
PORT 							= 3000
```

- Change the name of the `go-boilerplate.yaml` file to match the name of your project. You will want to make a copy it to `~/.{{ .Project.Name }}` to that you don't have to run every command with the `--config="~/.{{ .Project.Name }}"` flag.
- Fix the project name in the golang code:

```bash
$ make rename
```

### Build and Run

To build and run this boilerplate run:

```bash
$ make install
$ go-boilerplate
A boilerplate for an API written in Golang

Usage:
  go-boilerplate [command]

Available Commands:
  help        Help about any command
  serve       Runs the server
  version     Prints version information

Flags:
      --config string   config file (default is $HOME/.go-boilerplate.yaml)
  -h, --help            help for go-boilerplate
  -t, --toggle          Help message for toggle

Use "go-boilerplate [command] --help" for more information about a command.
```

### Docker

To build the docker image with appropriate tags:

```bash
$ make docker
```

To push the docker image to the configured repo:

```bash
$ make docker-push
```

To run the just built docker image with the local config loaded and the proper port exposed:

```bash
$ make docker-run
```

### Inspiration and Contributors:

https://github.com/silven/go-example/blob/master/Makefile
https://vic.demuzere.be/articles/golang-makefile-crosscompile/
