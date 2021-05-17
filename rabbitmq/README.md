# rabbitmq
--
    import "github.com/peter-mount/golib/rabbitmq"

A simple amqp library for connecting to RabbitMQ

This is a wrapper around the github.com/streadway/amqp library.

## Usage

#### type RabbitMQ

```go
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
}
```


#### func (*RabbitMQ) Connect

```go
func (s *RabbitMQ) Connect() error
```
Connect connects to the RabbitMQ instace thats been configured.

#### func (*RabbitMQ) Consume

```go
func (r *RabbitMQ) Consume(channel *amqp.Channel, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
```
Consume adds a consumer to the queue and returns a GO channel

#### func (*RabbitMQ) NewChannel

```go
func (s *RabbitMQ) NewChannel() (*amqp.Channel, error)
```

#### func (*RabbitMQ) Publish

```go
func (s *RabbitMQ) Publish(routingKey string, msg []byte) error
```
Publish a message

#### func (*RabbitMQ) QueueBind

```go
func (r *RabbitMQ) QueueBind(channel *amqp.Channel, name, key, exchange string, noWait bool, args amqp.Table) error
```
QueueBind binds a queue to an exchange & routingKey

#### func (*RabbitMQ) QueueDeclare

```go
func (r *RabbitMQ) QueueDeclare(channel *amqp.Channel, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
```
QueueDeclare declares a queue
