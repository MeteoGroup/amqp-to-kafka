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
  "strings"
  "github.com/Shopify/sarama"
)

type Producer struct {
  p sarama.SyncProducer
}

func createKafkaProducer() Producer {
  config := sarama.NewConfig()
  config.Producer.RequiredAcks = sarama.WaitForAll
  config.Producer.Retry.Max = 5
  config.Producer.Return.Successes = true

  brokers := strings.Split(kafkaBrokers, ",")
  producer, err := sarama.NewSyncProducer(brokers, config)
  logAndPanic(err)
  return Producer{p: producer}
}

func (producer Producer)publishMessage(message sarama.Encoder) (partition int32, offset int64, err error) {
  partition, offset, err = producer.p.SendMessage(&sarama.ProducerMessage{
    Topic: kafkaTopic,
    Value: message,
  })
  return
}
