package mqproducer

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Configuration for RabbitMQ connection
type RabbitConfig struct {
	URL            string
	ReconnectDelay time.Duration
	ExchangeName   string
	Queues         map[string]QueueConfig // Map of queues by name
}

type QueueConfig struct {
	Name       string
	RoutingKey string
	Durable    bool
}

// RabbitMQProducer manages the connection and publishing to RabbitMQ
type RabbitMQProducer struct {
	config     RabbitConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

// ---------------------------------------------------------------------

// Run starts the connection and the publishing process
func (p *RabbitMQProducer) Run() {
	for {
		if err := p.connect(); err != nil {
			log.Println("Failed to connect to RabbitMQ:", err)
			time.Sleep(p.config.ReconnectDelay)
			continue
		}
		break
	}

	// Simulate message sending
	p.sendMessage("testQueue", "Hello, RabbitMQ!")
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

	// Declare the exchange
	err = p.channel.ExchangeDeclare(
		p.config.ExchangeName, // exchange
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
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

		err = p.channel.QueueBind(
			queue.Name,            // queue name
			queue.RoutingKey,      // routing key
			p.config.ExchangeName, // exchange
			false,                 // no-wait
			nil,                   // arguments
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// sendMessage sends a message to a specified queue
func (p *RabbitMQProducer) sendMessage(queueName, message string) {
	queueConfig, exists := p.config.Queues[queueName]
	if !exists {
		log.Printf("Queue configuration not found for %s\n", queueName)
		return
	}

	if err := p.channel.Publish(
		p.config.ExchangeName,  // exchange
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
