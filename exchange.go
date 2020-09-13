package jackrabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
)

// helper function used to check for exchange config
func isValidExchangeConfig(config RabbitConnectionConfig) bool {
	return config.ExchangeName != "" && config.ExchangeType != ""
}

// function used to generate new exchange to deliver message over exchange
func ConnectAndDeliverOverExchange(config RabbitConnectionConfig, payload []byte) error {
	// connect to rabbitMQ server and generate new channel
	conn, err := NewRabbitConnection(config)
	if err != nil {
		return err
	}
	// defer closing of both channel and connection to rabbitMQ server
	defer conn.Connection.Close()
	defer conn.Channel.Close()

	// check that given configuration has neccessary variables
	if !isValidExchangeConfig(config) {
		log.Error("received invalid exchange configuration")
		return InvalidRabbitMQConfig
	}

	// declare events exchange with fanout type
	err = conn.Channel.ExchangeDeclare(config.ExchangeName, config.ExchangeType, false, false, false, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq exchange: %s", err))
		return err
	}
	// construct new message and send over rabbitMQ channel
	message := amqp.Publishing{Body: payload}
	err = conn.Channel.Publish(config.ExchangeName, "", false, false, message)
	if err != nil {
		log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
		return err
	}
	log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ exchange", message))
	return nil
}

// function used to generate new exchange to deliver message over exchange
func ConnectAndDeliverOverDurableExchange(config RabbitConnectionConfig, payload []byte) error {
	// connect to rabbitMQ server and generate new channel
	conn, err := NewRabbitConnection(config)
	if err != nil {
		return err
	}
	// defer closing of both channel and connection to rabbitMQ server
	defer conn.Connection.Close()
	defer conn.Channel.Close()

	// check that given configuration has neccessary variables
	if !isValidExchangeConfig(config) {
		log.Error("received invalid exchange configuration")
		return InvalidRabbitMQConfig
	}

	// declare events exchange with fanout type
	err = conn.Channel.ExchangeDeclare(config.ExchangeName, config.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq exchange: %s", err))
		return err
	}
	// construct new message and send over rabbitMQ channel
	message := amqp.Publishing{Body: payload}
	err = conn.Channel.Publish(config.ExchangeName, "", false, false, message)
	if err != nil {
		log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
		return err
	}
	log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ exchange", message))
	return nil
}