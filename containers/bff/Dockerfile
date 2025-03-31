# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.23-bookworm
ARG FINAL_IMAGE=debian:bookworm-slim

# Build Stage
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} as builder

ARG TARGETOS
ARG TARGETARCH
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
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/bff -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/bff
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/bffutil -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/bffutil
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /go/bin/trtl -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$GIT_REVISION'" ./cmd/trtl

# Final Stage
FROM --platform=${BUILDPLATFORM} ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="TRISA GDS Backend for Front-End"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/bff /usr/local/bin/bff
COPY --from=builder /go/bin/bffutil /usr/local/bin/bffutil
COPY --from=builder /go/bin/trtl /usr/local/bin/trtl

CMD [ "/usr/local/bin/bff", "serve" ]
