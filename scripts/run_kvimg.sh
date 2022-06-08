#! /bin/bash
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )"

source "${SCRIPT_DIR}/setup.sh"

"${SCRIPT_DIR}/network.sh"

NAME="dc_dev_kvimg"

$DOCKER_CMD create \
    --rm \
    --name ${NAME} \
    -v ${TOP_LEVEL}/db:/db \
    -p 4000:3000 \
    ${KVIMG_IMAGE_NAME}:${KVIMG_IMAGE_TAG} \
    /main -port=3000 -action=server -servers=http://nginx_vol_vol1:80

$DOCKER_CMD network connect \
    ${NET_NAME} \
    ${NAME}

$DOCKER_CMD start ${NAME}