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