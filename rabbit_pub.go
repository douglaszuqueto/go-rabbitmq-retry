package rabbit

import (
	"sync"
	"time"

	"github.com/rafaeljesus/retry-go"
	"github.com/streadway/amqp"
)

// NewPublisher NewPublisher
func NewPublisher(client *Client) (*Publisher, error) {
	c := &Publisher{
		cli: client,
		rw:  sync.Mutex{},
	}

	channel, err := c.cli.Channel(func() {})
	if err != nil {
		return nil, err
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	c.ch = channel

	return c, nil
}

// Publish Publish
func (c *Publisher) Publish(exchange, routingKey string, body []byte) error {
	attempts := 5
	sleepTime := time.Second * 2

	return retry.Do(c.publish(exchange, routingKey, body),
		attempts,
		sleepTime,
	)
}

func (c *Publisher) publish(exchange, routingKey string, body []byte) func() error {
	return func() error {
		return c.ch.GetChannel().Publish(
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Body: body,
			},
		)
	}
}

// Close Close
func (c *Publisher) Close() error {
	return c.ch.Close()
}
