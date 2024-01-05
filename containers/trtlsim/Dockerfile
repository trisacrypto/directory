# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.21-bookworm
ARG FINAL_IMAGE=debian:bookworm-slim

# Build stage
FROM ${BUILDER_IMAGE} as builder

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
RUN go build -v -o /go/bin/trtlsim ./cmd/trtlsim

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="GDS Replicated Data Backend (TRTL)"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/trtlsim /usr/local/bin/trtlsim

CMD [ "/usr/local/bin/trtlsim" ]