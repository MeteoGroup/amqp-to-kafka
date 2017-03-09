#!/bin/sh -e

if [ "$TRAVIS" != true ]; then
  echo "This script is intended to run within the travis build only" 1>&2
  exit 1
fi

( # temporarily disable shell debugging, docker hub password would leak otherwise
  set +x
  docker login -u "$DOCKER_USER" -p "$DOCKER_PASS"
)

set -x
IMAGE_TAG="${TRAVIS_TAG:-"`date '+%Y%m%dT%H%M%S'`"}"
docker tag "meteogroup/amqp-to-kafka:$COMMIT" "meteogroup/amqp-to-kafka:$IMAGE_TAG"
docker tag "meteogroup/amqp-to-kafka:$COMMIT" "meteogroup/amqp-to-kafka:latest"
docker push "meteogroup/amqp-to-kafka:$IMAGE_TAG"
docker push "meteogroup/amqp-to-kafka:latest"
