/*
Copyright © 2017 MeteoGroup Deutschland GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */
package main

import (
  "flag"
  "os"
)

var (
  amqpUri = ""
  amqpExchange = ""
  amqpQueue = ""
  amqpBindingKey = ""
  amqpConsumerTag = ""
  kafkaBrokers = ""
  kafkaTopic = ""
  metricsAddress = ""
)

func loadConfig() {
  flag.StringVar(&amqpUri, "amqp-uri", os.Getenv("AMQP_URI"), "AMQP URI")
  flag.StringVar(&amqpExchange, "amqp-exchange", os.Getenv("AMQP_EXCHANGE"), "Durable, non-auto-deleted AMQP exchange name")
  flag.StringVar(&amqpQueue, "amqp-queue", os.Getenv("AMQP_QUEUE"), "Ephemeral AMQP queue name")
  flag.StringVar(&amqpBindingKey, "amqp-binding-key", os.Getenv("AMQP_BINDING_KEY"), "AMQP binding key")
  flag.StringVar(&amqpConsumerTag, "amqp-consumer-tag", os.Getenv("AMQP_CONSUMER_TAG"), "AMQP consumer tag (should not be blank)")
  flag.StringVar(&kafkaBrokers, "kafka-brokers", os.Getenv("KAFKA_BROKERS"), "list of Kafka brokers used for bootstrapping")
  flag.StringVar(&kafkaTopic, "kafka-topic", os.Getenv("KAFKA_TOPIC"), "Kafka topic for outgoing messages")
  flag.StringVar(&metricsAddress, "metrics-address", os.Getenv("METRICS_ADDRESS"), "Listening address to serve metrics")
  flag.Parse()

  if amqpUri == "" {
    panic("Required parameter `amqp-uri` is missing or empty.")
  }
  if amqpExchange == "" {
    panic("Required parameter `amqp-exchange` is missing or empty.")
  }
  if kafkaBrokers == "" {
    panic("Required parameter `kafka-brokers` is missing or empty.")
  }
  if kafkaTopic == "" {
    panic("Required parameter `kafka-topic` is missing or empty.")
  }
  if amqpQueue != "" && (amqpBindingKey != "" || amqpExchange != "") {
    panic("AMQP queue cannot be specified together with AMQP binding key or AMQP exchange")
  }

  logConfig()
}

func logConfig() {
  parameters := []interface{}{}
  appendIfDefined := func(name string, value string) {
    if (value != "") {
      parameters = append(parameters, name, value)
    }
  }
  appendIfDefined("amqpUri", amqpUri)
  appendIfDefined("amqpExchange", amqpExchange)
  appendIfDefined("amqpQueue", amqpQueue)
  appendIfDefined("amqpBindingKey", amqpBindingKey)
  appendIfDefined("amqpConsumerTag", amqpConsumerTag)
  appendIfDefined("kafkaBrokers", kafkaBrokers)
  appendIfDefined("kafkaTopic", kafkaTopic)
  appendIfDefined("metricsAddress", metricsAddress)
  logInfo("starting with configuration", parameters...)
}
