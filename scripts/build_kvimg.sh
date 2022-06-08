#! /bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )"

source "${SCRIPT_DIR}/setup.sh"

DOCKER_FILE="${BUILD_DIR}/Dockerfile.kvimg"

$DOCKER_CMD build \
    -t ${KVIMG_IMAGE_NAME}:${KVIMG_IMAGE_TAG} \
    -t ${KVIMG_IMAGE_NAME}:latest \
    -f ${DOCKER_FILE} \
    ${TOP_LEVEL}