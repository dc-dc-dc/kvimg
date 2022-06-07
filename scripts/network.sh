#! /bin/bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )"

source "${SCRIPT_DIR}/setup.sh"

NET_CHECK=$( docker network inspect ${NET_NAME} 2>&1 1>/dev/null)

if [ "$NET_CHECK" != "" ]; then
    echo "No network found, creating...";

    ${DOCKER_CMD} network create \
        -d bridge \
        ${NET_NAME}
fi
