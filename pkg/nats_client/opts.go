package natsclient

import "time"

type Option func(*NATS)

func WithMaxReconnects(maxReconnects int) Option {
	return func(n *NATS) {
		n.maxReconnects = maxReconnects
	}
}

func WithReconnectWait(reconnectWait time.Duration) Option {
	return func(n *NATS) {
		n.reconnectWait = reconnectWait
	}
}

func WithConnTimeout(connTimeout time.Duration) Option {
	return func(n *NATS) {
		n.connTimeout = connTimeout
	}
}
