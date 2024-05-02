# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.19-bullseye
ARG FINAL_IMAGE=debian:bullseye-slim

# Build stage
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
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/gds -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/gds
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/gdsutil -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/gdsutil
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/trtl -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/trtl

# Final Stage
FROM --platform=${BUILDPLATFORM} ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="Global TRISA Directory Service"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/gds /usr/local/bin/gds
COPY --from=builder /go/bin/gdsutil /usr/local/bin/gdsutil
COPY --from=builder /go/bin/trtl /usr/local/bin/trtl

CMD [ "/usr/local/bin/gds", "serve" ]