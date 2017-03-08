#!/bin/bash -ex

cd "$(dirname "$0")"

TEST_STAGE_ID="$(printf '%04x'  "$RANDOM" "$RANDOM")"
TEST_STAGE="amqp2kafkatest$TEST_STAGE_ID"
export IMAGE_TAG="${1:-test}"

main() {
  trap 'tear-down' EXIT
  setup

  docker-compose -p "$TEST_STAGE" run test
}

setup() {
  docker-compose -p "$TEST_STAGE" up -d amqp-to-kafka
}

tear-down() {
  result="$?"
  docker-compose -p "$TEST_STAGE" stop
  if (( result )); then
    docker-compose -p "$TEST_STAGE" logs
  fi
  docker-compose -p "$TEST_STAGE" down --rmi local --volumes --remove-orphans
  return "$result"
}

main "$@"
