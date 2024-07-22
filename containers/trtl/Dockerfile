# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.22-bookworm
ARG FINAL_IMAGE=debian:bookworm-slim

# Build stage1
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} as builder

ARG TARGETOS
ARG TARGETARCH
ARG GIT_REVISION=""

# Ensure ca-certificates are up to date
RUN update-ca-certificates

WORKDIR $GOPATH/src/github.com/trisacrypto/directory

# Use modules for dependencies
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy only what is needed (ignore web and fixtures)
COPY cmd ./cmd
COPY pkg ./pkg
COPY proto ./proto

# Build binaries
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/trtl -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/trtl
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/gdsutil -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/gdsutil

# Final Stage
FROM --platform=${BUILDPLATFORM} ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="GDS Replicated Data Backend (TRTL)"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/trtl /usr/local/bin/trtl
COPY --from=builder /go/bin/gdsutil /usr/local/bin/gdsutil

CMD [ "/usr/local/bin/trtl", "serve" ]