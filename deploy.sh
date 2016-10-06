#!/bin/bash

_tag=$1

if [ -z "${_tag}" ]; then
    source _VERSION
    _tag=${_VERSION}
fi

docker tag "orangesys/alpine-slack8s:${tag}" "asia.gcr.io/orange-sys/alpine-slack8s:${tag}"
docker login -e $DOCKER_EMAIL -u _json_key -p "$(cat ${HOME}/account-auth.json)" https://asia.gcr.io
docker push asia.gcr.io/orange-sys/alpine-slack8s:0.1.3
