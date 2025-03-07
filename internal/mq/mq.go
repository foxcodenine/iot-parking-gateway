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
	Queues         map[string]Queue
	Exchanges      map[string]Exchange
	RoutingKey     map[string]Bind
}

// func (r *RabbitConfig) GetAllExchanges() []Exchange {
// 	exchangeMap := make(map[string]Exchange)
// 	for _, queueConfig := range r.Queues {
// 		for _, exchange := range queueConfig.Exchanges {
// 			// Use exchange name as key to ensure uniqueness
// 			exchangeMap[exchange.Name] = exchange
// 		}
// 	}

// 	// Convert the map to a slice
// 	exchanges := make([]Exchange, 0, len(exchangeMap))
// 	for _, exchange := range exchangeMap {
// 		exchanges = append(exchanges, exchange)
// 	}

// 	return exchanges
// }

type QueueConfig struct {
	Name       string
	RoutingKey string
	Durable    bool
	Exchanges  []Exchange
}

type Exchange struct {
	Name    string
	Type    string // "direct" or "fanout"
	Durable bool
}

type Queue struct {
	Name    string
	Durable bool
}

type Bind struct {
	Exchange string
	Queue    string
}

var Queues = map[string]Queue{
	"event_logs": {
		Name:    "event_logs",
		Durable: true,
	},
	"thingsboard_event_logs": {
		Name:    "thingsboard_event_logs",
		Durable: true,
	},
}

var Excanges = map[string]Exchange{
	"event_logs": {
		Name:    "event_logs",
		Type:    "fanout",
		Durable: true,
	},
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

		Queues:    Queues,
		Exchanges: Excanges,

		RoutingKey: map[string]Bind{
			"event_logs": {
				Exchange: "event_logs",
				Queue:    "event_logs",
			},
			"thingsboard_event_logs": {
				Exchange: "event_logs",
				Queue:    "thingsboard_event_logs",
			},
		},
	}
}
