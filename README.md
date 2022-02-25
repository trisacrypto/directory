# Global TRISA VASP Directory

**TRISA implementation of the Global VASP Directory Service**

## Code Generation

To run the code-generation utilities in this repository, you'll need the following tools installed.

- [protoc](https://github.com/protocolbuffers/protobuf/releases)
- [go-bindata](https://github.com/kevinburke/go-bindata)

If you're on OS X - then the easiest option is to install these tools using Homebrew as follows:

```
$ brew install go-bindata protobuf
```

On the server-side, the Go code requires code generation for the API protocol buffers and email template data. From the project root:

```
$ go generate ./...
```

On the web-ui side, there is a bash script to generate the Javascript code from the protocol buffers, which is also set as a package script.

```
$ cd web
$ npm run protos
```

## Docker Compose

The simplest way to develop the GDS UI is to run the GDS service and Envoy gRPC-Proxy using `docker compose`. The containers directory has the configuration and Dockerfiles to build the images:

```
$ docker compose -f ./containers/docker-compose.yaml --profile=all build
```

Note the `--profile` flag, this allows you to build or run only some of the images or services to help facilitate development. For example, if you're working on the front-end, you probably only need to run `--profile=api`. The profiles are as follows:

- `all`: All services defined by docker-compose.yaml
- `gds`: Just the GDS server(s)
- `api`: The GDS server(s) and the Envoy Proxy for grpc-web
- `ui`: Both the gds-ui and the gds-admin-ui
- `user`: Just the gds-ui React app
- `admin`: Just the gds-admin-ui React app

Once the containers are built, you will need to supply a configuration `.env` file in the `./containers` directory (e.g. adjacent to the `docker-compose.yaml` file). Copy the template .env file as follows:

```
$ cp ./containers/.env.template ./containers/.env
```

You shouldn't have to change anything to run the GDS server, but if you'd like your local environment to test against an external service, please obtain developer credentials from one of the maintainers.

The fixtures directory has been configured with example/test keys in `fixtures/cred` and has some other required directories revisioned with `.gitkeep` such as `fixtures/backups` and `fixtures/certs` -- these folders are git ignored, but you'll see data appear in them when you're running the GDS service.

By default, GDS will create an empty `fixtures/db` directory for the local database. If you would like to start with some test data, request the data from the maintainers and then unzip it to `fixtures/db`. Note, we're working on creating synthetic test data for testing purposes and this will be committed to the repo soon.

Finally, to run the services:

```
$ docker compose -f ./containers/docker-compose.yaml --profile=all up
```

## CLI Tools

There are several CLI helper tools within the repository designed to make it easier to test out the RPCs and experiment locally. The code for these CLI tools is contained in the `cmd` folder. (*Note: CLI tools are not maintained to the same degree as the rest of the codebase; if you notice any odd behavior, don't hesitate to ask.*)

### GDS CLI

This program will allow you to interact with the TRISA global directory service (GDS) and some GDS Admin RPCs experimentally (e.g. looking up a VASP, reviewing a certificate request).

To install the GDS CLI, run this from the project root:

```bash
go install ./cmd/gds
```

### GDS Utils CLI

This program provides utilities for operating the GDS service and database.

To install the GDS Utils CLI, run this from the project root:

```bash
go install ./cmd/gdsutil
```

### Trtl CLI

This program provides tools for interacting with the Trtl data replication service.

To install the Trtl CLI, run this from the project root:

```bash
go install ./cmd/trtl
```


## Configuration

Server side configuration is done with the environment. Please see the instructions in the .env.template file for creating a local .env file to get started with development.

### Profiles

Client-side configuration is setup using profiles. A profile is a set of related configurations for both development and production. For example, the default environments are "production" to connect to vaspdirectory.net, "testnet" to connect to trisatest.net, and "localhost" to connect to locally running development servers. The profiles are configured in a YAML file that is stored in an OS-specific configuration directory.

You must first install the CLI Tools described in the section above. After that, install the profiles helper tool:


```
$ gds profiles --install
```

This will create YAML files in your OS-specific configuration directory. To view the path of the configuration file:

```
$ gds profiles --path
```

Basic usage is as follows:

1. `gds profiles` - show the configuration of the currently active profile
2. `gds profiles --list` - show a list of available profiles
3. `gds profiles --activate [name]` - activate the specified profile and use it
4. `gds profiles --edit` - edit the profiles using a command line editor

The easiest way to edit the profiles is to use `gds profiles --edit`, which will use the editor specified in the environment variable `$EDITOR`, or search for an editor in your `PATH` if none is specified. You must use a command line editor, e.g. `vim`, `emacs`, or `nano` so that the profile editor can verify the contents of the profiles before saving it.

If you would like to specify a different editor on the command line but do not want to set the `$EDITOR` environment variable, you can use a command-specific environment:

```
$ EDITOR=nano gds profiles --edit
```

If you would like to edit the profiles using VSCode or a GUI based editor, use the following command (note there will be no verification of correctness using this method):

```
$ code "$(gds profiles --path)"
```