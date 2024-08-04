#!/bin/bash
CUR_DIR=$(cd $(dirname $0); pwd)

export SERVICE_NAME=simple_short_url
export PORT=8080
export ACCESS_TOKEN="short_url_custom_access_token"

echo "CUR_DIR/bin/${SERVICE_NAME}"
exec ${CUR_DIR}/bin/${SERVICE_NAME}