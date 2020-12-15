package rabbit

import (
	"sync"

	"github.com/streadway/amqp"
)

// SubFunction SubFunction
type SubFunction func(amqp.Delivery) (bool, error)

// NewSubscriber NewSubscriber
func NewSubscriber(client *Client) *Subscriber {
	c := &Subscriber{
		cli: client,
		rw:  sync.Mutex{},
	}

	return c
}

// Subscribe Subscribe
func (s *Subscriber) Subscribe(queueName string, sub SubFunction) error {
	channel, err := s.cli.Channel(func() {
		s.doSub()
	})
	if err != nil {
		return err
	}

	s.rw.Lock()
	s.ch = channel
	s.rw.Unlock()

	s.queueName = queueName
	s.fn = sub

	return s.doSub()
}

func (s *Subscriber) doSub() error {
	queue, err := s.consume(s.queueName)
	if err != nil {
		return err
	}

	go s.doProcess(queue)

	return nil
}

func (s *Subscriber) doProcess(queue <-chan amqp.Delivery) {
	for msg := range queue {
		reject, err := s.fn(msg)

		if err != nil {
			msg.Reject(reject)
			continue
		}

		if reject {
			msg.Reject(true)
			continue
		}

		msg.Ack(true)
	}
}

// Close Close
func (s *Subscriber) Close() error {
	return s.ch.Close()
}

func (s *Subscriber) consume(queue string) (<-chan amqp.Delivery, error) {
	return s.ch.GetChannel().Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
