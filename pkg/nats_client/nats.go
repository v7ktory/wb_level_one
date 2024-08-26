package natsclient

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	defaultMaxReconnects = 2
	defaultReconnectWait = 2 * time.Second
	defaultConnTimeout   = 2 * time.Second
)

type NATS struct {
	maxReconnects int
	reconnectWait time.Duration
	connTimeout   time.Duration

	Conn *nats.Conn
}

func New(url string, opts ...Option) (*NATS, error) {
	const op = "natsclient.nats.go - New"
	n := &NATS{
		maxReconnects: defaultMaxReconnects,
		reconnectWait: defaultReconnectWait,
		connTimeout:   defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(n)
	}

	options := []nats.Option{
		nats.MaxReconnects(n.maxReconnects),
		nats.ReconnectWait(n.reconnectWait),
		nats.Timeout(n.connTimeout),

		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS: client disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS: client reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS: client closed")
		}),
		nats.ConnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS: client connected")
		}),
	}

	var err error
	n.Conn, err = nats.Connect(url, options...)
	if err != nil {
		return nil, fmt.Errorf("%s - nats.Connect: %w", op, err)
	}
	return n, nil
}
func (n *NATS) Close() {
	if n.Conn != nil {
		if err := n.Conn.Drain(); err != nil {
			fmt.Printf("Error draining connection: %v\n", err)
			n.Conn.Close()
		}
	}
}
