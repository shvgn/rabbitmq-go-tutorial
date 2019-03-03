// https://www.rabbitmq.com/tutorials/tutorial-two-go.html

package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// The connection abstracts the socket connection, and takes care of
	// protocol version negotiation and authentication and so on for us.
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// To send, we must declare a queue for us to send to.
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Now we can publish a message to the queue.  Declaring a queue is
	// idempotent - it will only be created if it doesn't exist already. The
	// message content is a byte array, so you can encode whatever you like
	// there.

	// We'll be sending strings that stand for complex tasks. We don't have a
	// real-world task, like images to be resized or pdf files to be rendered,
	// so let's fake it by just pretending we're busy - by using the
	// `time.Sleep` function. We'll take the number of dots in the string as its
	// complexity; every dot will account for one second of "work". For example,
	// a fake task described by `Hello...` will take three seconds.
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [×] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello " + (time.Now()).String()
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
