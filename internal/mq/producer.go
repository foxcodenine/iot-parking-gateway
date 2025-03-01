package mq

import (
	"errors"
	"fmt"

	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// RabbitMQProducer manages the connection and publishing to RabbitMQ
type RabbitMQProducer struct {
	config     RabbitConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

var AppRabbitMQProducer *RabbitMQProducer

// NewRabbitMQProducer creates a new producer instance
func NewRabbitMQProducer(config RabbitConfig) *RabbitMQProducer {
	AppRabbitMQProducer = &RabbitMQProducer{
		config: config,
	}

	return AppRabbitMQProducer
}

// ---------------------------------------------------------------------

// Run starts the connection and the publishing process
func (p *RabbitMQProducer) Run() {
	for {
		if p.connect() {
			helpers.LogInfo("Successfully connected to RabbitMQ")
			p.monitorConnection() // Start monitoring the connection for closures
			break                 // Exit the loop after a successful connection
		}

		// Wait before retrying the connection
		helpers.LogInfo(fmt.Sprintf("Retrying connection to RabbitMQ in %s...", p.config.ReconnectDelay))
		time.Sleep(p.config.ReconnectDelay)
	}
}

// connect handles the connection and channel setup, including declaring multiple queues
func (p *RabbitMQProducer) connect() bool {
	var err error
	p.connection, err = amqp.Dial(p.config.URL)
	if err != nil {
		helpers.LogError(err, "Failed to connect to RabbitMQ")
		return false
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		if p.connection != nil {
			p.connection.Close()
		}
		helpers.LogError(err, "Failed to open a channel")
		return false
	}

	if err := p.setupExchangesAndQueues(); err != nil {
		return false
	}

	return true
}

func (p *RabbitMQProducer) setupExchangesAndQueues() error {
	// Declare Exchanges
	for _, exchange := range p.config.Exchanges {
		err := p.channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	// Declare Queues
	for _, queue := range p.config.Queues {
		_, err := p.channel.QueueDeclare(queue.Name, queue.Durable, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	// Ending Queses to Excahnges
	for routingKey, bind := range p.config.RoutingKey {
		err := p.channel.QueueBind(bind.Queue, routingKey, bind.Exchange, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *RabbitMQProducer) monitorConnection() {
	go func() {
		for {
			reason, ok := <-p.connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				helpers.LogInfo("Channel and connection closed")
				break
			}
			helpers.LogError(fmt.Errorf("connection closed: %s", reason), "Trying to reconnect...")
			for {
				if p.connect() {
					helpers.LogInfo("Reconnection successful")
					return
				}
				time.Sleep(p.config.ReconnectDelay)
			}
		}
	}()
}

func (p *RabbitMQProducer) SendMessageToExchange(exchangeName, message string) {
	exchange, ok := p.config.Exchanges[exchangeName]
	if !ok {
		helpers.LogError(
			errors.New("exchange not found"),
			fmt.Sprintf("Exchange '%s' does not exist", exchangeName),
		)
		return
	}

	err := p.channel.Publish(
		exchange.Name, // exchange name to publish to
		"",            // routing key (empty if not needed)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			MessageId:   uuid.NewString(),
		},
	)

	if err != nil {
		helpers.LogError(
			err,
			fmt.Sprintf("Failed to publish message to exchange '%s'", exchangeName),
		)
	}
}

func (p *RabbitMQProducer) SendMessage(routingKeyName, exchangeName, message string) {
	routingKey, ok := p.config.RoutingKey[routingKeyName]
	if !ok {
		helpers.LogError(
			errors.New("routingKey not found"),
			fmt.Sprintf("Routing key '%s' does not exist", routingKeyName),
		)
		return
	}
	if routingKey.Exchange != exchangeName {
		helpers.LogError(
			errors.New("exchange mismatch"),
			fmt.Sprintf("Routing key '%s' is configured for exchange '%s', but received exchange '%s'",
				routingKeyName, routingKey.Exchange, exchangeName),
		)
		return
	}

	err := p.channel.Publish(
		exchangeName,
		routingKeyName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			MessageId:   uuid.NewString(),
		},
	)

	if err != nil {
		helpers.LogError(
			err,
			fmt.Sprintf("Failed to publish message using routing key '%s' on exchange '%s'", routingKeyName, exchangeName),
		)
	}
}

// Close cleanly closes the channel and connection
func (p *RabbitMQProducer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		p.connection.Close()
	}
}
