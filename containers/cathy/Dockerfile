# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.19-buster
ARG FINAL_IMAGE=debian:buster-slim

# Build Stage
FROM ${BUILDER_IMAGE} AS builder

ARG GIT_REVISION=""

# Ensure ca-certificates are up to date
RUN update-ca-certificates

# Use modeules for dependencies
WORKDIR $GOPATH/src/github.com/trisacrypto/directory
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy only what is needed (ignore web and fixtures)
COPY cmd ./cmd
COPY pkg ./pkg

# Build binaries
RUN go build -v -o /go/bin/sectigo -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/sectigo

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="TRISA Cathy CA for Integration Testing and Staging QA"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/sectigo /usr/local/bin/sectigo

CMD [ "/usr/local/bin/sectigo", "serve" ]
