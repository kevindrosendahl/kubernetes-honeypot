#! /bin/bash

set -eu

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

sudo mkdir -p /var/log/nginx

sudo docker run \
    -d \
    --rm \
    -v ${DIR}/nginx.conf:/etc/nginx/nginx.conf:ro \
    -v /var/log:/var/log \
    -v /home/eprokop/.minikube:/etc/certs \
    --net host \
    --name nginx \
    nginx
