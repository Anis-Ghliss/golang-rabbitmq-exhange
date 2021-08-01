package main

import (
	"log"

	"github.com/streadway/amqp"
)

func exit_on_error(err error) {
        if err != nil {
                log.Fatal(err)
        }
}

func main() {
        conn, err := amqp.Dial("url")
        exit_on_error(err)
        defer conn.Close()

        ch, err := conn.Channel()
        exit_on_error(err)
        defer ch.Close()

        err = ch.ExchangeDeclare(
                "example.topic",   // name
                "topic", // type
                true,     // durable
                false,    // auto-deleted
                false,    // internal
                false,    // no-wait
                nil,      // arguments
        )
        exit_on_error(err)

        q, err := ch.QueueDeclare(
                "isporit-videos-feedbacks",    // name
                false, // durable
                false, // delete when usused
                true,  // exclusive
                false, // no-wait
                nil,   // arguments
        )
        exit_on_error(err)

        routingkey := "eu.de.*"
        err = ch.QueueBind(
                q.Name, // queue name
                routingkey,     // routing key
                "example.topic", // exchange
                false,
                nil,
        )
        exit_on_error(err)

        msgs, err := ch.Consume(
                q.Name, // queue
                "",     // consumer
                true,   // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
        )
        exit_on_error(err)

        forever := make(chan bool)

        go func() {
                for d := range msgs {
                        log.Printf(" [x] %s", d.Body)
                }
        }()

        log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
        <-forever
}