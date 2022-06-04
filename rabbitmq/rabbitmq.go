// A simple amqp library for connecting to RabbitMQ
//
// This is a wrapper around the github.com/streadway/amqp library.
//
package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQ struct {
	// The amqp url to connect to
	Url string `yaml:"url"`
	// The schange for publishing, defaults to amq.topic
	Exchange string `yaml:"exchange"`
	// The name of the connection that appears in the management plugin
	ConnectionName string `yaml:"connectionName"`
	// The heartBeat in seconds. Defaults to 10
	HeartBeat int `yaml:"heartBeat"`
	// The product name in the management plugin (optional)
	Product string `yaml:"product"`
	// The product version in the management plugin (optional)
	Version string `yaml:"version"`
	// ===== Internal
	connection *amqp.Connection `yaml:"-"` // amqp connection
	channel    *amqp.Channel    `yaml:"-"` // amqp channel
}

// called by main() ensure mandatory config is present
func (s *RabbitMQ) url() string {
	if s.Url == "" {
		log.Fatal("amqp.url is mandatory")
	}
	return s.Url
}

func (s *RabbitMQ) exchange() string {
	if s.Exchange == "" {
		return "amq.topic"
	}
	return s.Exchange
}

// Connect connects to the RabbitMQ instace thats been configured.
func (s *RabbitMQ) Connect() error {
	if s.connection != nil {
		return nil
	}

	log.Println("Connecting to amqp")

	var heartBeat = s.HeartBeat
	if heartBeat == 0 {
		heartBeat = 10
	}

	var product = s.Product
	if product == "" {
		product = "Area51 GO"
	}

	var version = s.Version
	if version == "" {
		version = "0.3β"
	}

	// Use the user provided client name
	if connection, err := amqp.DialConfig(s.url(), amqp.Config{
		Heartbeat: time.Duration(heartBeat) * time.Second,
		Properties: amqp.Table{
			"product":         product,
			"version":         version,
			"connection_name": s.ConnectionName,
		},
		Locale: "en_US",
	}); err != nil {
		return err
	} else {
		s.connection = connection
	}

	if channel, err := s.NewChannel(); err != nil {
		return err
	} else {
		s.channel = channel
	}

	log.Println("AMQP Connected")

	return s.channel.ExchangeDeclare(s.exchange(), "topic", true, false, false, false, nil)
}

func (s *RabbitMQ) NewChannel() (*amqp.Channel, error) {
	if channel, err := s.connection.Channel(); err != nil {
		return nil, err
	} else {
		return channel, nil
	}
}

// Publish a message
func (s *RabbitMQ) Publish(routingKey string, msg []byte) error {
	return s.channel.Publish(
		s.exchange(),
		routingKey,
		false,
		false,
		amqp.Publishing{
			Body: msg,
		})
}

// QueueDeclare declares a queue
func (r *RabbitMQ) QueueDeclare(channel *amqp.Channel, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

// QueueBind binds a queue to an exchange & routingKey
func (r *RabbitMQ) QueueBind(channel *amqp.Channel, name, key, exchange string, noWait bool, args amqp.Table) error {
	return channel.QueueBind(name, key, exchange, noWait, args)
}

// Consume adds a consumer to the queue and returns a GO channel
func (r *RabbitMQ) Consume(channel *amqp.Channel, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
