package rabbit

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/douglaszuqueto/go-rabbit-retry/pkg/logger"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Fn Fn
type Fn func()

func makeURL(cfg Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.IP,
		cfg.Port,
		cfg.VirtualHost,
	)
}

// New new
func New(cfg Config) (*Client, error) {
	url := makeURL(cfg)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn: conn,
		url:  url,
		rw:   sync.Mutex{},
	}

	go client.loop()

	return client, nil
}

// GetChannel GetChannel
func (s *Client) GetChannel() (*amqp.Channel, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.conn.Channel()
}

// Channel Channel
func (s *Client) Channel(fn Fn) (*Channel, error) {
	return NewChannel(s, fn)
}

func (s *Client) loop() {
WATCH:

	connErr, ok := <-s.conn.NotifyClose(make(chan *amqp.Error))
	if !ok || s.IsClosed() {
		s.Close()

		logger.Warning("rabbit.notify.error", zap.Bool("state", ok))
		return
	}

	if connErr == nil {
		logger.Warning("rabbit.conn.error", zap.String("error", "conexão perdida"))
	}

	logger.Warning("rabbit.notify.error", zap.String("error", connErr.Error()))

	for {
		conn, err := amqp.Dial(s.url)
		if err != nil {
			logger.Warning("rabbit.conn.reconnect", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)

			continue
		}

		s.rw.Lock()
		s.conn = conn
		s.rw.Unlock()

		logger.Info("rabbit.conn.retry")

		goto WATCH
	}
}

// Close lcose
func (s *Client) Close() error {
	logger.Info("Closing rabbit connection", zap.String("pkg", "rabbit"))

	if s.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&s.closed, 1)

	return s.conn.Close()
}

// NewConn NewConn
func NewConn() (*Client, error) {
	config := Config{
		IP:          os.Getenv("RABBITMQ_IP"),
		Port:        os.Getenv("RABBITMQ_PORT"),
		Username:    os.Getenv("RABBITMQ_USERNAME"),
		Password:    os.Getenv("RABBITMQ_PASSWORD"),
		VirtualHost: os.Getenv("RABBITMQ_VIRTUALHOST"),
	}

	return New(config)
}

// IsClosed IsClosed
func (s *Client) IsClosed() bool {
	return (atomic.LoadInt32(&s.closed) == 1)
}
