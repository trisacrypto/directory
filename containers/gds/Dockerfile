FROM golang:1.16 AS builder

WORKDIR /srv/build

COPY . .
RUN go test ./... && go build -v ./cmd/gds
RUN go build -v ./cmd/sectigo
RUN go build -v ./cmd/debug

FROM ubuntu:bionic

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="Global TRISA Directory Service"

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get update && apt-get install -y wget gnupg

COPY --from=builder /srv/build/gds /bin/
COPY --from=builder /srv/build/sectigo /bin/
COPY --from=builder /srv/build/debug /bin/

ENTRYPOINT [ "/bin/gds", "serve" ]