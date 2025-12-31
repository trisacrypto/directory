ARG BASE_IMAGE=nginx:stable
FROM ${BASE_IMAGE} AS final

LABEL maintainer="TRISA <admin@travelrule.io>"
LABEL description="TRISA GDS Maintenance Mode Placeholder"

COPY web/maintenance /usr/share/nginx/html
COPY containers/maintenance/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
