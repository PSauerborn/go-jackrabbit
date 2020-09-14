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
    config := rabbit.RabbitConnectionConfig{
        QueueURL: "amqp://guest:guest@192.168.99.100:5672/",
        QueueName: "new-testing-queue",
        ExchangeName: "testing-exchange",
        ExchangeType: "fanout",
    }
    go rabbit.ListenOnQueueWithExchange(config, ExampleHandler)

    for {
        err := rabbit.ConnectAndDeliverOverQueue(config, []byte("testing queue message"))
        if err != nil {
            log.Error(fmt.Errorf("unable to send message over Rabbit Queue: %v", err))
        }
        time.Sleep(5 * time.Second)
    }
}