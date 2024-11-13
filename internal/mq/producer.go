package mq

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// RabbitMQProducer manages the connection and publishing to RabbitMQ
type RabbitMQProducer struct {
	config     RabbitConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQProducer creates a new producer instance
func NewRabbitMQProducer(config RabbitConfig) *RabbitMQProducer {
	return &RabbitMQProducer{
		config: config,
	}
}

// ---------------------------------------------------------------------

// Run starts the connection and the publishing process
func (p *RabbitMQProducer) Run() {
	attempt := 0
	for {
		if err := p.connect(); err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v, attempt: %d\n", err, attempt)
			if attempt > 5 { // Maximum of 5 attempts
				log.Println("Max reconnect attempts reached, exiting.")
				return
			}
			time.Sleep(p.config.ReconnectDelay * time.Duration(attempt)) // Exponential backoff
			attempt++
			continue
		}
		break
	}

	// Assume operational status, proceed to send a message
	p.sendMessage("testExchange", "testQueue", "Hello, RabbitMQ!")
}

// connect handles the connection and channel setup, including declaring multiple queues
func (p *RabbitMQProducer) connect() error {
	var err error
	p.connection, err = amqp.Dial(p.config.URL)
	if err != nil {
		return err
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		p.connection.Close()
		return err
	}

	// Declare all the exchange
	for _, exchange := range p.config.GetAllExchanges() {
		err = p.channel.ExchangeDeclare(
			exchange.Name, // exchange
			"direct",      // type
			true,          // durable
			false,         // auto-deleted
			false,         // internal
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			return err
		}
	}

	// Declare and bind all queues
	for _, queue := range p.config.Queues {
		_, err = p.channel.QueueDeclare(
			queue.Name,    // queue name
			queue.Durable, // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			return err
		}

		for _, exchange := range queue.Exchanges {
			err = p.channel.QueueBind(
				queue.Name,       // queue name
				queue.RoutingKey, // routing key
				exchange.Name,    // exchange
				false,            // no-wait
				nil,              // arguments
			)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// sendMessage sends a message to a specified queue
func (p *RabbitMQProducer) sendMessage(exchangeName, queueName, message string) {
	queueConfig, exists := p.config.Queues[queueName]
	if !exists {
		log.Printf("Queue configuration not found for %s\n", queueName)
		return
	}

	if err := p.channel.Publish(
		exchangeName,           // exchange
		queueConfig.RoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			MessageId:   uuid.NewString(),
		},
	); err != nil {
		log.Printf("Failed to publish a message to queue %s: %s\n", queueName, err)
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
