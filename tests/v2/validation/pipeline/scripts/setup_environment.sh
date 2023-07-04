#!/bin/bash
set -ex
cd $(dirname $0)/../../../../../

echo "build corral packages"
sh tests/v2/validation/pipeline/scripts/build_corral_packages.sh

echo | corral config

echo "build rangerHA images"
sh tests/v2/validation/pipeline/scripts/build_rangerha_images.sh

corral list

echo "running ranger corral"
tests/v2/validation/registries/bin/rangerha
