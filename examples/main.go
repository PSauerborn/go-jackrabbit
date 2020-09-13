package main

import (
	"fmt"
	"time"
	rabbit "github.com/PSauerborn/go-jackrabbit"
	log "github.com/sirupsen/logrus"
)


func ExampleHandler(msg []byte) {
	log.Info(fmt.Sprintf("received queue message: %s", string(msg)))
}

func main() {
	// create new rabbiMQ connection settings and send message over queue
	config := rabbit.RabbitConnectionConfig{QueueURL: "amqp://guest:guest@192.168.99.100:5672/", QueueName: "new-testing-queue"}
	go rabbit.ListenOnQueue(config, ExampleHandler)

	for {
		rabbit.ConnectAndDeliverOverQueue(config, []byte("testing queue message"))
		time.Sleep(5 * time.Second)
	}
}