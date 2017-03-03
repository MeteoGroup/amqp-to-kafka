package test

import (
  "github.com/streadway/amqp"
  "github.com/Shopify/sarama"
  envelope "github.com/meteogroup/kafka-envelope/go/kafka_envelope"
  "testing"
  "strings"
  "time"
)

func TestDeliveryToKafka(t *testing.T) {
  receive, err := createKafkaMessageChannel("kafka:9092", "amqp-messages")
  if (err != nil) {
    t.Fatal(err)
  }
  err = publishMessageToAMQP("amqp://amqp", "public", "/messages", "a message")
  if (err != nil) {
    t.Fatal(err)
  }

  timeout := make(chan interface{});
  go func() {
    time.Sleep(10 * time.Second)
    timeout <- nil
  }()
  select {
  case message := <-receive:
    if (string(message.Payload) != "a message") {
      t.Fatal("Unexpected message, expected: \"a message\", got: \"" + string(message.Payload) + "\"")
    }
  case <-timeout:
    t.Fatal("Time out")
  }
}

func publishMessageToAMQP(amqpURI, exchange, routingKey, message string) (err error) {
  connection, err := amqp.Dial(amqpURI)
  if (err != nil) {
    return
  }

  channel, err := connection.Channel()
  if (err != nil) {
    return
  }

  err = channel.Publish(
    exchange, // publish to an exchange
    routingKey, // routing to 0 or more queues
    false, // mandatory
    false, // immediate
    amqp.Publishing{
      Headers:         amqp.Table{},
      ContentType:     "text/plain",
      ContentEncoding: "utf-8",
      Body:            []byte(message),
      DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
      Priority:        0, // 0-9
      // a bunch of application/implementation-specific fields
    },
  );
  return
}

func createKafkaMessageChannel(kafkaBrokers, topic string) (<-chan envelope.Envelope, error) {
  config := sarama.NewConfig()
  config.Producer.RequiredAcks = sarama.WaitForAll
  config.Producer.Retry.Max = 5
  config.Producer.Return.Successes = true

  brokers := strings.Split(kafkaBrokers, ",")
  consumer, err := sarama.NewConsumer(brokers, config)
  if (err != nil) {
    return nil, err
  }
  partitionConsumer, err := consumer.ConsumePartition(topic, 0, 0)
  if (err != nil) {
    return nil, err
  }
  messages := make(chan envelope.Envelope)
  rawMessages := partitionConsumer.Messages()
  go func() {
    for in := range rawMessages {
      out, err := envelope.ConsumerEnvelope(in.Value).Unmarshal()
      if (err != nil) {
        panic(err)
      }
      messages <- out
    }
  }()
  return messages, err
}
