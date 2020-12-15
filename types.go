package rabbit

import (
	"sync"

	"github.com/streadway/amqp"
)

// Config config
type Config struct {
	Username    string
	Password    string
	IP          string
	Port        string
	VirtualHost string
}

// Client client
type Client struct {
	conn *amqp.Connection
	url  string
	rw   sync.Mutex

	closed int32
}

// Channel Channel
type Channel struct {
	cli *Client
	ch  *amqp.Channel

	closed int32
	rw     sync.Mutex
}

// Publisher Publisher
type Publisher struct {
	cli *Client
	ch  *Channel

	rw sync.Mutex
}

// Subscriber Subscriber
type Subscriber struct {
	cli *Client
	ch  *Channel

	queueName string
	fn        SubFunction

	rw sync.Mutex
}
