package rabbit

import (
	"sync/atomic"
	"time"

	"github.com/douglaszuqueto/go-rabbit-retry/pkg/logger"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// NewChannel NewChannel
func NewChannel(cli *Client, fn Fn) (*Channel, error) {
	c := &Channel{
		cli: cli,
	}

	if err := c.getChannel(); err != nil {
		return nil, err
	}

	go c.reconnect(fn)

	return c, nil
}

// GetChannel GetChannel
func (c *Channel) GetChannel() *amqp.Channel {
	c.rw.Lock()
	defer c.rw.Unlock()

	return c.ch
}

func (c *Channel) reconnect(fn Fn) {
WATCH:

	connErr, ok := <-c.ch.NotifyClose(make(chan *amqp.Error))
	if !ok || c.IsClosed() {
		c.Close()

		logger.Warning("rabbit.ch.error", zap.Bool("state", ok))
		return
	}

	if connErr == nil {
		logger.Warning("rabbit.ch.error", zap.String("error", "conexão perdida"))
		return
	}

	logger.Warning("rabbit.ch.error", zap.String("error", connErr.Error()))

	for {
		if err := c.getChannel(); err != nil {
			logger.Warning("rabbit.ch.getChannel", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)

			continue
		}

		logger.Info("rabbit.ch.retry")

		fn()

		goto WATCH
	}
}

func (c *Channel) getChannel() error {
	ch, err := c.cli.GetChannel()
	if err != nil {
		return err
	}

	c.rw.Lock()
	c.ch = ch
	c.rw.Unlock()

	return nil
}

// IsClosed IsClosed
func (c *Channel) IsClosed() bool {
	return (atomic.LoadInt32(&c.closed) == 1)
}

// Close ensure closed flag set
func (c *Channel) Close() error {
	if c.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&c.closed, 1)

	return c.ch.Close()
}
