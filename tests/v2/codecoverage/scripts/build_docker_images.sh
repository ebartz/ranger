#!/bin/bash
set -e

ARCH=${ARCH:-"amd64"}
REPO=rangertest
TAG=v2.7-head
SYSTEM_CHART_REPO_DIR=build/system-charts
SYSTEM_CHART_DEFAULT_BRANCH=${SYSTEM_CHART_DEFAULT_BRANCH:-"dev-v2.7"}
CHART_REPO_DIR=build/charts
CHART_DEFAULT_BRANCH=${CHART_DEFAULT_BRANCH:-"dev-v2.7"}

cd $(dirname $0)/../package

../scripts/k3s-images.sh

cp ../bin/ranger ../bin/agent ../bin/data.json ../bin/k3s-airgap-images.tar .

# Make sure the used data.json is a release artifact
cp ../bin/data.json ../bin/ranger-data.json

IMAGE=${REPO}/ranger:${TAG}
AGENT_IMAGE=${REPO}/ranger-agent:${TAG}

echo "building ranger test docker image"
docker build --build-arg VERSION=${TAG} --build-arg ARCH=${ARCH} --build-arg IMAGE_REPO=${REPO} -t ${IMAGE} -f Dockerfile . --no-cache

echo "building agent test docker image"
docker build --build-arg VERSION=${TAG} --build-arg ARCH=${ARCH} --build-arg RANCHER_TAG=${TAG} --build-arg RANCHER_REPO=${REPO} -t ${AGENT_IMAGE} -f Dockerfile.agent . --no-cache

echo ${DOCKERHUB_PASSWORD} | docker login --username ${DOCKERHUB_USERNAME} --password-stdin

echo "docker push ranger"
docker image push ${IMAGE}

echo "docker push agent"
docker image push ${AGENT_IMAGE}
