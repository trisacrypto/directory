# Dynamic Builds
ARG BUILDER_IMAGE=node:18
ARG FINAL_IMAGE=nginx:stable

# Build stage
FROM ${BUILDER_IMAGE} as builder

WORKDIR /usr/src/app

# Dependencies
COPY web/gds-user-ui/package.json ./
COPY web/gds-user-ui/yarn.lock ./
RUN yarn install --frozen-lockfile --check-files

# Set production environment variable for build context
ARG REACT_APP_TRISA_BASE_URL="https://bff.trisa.directory/v1/"
ARG REACT_APP_ANALYTICS_ID
ARG REACT_APP_VERSION_NUMBER
ARG REACT_APP_GIT_REVISION
ARG REACT_APP_AUTH0_DOMAIN
ARG REACT_APP_AUTH0_CLIENT_ID
ARG REACT_APP_AUTH0_SCOPE="openid profile email"
ARG REACT_APP_AUTH0_REDIRECT_URI="https://trisa.directory/auth/callback"
ARG REACT_APP_AUTH0_AUDIENCE="https://bff.trisa.directory"
ARG REACT_APP_SENTRY_DSN
ARG REACT_APP_SENTRY_ENVIRONMENT
ARG REACT_APP_USE_DASH_LOCALE

# Build app with react-scripts
COPY web/gds-user-ui/ ./
RUN yarn build

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="TRISA <admin@travelrule.io>"
LABEL description="Global TRISA Directory Service Advanced User Interface"

COPY --from=builder /usr/src/app/build/ /usr/share/nginx/html
COPY containers/gds-user-ui/nginx.conf /etc/nginx/conf.d/default.conf

ENV NGINX_ENTRYPOINT_QUIET_LOGS=1
CMD [ "nginx", "-g", "daemon off;" ]
