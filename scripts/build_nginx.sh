#! /bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )"

source "${SCRIPT_DIR}/setup.sh"

DOCKER_FILE="${BUILD_DIR}/Dockerfile.nginx"

$DOCKER_CMD build \
    -t ${NGINX_IMAGE_NAME}:${NGINX_IMAGE_TAG} \
    -t ${NGINX_IMAGE_NAME}:latest \
    -f ${DOCKER_FILE} \
    ${BUILD_DIR}