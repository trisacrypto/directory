FROM alpine:latest

LABEL maintainer="TRISA <info@travelrule.io>"
LABEL description="StatefulSet initialization for trtl replicas"

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
