package main

import (
	"log"
	"os"

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

        routingkey := os.Args[1]
        message := os.Args[2]

        err = ch.Publish(
                "example.topic", // exchange
                routingkey,     // routing key
                false,  // mandatory
                false,  // immediate
                amqp.Publishing {
                        ContentType: "text/plain",
                        Body:        []byte(message),
                },
        )
        exit_on_error(err)

        log.Printf(" [x] Sent %s", message)
}