# Global TRISA VASP Directory

**TRISA implementation of the Global VASP Directory Service**


## Docker Compose

The simplest way to develop the GDS UI is to run the GDS service and Envoy gRPC-Proxy using `docker compose`. The containers directory has the configuration and Dockerfiles to build the images:

```
$ docker compose -f ./containers/docker-compose.yaml build
```

Once the containers are built, you will need to supply a configuration `.env` file in the `./containers` directory (e.g. adjacent to the `docker-compose.yaml` file) with at least the following environment variables:

- `$SECTIGO_USERNAME`
- `$SECTIGO_PASSWORD`
- `$SENDGRID_API_KEY`

And then to run the services:

```
$ docker compose -f ./containers/docker-compose.yaml up
```

This will not run the UI container, just the backend containers.