AMQP to Kafka forwarder [![build status](https://travis-ci.org/MeteoGroup/amqp-to-kafka.svg)](https://travis-ci.org/MeteoGroup/amqp-to-kafka)
=======================

Fetches messages from AMQP and forwards them to Kafka.

 ## Build

Given a properly setup go develop environment, just run

```
  go get -d -t .
  go build
```

The resulting binary should be named `amqp-to-kafka` (or `amqp-to-kafka.exe`
on Windows). For further detail, e.g. cross-platform builds and custom build
parameters, please refer to the go-documentation.


## Usage

`amqp-to-kafka` takes the following commandline parameters:

  - `--amqp-uri`: URI of the AMQP broker
  - `--amqp-queue`: name of the queue to forward
  - `--amqp-exchange`: name of the AMQP exchange to forward messages from
  - `--amqp-binding-key`: binding key to select forwarded messages
  - `--amqp-consumer-tag`: consumer tag, see AMQP documentation
  - `--kafka-brokers=STRING`: list of Kafka brokers used for bootstrapping
  - `--kafka-topic=STRING`: Kafka topic for outgoing messages
  - `--metrics-address=HOST:PORT`: Listening address to serve metrics

`amqp-uri`, `amqp-exchange`, `kafka-brokers` and `kafka-topic` are mandatory,
everything else is optional. When `metrics-address` is given `amqp-to-kafka`
binds that address to export prometheus compatible metrics. Parameters may set
via environment variables as well

  - `AMQP_URI`: URI of the AMQP broker
  - `AMQP_QUEUE`: name of the queue to forward
  - `AMQP_EXCHANGE`: name of the AMQP exchange to forward messages from
  - `AMQP_BINDING_KEY`: binding key to select forwarded messages
  - `AMQP_CONSUMER_TAG`: consumer tag, see AMQP documentation
  - `KAFKA_BROKERS`: list of Kafka brokers used for bootstrapping
  - `KAFKA_TOPIC`: Kafka topic for outgoing messages
  - `METRICS_ADDRESS`: Listening address to serve metrics

When both are specified, commandline parameters take
precedence over environment variables.

Either `amqp-queue` or `amqp-exchange` and `amqp-binding-key` may be specified.
In case `amqp-queue` is used all messages to the given queue are forwarded.
In case `amqp-exchange` and `amqp-binding-key` are specified a temporary queue
is declared and bound to the exchange using the binding key to select messages
forwarded.


## Docker

A docker-ized variant is available `meteogroup/amqp-to-kafka`. Metrics are
exposed on port `8080`. To run use

```
docker run -P meteogroup/amqp-to-kafka <additional commandline arguments>
```


## License

Copyright © 2017 MeteoGroup Deutschland GmbH,
all the files in this repository are released under the terms of
[Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0).
