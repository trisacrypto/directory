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
ARG REACT_APP_GOOGLE_CLIENT_ID
ARG REACT_APP_GDS_API_ENDPOINT="https://api.admin.vaspdirectory.net"
ARG REACT_APP_GDS_IS_TESTNET=false

# Build app with react-scripts
COPY web/gds-admin-ui/ ./
RUN yarn build

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <admin@trisa.io>"
LABEL description="GDS Administrative User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-admin-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
