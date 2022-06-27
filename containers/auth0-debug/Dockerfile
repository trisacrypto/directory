# Dynamic Builds
ARG PYTHON_IMAGE=python:3.10.5-slim-buster

FROM ${PYTHON_IMAGE} as final

WORKDIR /app

COPY web/auth0-debug .

CMD [ "python3", "serve.py" ]