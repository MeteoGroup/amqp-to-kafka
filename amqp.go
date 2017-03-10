package main

import (
  "fmt"
  "github.com/streadway/amqp"
)

type Consumer struct {
  connection *amqp.Connection
  channel    *amqp.Channel
  deliveries <-chan amqp.Delivery
  tag        string
}

func openDeliveryChannel(amqpURI, exchange, queueName, bindingKey, consumerTag string) (consumer Consumer) {
  consumer = Consumer{connection: nil, channel: nil, deliveries: nil, tag: consumerTag}
  var err error

  consumer.connection, err = amqp.Dial(amqpURI)
  logAndPanic(err)

  consumer.channel, err = consumer.connection.Channel()
  logAndPanic(err)

  if (queueName == "") {
    if (bindingKey == "") {
      bindingKey = "#"
    }
    queue, err := consumer.channel.QueueDeclare(
      "", // name of the queue
      false, // durable
      true, // delete when unused
      false, // exclusive
      false, // noWait
      nil, // arguments
    )
    logAndPanic(err)
    queueName = queue.Name

    err = consumer.channel.QueueBind(
      queueName, // name of the queue
      bindingKey, // bindingKey
      exchange, // sourceExchange
      false, // noWait
      nil, // arguments
    )
    logAndPanic(err)
  }

  deliveries, err := consumer.channel.Consume(
    queueName, // name
    consumer.tag, // consumerTag,
    false, // noAck
    false, // exclusive
    false, // noLocal
    false, // noWait
    nil, // arguments
  )
  logAndPanic(err)

  consumer.deliveries = deliveries
  return
}

func (c Consumer) shutdown() error {
  // will close() the deliveries channel
  if err := c.channel.Cancel(c.tag, true); err != nil {
    return fmt.Errorf("Consumer cancel failed: %s", err)
  }

  if err := c.connection.Close(); err != nil {
    return fmt.Errorf("AMQP connection close error: %s", err)
  }

  return nil
}
