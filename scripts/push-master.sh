#!/bin/bash
set -e
cd $(dirname $0)/..
mkdir -p dist
./scripts/build
export REPO=${REPO:-ranger}
TAG=master ./scripts/package
docker push $REPO/server:master
