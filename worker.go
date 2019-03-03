// https://www.rabbitmq.com/tutorials/tutorial-two-go.html

// Fake a second of work for every dot in the message body. It will pop messages
// from the queue and perform the task
package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Setting up is the same as the publisher; we open a connection and a channel,
// and declare the queue from which we're going to consume. Note this matches up
// with the queue that `send.go` publishes to.

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Note that we declare the queue here, as well. Because we might start the
	// consumer before the publisher, we want to make sure the queue exists
	// before we try to consume messages from it.
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// We can set the prefetch count with the value of 1. This tells RabbitMQ
	// not to give more than one message to a worker at a time. Or, in other
	// words, don't dispatch a new message to a worker until it has processed
	// and acknowledged the previous one. Instead, it will dispatch it to the
	// next worker that is not still busy.
	//
	// NOTE: If all the workers are busy, your queue can fill up. You will want
	// to keep an eye on that, and maybe add more workers, or have some other
	// strategy.
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// We're about to tell the server to deliver us the messages from the queue.
	// Since it will push us messages asynchronously, we will read the messages
	// from a channel (returned by amqp::Consume) in a goroutine.
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			// Note that our fake task simulates execution time.
			time.Sleep(t * time.Second)
			log.Printf("Done")
			err := d.Ack(false)
			failOnError(err, "Failed to acknowledge a delivery")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
