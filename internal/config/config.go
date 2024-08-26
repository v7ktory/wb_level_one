package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	// HTTP
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second

	// Postgres
	maxPoolSize  = 1
	connAttempts = 2
	connTimeout  = time.Second

	// NATS
	maxReconnects = 3
	reconnectWait = time.Second * 3
	timeout       = time.Second * 5
	streamName    = "example-stream"
	subject       = "example-subject"
	consumerName  = "example-consumer-group-name"
)

type (
	Config struct {
		HTTP HTTP
		PG   Postgres
		NATS NATS
	}

	HTTP struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	Postgres struct {
		URL          string
		MaxPoolSize  int
		ConnAttempts int
		ConnTimeout  time.Duration
	}
	NATS struct {
		URL           string
		MaxReconnects int
		ReconnectWait time.Duration
		Timeout       time.Duration

		StreamName   string
		Subject      string
		ConsumerName string
	}
)

func Load(path string) (config Config, err error) {
	err = godotenv.Load(path)
	if err != nil {
		log.Fatalf("unable to load .env file: %v", err)
	}
	// HTTP
	config.HTTP.Port = os.Getenv("HTTP_PORT")
	config.HTTP.ReadTimeout = readTimeout
	config.HTTP.WriteTimeout = writeTimeout

	// Postgres
	config.PG.URL = os.Getenv("PG_URL")
	config.PG.MaxPoolSize = maxPoolSize
	config.PG.ConnAttempts = connAttempts
	config.PG.ConnTimeout = connTimeout

	// NATS
	config.NATS.URL = os.Getenv("NATS_URL")
	config.NATS.MaxReconnects = maxReconnects
	config.NATS.ReconnectWait = reconnectWait
	config.NATS.Timeout = timeout

	// JetStream
	config.NATS.StreamName = streamName
	config.NATS.Subject = subject
	config.NATS.ConsumerName = consumerName

	return
}
