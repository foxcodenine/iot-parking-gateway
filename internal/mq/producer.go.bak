package mq

import (
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
	// Declare exchanges
	for _, exchange := range p.config.GetAllExchanges() {
		if err := p.channel.ExchangeDeclare(
			exchange.Name, exchange.Type, true, false, false, false, nil); err != nil {
			return err
		}
	}

	// Declare and bind queues
	for _, queue := range p.config.Queues {
		if _, err := p.channel.QueueDeclare(queue.Name, queue.Durable, false, false, false, nil); err != nil {
			return err
		}
		for _, exchange := range queue.Exchanges {
			routingKey := queue.RoutingKey
			if exchange.Type == "fanout" {
				routingKey = "" // No routing key needed for fanout exchanges
			}
			if err := p.channel.QueueBind(queue.Name, routingKey, exchange.Name, false, nil); err != nil {
				return err
			}
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

// sendMessage sends a message to a specified queue using a direct exchange
func (p *RabbitMQProducer) SendMessage(exchangeName, queueName, message string) {

	var routingKey string = ""
	if queueName != "" {

		queueConfig, exists := p.config.Queues[queueName]
		if !exists {
			helpers.LogInfo("Queue configuration not found for %s", queueName)
			return
		}
		routingKey = queueConfig.RoutingKey
	}

	if err := p.channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			MessageId:   uuid.NewString(),
		},
	); err != nil {
		helpers.LogError(err, "Failed to publish a message to queue "+queueName)
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
