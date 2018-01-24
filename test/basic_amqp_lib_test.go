package test

import (
  "github.com/streadway/amqp"
  "testing"
)

func TestRedefineQueue(t *testing.T) {
  connection, err := amqp.Dial("amqp://amqp")
  if (err != nil) {
    t.Fatal(err)
  }

  channel, err := connection.Channel()
  if (err != nil) {
    t.Fatal(err)
  }

  queue, err := channel.QueueDeclare(
    "q_anonymous_test", // name of the queue
    false, // durable
    true, // delete when unused
    false, // exclusive
    false, // noWait
    nil, // arguments
  )
  if (err != nil) {
    t.Fatal(err)
  }

  if (queue.Name != "q_anonymous_test") {
    t.Error("Expected test queue name, but got ", queue.Name)
  }

  // declare queue a second time to learn error object if any
  queueNew, err := channel.QueueDeclare("q_anonymous_test", false, true, false, false, nil,)
  if (err != nil) {
    t.Fatal(err)
  }
  if (queue != queueNew) {
    t.Error("First Queue ", queue, ", Second queue ", queueNew)
  }
}

func TestDoubleBindQueue(t *testing.T) {
  connection, err := amqp.Dial("amqp://amqp")
  if (err != nil) {
    t.Fatal(err)
  }

  channel, err := connection.Channel()
  if (err != nil) {
    t.Fatal(err)
  }
  queue, err := channel.QueueDeclare("q_anonymous_test", false, true, false, false, nil,)
  if (err != nil) {
    t.Fatal(err)
  }

  err = channel.QueueBind(
    queue.Name, // name of the queue
    "#", // bindingKey
    "public", // sourceExchange
    false, // noWait
    nil, // arguments
  )
  if (err != nil) {
    t.Fatal(err)
  }
  // bind a second time
  err = channel.QueueBind(queue.Name, "#", "public", false, nil,)
  if (err != nil) {
    t.Fatal(err)
  }

}
