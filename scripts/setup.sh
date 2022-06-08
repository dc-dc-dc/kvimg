DOCKER_CMD="docker"
TOP_LEVEL=$(dirname ${SCRIPT_DIR})

BUILD_DIR="${TOP_LEVEL}/build"
DATA_DIR="${TOP_LEVEL}/data"

NET_NAME="dc_dev_net"
NGINX_IMAGE_NAME="dc-dc-dc/nginx"
NGINX_IMAGE_TAG="0.1.0"
KVIMG_IMAGE_NAME="dc-dc-dc/kvimg"
KVIMG_IMAGE_TAG="0.1.0"