package jackrabbit

import (
	"fmt"
	"errors"
	"github.com/streadway/amqp"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)


var (
	validate = validator.New()
	InvalidRabbitMQConfig = errors.New("invalid RabbitMQ configuration")
	RabbitMQConnectionError = errors.New("cannot connect to RabbitMQ server")
)


type RabbitConnectionConfig struct {
	QueueURL     string `json:"queue_url" validate:"required"`
	QueueName    string `json:"queue_name"`
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
}

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// function used to create new rabbitMQ channel using AMQP library.
// channels can the be used to create new queues and exchanges
func NewRabbitConnection(config RabbitConnectionConfig) (*RabbitMQConnection, error) {
	// validate configuration for missing fields
	err := validate.Struct(config)
	if err != nil {
		log.Error(fmt.Errorf("invalid RabbitMQ configuration: %v", err))
		return &RabbitMQConnection{}, InvalidRabbitMQConfig
	}

	log.Info(fmt.Sprintf("connecting to rabbitmq server at %s", config.QueueURL))
	// connect to rabbitmq server using queue url
	conn, err := amqp.Dial(config.QueueURL)
	if err != nil {
		log.Error(fmt.Errorf("unable to connect to rabbitmq server: %s", err))
		return &RabbitMQConnection{}, err
	}
	// create new rabbitMQ channel
	channel, err := conn.Channel()
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq channel: %s", err))
		// ensure that connection is properly closed if channel cannot be generated
		defer conn.Close()
		return &RabbitMQConnection{}, err
	}
	return &RabbitMQConnection{conn, channel}, nil
}