#!/bin/bash
set -e
cd $(dirname $0)/../../../../../

echo "building ranger HA corral bin"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tests/v2/validation/registries/bin/rangerha ./tests/v2/validation/pipeline/rangerha

echo "building ranger cleanup bin"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tests/v2/validation/registries/bin/rangercleanup ./tests/v2/validation/pipeline/rangercleanup
