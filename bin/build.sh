#!/usr/bin/env bash
set -e

if [ $# -eq 0 ]
  then
    tag="latest"
  else
    tag=$1
fi

echo "Using tag $tag"

dingo -src=./internal/config/di -dest=./generated

docker build . -t navexplorer/indexer:$tag
docker push navexplorer/indexer:$tag
