#!/bin/bash
set -e

cd $(dirname $0)

echo "build ranger"
./build-ranger.sh
echo "build agent"
./build-agent.sh