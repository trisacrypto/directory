FROM node:16 AS builder

WORKDIR /usr/src/app

# Dependencies
COPY web/gds-admin-ui/package*.json ./
RUN npm install

# Set production environment variable for build context
ARG REACT_APP_GDS_API_ENDPOINT="https://api.admin.vaspdirectory.net"
ARG REACT_APP_GDS_IS_TESTNET=false

# Build app with browserify
COPY web/gds-admin-ui/ ./
RUN npm run build

FROM nginx:stable

LABEL maintainer="TRISA <admin@trisa.io>"
LABEL description="GDS Administrative User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-admin-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
