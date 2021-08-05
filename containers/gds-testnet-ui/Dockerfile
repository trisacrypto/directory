FROM node:16 AS builder

WORKDIR /usr/src/app

# Dependencies
COPY web/package*.json ./
RUN npm install

# Set production environment variable for build context
ENV REACT_APP_GDS_API_ENDPOINT="https://proxy.trisatest.net"
ENV REACT_APP_GDS_IS_TESTNET=true

# Build app with browserify
COPY web/ ./
RUN npm run build

FROM nginx:stable

LABEL maintainer="TRISA <admin@trisa.io>"
LABEL description="Global TRISA Directory Service User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-testnet-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
