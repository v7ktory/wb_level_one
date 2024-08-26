package subscriber

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/v7ktory/wb_task_one/internal/entity"
	"github.com/v7ktory/wb_task_one/internal/model"
	"github.com/v7ktory/wb_task_one/internal/repo/cache"
	"github.com/v7ktory/wb_task_one/internal/repo/pgdb"
)

type Subscriber struct {
	jetStr jetstream.JetStream
	pgRepo *pgdb.PgRepo
	cache  cache.CacheRepo[string, *entity.Order]
	logger *slog.Logger
}

func New(jetStr jetstream.JetStream, pgRepo *pgdb.PgRepo, cache cache.CacheRepo[string, *entity.Order], logger *slog.Logger) *Subscriber {
	return &Subscriber{
		jetStr: jetStr,
		pgRepo: pgRepo,
		cache:  cache,
		logger: logger,
	}
}

// Create stream and subscribe to NATS stream and consume incoming messages
func (s *Subscriber) Subscribe(ctx context.Context, streamName, consumerName, subject string) error {
	const op = "subscriber.subscriber.go - Subscribe"
	_, err := s.createStream(ctx, streamName, subject)
	if err != nil {
		s.logger.Error("Failed to create stream", op, slog.Any("error", err.Error()))
		return fmt.Errorf("%s - createStream: %w", op, err)
	}

	c, err := s.createConsumer(ctx, streamName, consumerName)
	if err != nil {
		s.logger.Error("Failed to create consumer", op, slog.Any("error", err.Error()))
		return fmt.Errorf("%s - jetstream.CreateConsumer: %w", op, err)
	}

	cons, err := c.Consume(func(msg jetstream.Msg) {
		defer msg.Ack()

		s.logger.Debug("Message received", slog.Any("data", string(msg.Data())))

		// Message validation
		orderRequest, problems, err := DecodeNATSReq[model.Order](msg.Data())
		if err != nil {
			if len(problems) > 0 {
				for _, problem := range problems {
					s.logger.Error("Validation error", slog.Any("problem", problem))
				}
				return
			} else {
				s.logger.Error("Decode error", slog.Any("error", err.Error()))
				return
			}
		}

		// Convert NATS request and save in Postgres and cache
		order := ConvertNATSReq(orderRequest)

		uid, err := s.pgRepo.Save(ctx, order)
		if err != nil {
			s.logger.Error("Failed to save orderr", slog.Any("error", err.Error()))
			return
		}
		s.cache.Put(uid, order)

		s.logger.Debug("Order saved successfully", slog.Any("order_uid", uid))
	})
	if err != nil {
		s.logger.Error("Failed to consume messages", op, slog.Any("error", err.Error()))
		return fmt.Errorf("%s - jetstream.Consume: %w", op, err)
	}
	defer cons.Stop()

	<-ctx.Done()

	s.logger.Debug("Context canceled, stopping subscriber")
	return ctx.Err()
}

func (s *Subscriber) createStream(ctx context.Context, streamName, subject string) (jetstream.Stream, error) {
	const op = "subscriber.subscriber.go - createStream"
	stream, err := s.jetStr.CreateStream(ctx, jetstream.StreamConfig{
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
		s.logger.Error("Failed to create stream", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s - jetstream.CreateStream: %w", op, err)
	}

	return stream, nil
}

func (s *Subscriber) createConsumer(ctx context.Context, streamName, consumerName string) (jetstream.Consumer, error) {
	const op = "subscriber.subscriber.go - createConsumer"
	consumer, err := s.jetStr.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Durable:       consumerName,                // durable name is the same as consumer group name
		DeliverPolicy: jetstream.DeliverAllPolicy,  // deliver all messages, even if they were sent before the consumer was created
		AckPolicy:     jetstream.AckExplicitPolicy, // ack messages manually
		AckWait:       5 * time.Second,             // wait for ack for 5 seconds
		MaxAckPending: -1,
	})
	if err != nil {
		s.logger.Error("Failed to create consumer", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s - jetstream.CreateConsumer: %w", op, err)
	}
	return consumer, nil
}
