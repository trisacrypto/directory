# Dynamic Builds
ARG BUILDER_IMAGE=node:16
ARG FINAL_IMAGE=nginx:stable

# Build stage
FROM ${BUILDER_IMAGE} as builder

WORKDIR /usr/src/app

# Dependencies
COPY web/gds-admin-ui/package.json ./
COPY web/gds-admin-ui/yarn.lock ./
RUN yarn install --frozen-lockfile --check-files

# Set production environment variable for build context
ARG REACT_APP_VERSION_NUMBER
ARG REACT_APP_GIT_REVISION
ARG REACT_APP_GDS_API_ENDPOINT="https://api.admin.trisa.directory"
ARG REACT_APP_GDS_IS_TESTNET=false
ARG REACT_APP_GOOGLE_CLIENT_ID
ARG REACT_APP_SENTRY_DSN
ARG REACT_APP_SENTRY_ENVIRONMENT

# Build app with react-scripts
COPY web/gds-admin-ui/ ./
RUN yarn build

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <admin@travelrule.io>"
LABEL description="GDS Administrative User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-admin-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
