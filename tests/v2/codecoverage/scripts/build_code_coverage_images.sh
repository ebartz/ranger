#!/bin/bash
set -e
cd $(dirname $0)/../../../../


echo "building cluster setup bin"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tests/v2/codecoverage/bin/setupranger ./tests/v2/codecoverage/setupranger

echo "building ranger HA corral bin"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tests/v2/codecoverage/bin/rangercorral ./tests/v2/codecoverage/rangerha

echo "building ranger cleanup bin"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tests/v2/codecoverage/bin/rangercleanup ./tests/v2/codecoverage/rangercleanup
