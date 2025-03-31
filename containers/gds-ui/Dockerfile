# Dynamic Builds
ARG BUILDER_IMAGE=node:16
ARG FINAL_IMAGE=nginx:stable

# Build stage
FROM ${BUILDER_IMAGE} as builder

WORKDIR /usr/src/app

# Dependencies
COPY web/gds-ui/package.json ./
COPY web/gds-ui/yarn.lock ./
RUN yarn install --frozen-lockfile --check-files

# Set production environment variable for build context
ARG REACT_APP_GDS_API_ENDPOINT="https://proxy.trisa.directory"
ARG REACT_APP_GDS_IS_TESTNET=false
ARG REACT_APP_ANALYTICS_ID=""

# Build app with react-scripts
COPY web/gds-ui/ ./
RUN yarn build

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <admin@trisa.io>"
LABEL description="Global TRISA Directory Service User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
