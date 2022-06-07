#! /bin/bash
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )"

source "${SCRIPT_DIR}/setup.sh"

"${SCRIPT_DIR}/network.sh"

PORT=3000
VOLUME=${1:-vol1}
C_PATH="${DATA_DIR}/${VOLUME}"
NAME="nginx_vol_${VOLUME}"
if [ ! -d ${C_PATH} ]; then
    echo "no dir creating"
    chmod 777 ${C_PATH}
    mkdir -p ${C_PATH}
fi

$DOCKER_CMD create \
    --rm \
    --name ${NAME} \
    -v ${C_PATH}:/data \
    -p ${PORT}:80 \
    ${NGINX_IMAGE_NAME}:${NGINX_IMAGE_TAG}

$DOCKER_CMD network connect \
    ${NET_NAME} \
    ${NAME}

$DOCKER_CMD start ${NAME}