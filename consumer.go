package jackrabbit

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

// function used to start blocking goroutine that connects to rabbitMQ
// connection and executes some handler function when messages are send
// over the specified queue
func ListenOnQueue(config RabbitConnectionConfig, handler func(payload []byte)) error {
    // check that given configuration has necessary variables
    if !isValidQueueConfig(config) {
        log.Error("received invalid exchange configuration")
        return InvalidRabbitMQConfig
    }
    // connect to rabbitMQ server and generate new channel
    conn, err := NewRabbitConnection(config)
    if err != nil {
        log.Fatal(fmt.Errorf("unable to connect to Rabbit Broker: %v", err))
        return err
    }
    // defer closing of both channel and connection to rabbitMQ server
    defer conn.Connection.Close()
    defer conn.Channel.Close()

    // declare queue
    queue, err := conn.Channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }

    // consume to messages on channel
    messages, err := conn.Channel.Consume(queue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }
    // start goroutine to handle messages
    forever := make(chan bool)
    go func() { for d := range messages { handler(d.Body) }}()
    <-forever
    return nil
}

// function used to start blocking goroutine that connects to rabbitMQ
// connection and executes some handler function when messages are send
// over the specified queue or exchange
func ListenOnQueueWithExchange(config RabbitConnectionConfig, handler func(payload []byte)) error {
    // check that given configuration has necessary variables
    if !isValidExchangeConfig(config) || !isValidQueueConfig(config) {
        log.Error("received invalid exchange configuration")
        return InvalidRabbitMQConfig
    }
    // connect to rabbitMQ server and generate new channel
    conn, err := NewRabbitConnection(config)
    if err != nil {
        log.Fatal(fmt.Errorf("unable to connect to Rabbit Broker: %v", err))
        return err
    }
    // defer closing of both channel and connection to rabbitMQ server
    defer conn.Connection.Close()
    defer conn.Channel.Close()

    // declare events exchange with fanout type
    err = conn.Channel.ExchangeDeclare(config.ExchangeName, config.ExchangeType, false, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbitmq exchange: %s", err))
        return err
    }
    // declare queue
    queue, err := conn.Channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }
    // bind queue to exchange
    err = conn.Channel.QueueBind(queue.Name, "", config.ExchangeName, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to bind to exchange: %v", err))
        return err
    }
    // consume to messages on channel
    messages, err := conn.Channel.Consume(queue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
        return err
    }
    // start goroutine to handle messages
    forever := make(chan bool)
    go func() { for d := range messages { handler(d.Body) }}()
    <-forever

    return nil
}