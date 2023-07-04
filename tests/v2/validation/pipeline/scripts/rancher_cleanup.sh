#!/bin/bash
set -ex
cd $(dirname $0)/../../../../../

echo | corral config

corral list

echo "cleanup ranger"
tests/v2/validation/registries/bin/rangercleanup
