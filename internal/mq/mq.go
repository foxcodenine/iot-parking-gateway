package mq

import (
	// "os"
	"fmt"
	"os"
	"time"
)

// Configuration for RabbitMQ connection
type RabbitConfig struct {
	URL            string
	ReconnectDelay time.Duration
	Queues         map[string]QueueConfig // Map of queues by name
}

func (r *RabbitConfig) GetAllExchanges() []Exchange {
	exchangeMap := make(map[string]Exchange)
	for _, queueConfig := range r.Queues {
		for _, exchange := range queueConfig.Exchanges {
			// Use exchange name as key to ensure uniqueness
			exchangeMap[exchange.Name] = exchange
		}
	}

	// Convert the map to a slice
	exchanges := make([]Exchange, 0, len(exchangeMap))
	for _, exchange := range exchangeMap {
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}

type QueueConfig struct {
	Name       string
	RoutingKey string
	Durable    bool
	Exchanges  []Exchange
}

type Exchange struct {
	Name string
}

func SetupRabbitMQConfig() RabbitConfig {
	// Setup your RabbitMQ configuration here
	var rabbitmqPort string

	switch os.Getenv("GO_ENV") {
	case "production":
		rabbitmqPort = os.Getenv("RABBITMQ_PORT")
	default:
		rabbitmqPort = os.Getenv("RABBITMQ_PORT_EX")
	}
	urlStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"), rabbitmqPort)
	return RabbitConfig{
		URL:            urlStr,
		ReconnectDelay: 10 * time.Second,
		Queues: map[string]QueueConfig{
			"test_queue": {
				Name:       "test_queue",
				RoutingKey: "test_queue",
				Durable:    true,
				Exchanges: []Exchange{
					{Name: "test_exchange"},
				},
			},
			"nb_iot_event_logs_queue": {
				Name:       "nb_iot_event_logs_queue",
				RoutingKey: "nb_iot_event_logs_queue",
				Durable:    true,
				Exchanges: []Exchange{
					{Name: "nb_iot_event_logs_exchange"},
				},
			},
		},
	}
}
