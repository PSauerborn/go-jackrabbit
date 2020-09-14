package jackrabbit

import (
    "fmt"
    "github.com/streadway/amqp"
    log "github.com/sirupsen/logrus"
)

var (

)

// helper function used to check for valid queue config
func isValidQueueConfig(config RabbitConnectionConfig) bool {
    return config.QueueName != ""
}

// function used to generate new queue to deliver message
func ConnectAndDeliverOverQueue(config RabbitConnectionConfig, payload []byte) error {
    // connect to rabbitMQ server and generate new channel
    conn, err := NewRabbitConnection(config)
    if err != nil {
        return err
    }
    // defer closing of both channel and connection to rabbitMQ server
    defer conn.Connection.Close()
    defer conn.Channel.Close()

    // check that given configuration has necessary variables
    if !isValidQueueConfig(config) {
        log.Error("received invalid queue configuration")
        return InvalidRabbitMQConfig
    }
    // declare using rabbitMQ channel
    queue, err := conn.Channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
    if err != nil {
        log.Warn(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }
    // construct new message and send over rabbitMQ channel
    message := amqp.Publishing{Body: payload}
    err = conn.Channel.Publish("", queue.Name, false, false, message)
    if err != nil {
        log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
        return err
    }
    log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ queue", message))
    return nil
}

// function used to generate new queue to deliver message over a durable queue
func ConnectAndDeliverOverDurableQueue(config RabbitConnectionConfig, payload []byte) error {
    // connect to rabbitMQ server and generate new channel
    conn, err := NewRabbitConnection(config)
    if err != nil {
        return err
    }
    // defer closing of both channel and connection to rabbitMQ server
    defer conn.Connection.Close()
    defer conn.Channel.Close()

    // check that given configuration has necessary variables
    if !isValidQueueConfig(config) {
        log.Error("received invalid queue configuration")
        return InvalidRabbitMQConfig
    }
    // declare using rabbitMQ channel
    queue, err := conn.Channel.QueueDeclare(config.QueueName, true, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }
    // construct new message and send over rabbitMQ channel
    message := amqp.Publishing{Body: payload}
    err = conn.Channel.Publish("", queue.Name, false, false, message)
    if err != nil {
        log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
        return err
    }
    log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ queue", message))
    return nil
}

