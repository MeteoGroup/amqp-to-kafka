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
  local i
  docker-compose -p "$TEST_STAGE" up -d kafka amqp

  wait-for kafka 9092
  wait-for amqp 5672

  in-test-stage nasqueron/rabbitmqadmin declare exchange name=public type=topic --host=amqp
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

wait-for() {
  local service="$1" port="$2"
  for ((i = 0; i < 10; ++i)); do
    in-test-stage toolbelt/netcat -z "$service" "$port" && break
  done
}

in-test-stage() {
  docker run --rm --net "${TEST_STAGE}_default" "$@"
}

main "$@"