FROM alpine:latest

LABEL maintainer="TRISA <info@trisa.io>"
LABEL description="StatefulSet initialization for trtl replicas"

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]