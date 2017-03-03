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
  "runtime"
  "github.com/streadway/amqp"
  "github.com/Shopify/sarama"
  message "github.com/meteogroup/kafka-envelope/go/kafka_envelope"
  "time"
  "strconv"
)

func main() {
  loadConfig()

  amqpConsumer := openDeliveryChannel(amqpUri, amqpExchange, amqpQueue, amqpBindingKey, amqpConsumerTag)
  kafkaProducer := createKafkaProducer()
  startPrometheusHttpExporter()

  logInfo("lift off")
  for delivery := range amqpConsumer.deliveries {
    forwardToKafka(delivery, kafkaProducer)
    runtime.GC()
  }
}

func removeEmpty(m map[string]string) {
  for k, v := range m {
    if (v == "") {
      delete(m, k)
    }
  }
}

func forwardToKafka(delivery amqp.Delivery, kafkaProducer Producer) (forwarded []amqp.Delivery, skipped []amqp.Delivery) {
  receivedTime := time.Now()
  forwarded = []amqp.Delivery{}
  skipped = []amqp.Delivery{}
  deliveryMode := ""
  if (delivery.DeliveryMode == 1) {
    deliveryMode = "persistent"
  }
  if (delivery.DeliveryMode == 2) {
    deliveryMode = "non-persistent"
  }
  headers := map[string]string{
    "AppId": delivery.AppId,
    "ContentType": delivery.ContentType,
    "ContentEncoding": delivery.ContentEncoding,
    "ConsumerTag": delivery.ConsumerTag,
    "CorrelationId": delivery.CorrelationId,
    "DeliveryMode": deliveryMode,
    "Exchange": delivery.Exchange,
    "Expiration": delivery.Expiration,
    "MessageId": delivery.MessageId,
    "Priority": strconv.FormatUint(uint64(delivery.Priority), 10),
    "ReplyTo": delivery.ReplyTo,
    "RoutingKey": delivery.RoutingKey,
    "Type": delivery.Type,
    "UserId": delivery.UserId,
    "X-Received": receivedTime.UTC().Format(time.RFC3339),
    "X-Topic": kafkaTopic}
  removeEmpty(headers)
  if (delivery.DeliveryTag != 0) {
    headers["DeliveryTag"] = strconv.FormatUint(delivery.DeliveryTag, 16)
  }
  if (delivery.MessageCount > 0) {
    headers["MessageCount"] = strconv.FormatUint(uint64(delivery.MessageCount), 10)
  }
  if (delivery.Timestamp.UTC().Year() > 1980) {
    headers["Timestamp"] = delivery.Timestamp.UTC().Format(time.RFC3339)
  }
  partition, offset, err := kafkaProducer.publishMessage(message.ProducerEnvelope{
    Headers: headers,
    Payload: sarama.ByteEncoder(delivery.Body)})
  if (err != nil) {
    logError(err)
    skipped = append(skipped, delivery)
    messageCounter.WithLabelValues("skipped").Inc()
  } else {
    forwarded = append(forwarded, delivery)
    kafkaOffsets.WithLabelValues(strconv.FormatInt(int64(partition), 10)).Set(float64(offset))
    messageCounter.WithLabelValues("forwarded").Inc()
  }
  return
}
