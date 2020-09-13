# Introduction

High level utility library for Golang RabbitMQ connections. The library contains helper functions
used to send messages over rabbitMQ queues, as well as worker functions used to connect and listen
to rabbit queues

## Example Usage

All functions take a `RabbitConnectionConfig` parameter defining the connection settings, which is defined as

```go
type RabbitConnectionConfig struct {
	QueueURL     string `json:"queue_url" validate:"required"`
	QueueName    string `json:"queue_name"`
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
}
```

Once the connection setting is defined, the `ConnectAndDeliverOverQueue` and `ConnectAndDeliverOverExchange`
functions can be used to send messages over the defined queue. Additionally, the `ListenOnQueue` functions
can be used to listen on queues and handle any incoming messages

```go
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
```