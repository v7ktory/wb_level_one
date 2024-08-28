package natsjs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	jetStr jetstream.JetStream
	logger *slog.Logger
}

func NewPublisher(jetStr jetstream.JetStream, logger *slog.Logger) *Publisher {
	return &Publisher{
		jetStr: jetStr,
		logger: logger,
	}
}

// TODO we can add a publish method to produce messages into NATS stream
func (p *Publisher) CreateStream(ctx context.Context, streamName, subject string) (jetstream.Stream, error) {
	const op = "subscriber.subscriber.go - createStream"
	stream, err := p.jetStr.CreateStream(ctx, jetstream.StreamConfig{
		Name:              streamName,
		Subjects:          []string{subject},
		Retention:         jetstream.InterestPolicy, // remove acked messages
		Discard:           jetstream.DiscardOld,     // when the stream is full, discard old messages
		MaxAge:            7 * 24 * time.Hour,       // max age of stored messages is 7 days
		Storage:           jetstream.FileStorage,    // type of message storage
		MaxMsgsPerSubject: 100_000_000,              // max stored messages per subject
		MaxMsgSize:        4 << 20,                  // max single message size is 4 MB
		NoAck:             false,
	})
	if err != nil {
		p.logger.Error("Failed to create stream", slog.Any("error", err.Error()), slog.Any("operation", op))
		return nil, fmt.Errorf("%s - jetstream.CreateStream: %w", op, err)
	}

	return stream, nil
}
