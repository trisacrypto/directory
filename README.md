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

If you have issues you may also need to install these:

```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ go: downloading google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
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
$ docker compose -f ./containers/docker-compose.yaml build
```

Once the containers are built, you will need to supply a configuration `.env` file in the `./containers` directory (e.g. adjacent to the `docker-compose.yaml` file) with at least the following environment variables:

- `$SECTIGO_USERNAME`
- `$SECTIGO_PASSWORD`
- `$SENDGRID_API_KEY`
- `$GOOGLE_APPLICATION_CREDENTIALS`
- `$GOOGLE_PROJECT_NAME`

And then to run the services:

```
$ docker compose -f ./containers/docker-compose.yaml up
```

This will not run the UI container, just the backend containers.
