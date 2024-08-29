package natsjs

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
	jetStr    jetstream.JetStream
	orderRepo pgdb.Order
	cache     cache.Cache[string, *entity.Order]
	logger    *slog.Logger
}

func NewSubscriber(jetStr jetstream.JetStream, orderRepo pgdb.Order, cache cache.Cache[string, *entity.Order], logger *slog.Logger) *Subscriber {
	return &Subscriber{
		jetStr:    jetStr,
		orderRepo: orderRepo,
		cache:     cache,
		logger:    logger,
	}
}

// Subscribe to NATS stream and consume incoming messages
func (s *Subscriber) Subscribe(ctx context.Context, c jetstream.Consumer) error {
	const op = "subscriber.subscriber.go - Subscribe"

	cons, err := c.Consume(func(msg jetstream.Msg) {
		defer msg.Ack()

		if err := s.handleMessage(ctx, msg.Data()); err != nil {
			s.logger.Error("Message handling error", slog.Any("error", err.Error()), slog.Any("operation", op))
		}
	})
	if err != nil {
		s.logger.Error("Failed to consume messages", slog.Any("error", err.Error()), slog.Any("operation", op))
		return fmt.Errorf("%s - jetstream.Consume: %w", op, err)
	}
	defer cons.Stop()

	<-ctx.Done()
	s.logger.Debug("Context canceled, stopping subscriber", slog.Any("operation", op))
	return nil
}

func (s *Subscriber) handleMessage(ctx context.Context, data []byte) error {
	const op = "subscriber.subscriber.go - handleMessage"

	orderRequest, problems, err := decodeNATSReq[model.Order](data)
	if err != nil {
		if len(problems) > 0 {
			for _, problem := range problems {
				s.logger.Error("Validation error", slog.Any("problem", problem), slog.Any("operation", op))
			}
			return nil
		}
		return fmt.Errorf("%s - decodeNATSReq: %w", op, err)
	}

	order := convertNATSReq(orderRequest)
	uid, err := s.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("%s - orderRepo.SaveOrder: %w", op, err)
	}

	s.cache.Put(uid, order)
	s.logger.Debug("Order saved successfully", slog.Any("order_uid", uid), slog.Any("operation", op))
	return nil
}
func (s *Subscriber) CreateConsumer(ctx context.Context, streamName, consumerName string) (jetstream.Consumer, error) {
	const op = "subscriber.subscriber.go - createConsumer"
	consumer, err := s.jetStr.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Durable:       consumerName,                // durable name is the same as consumer group name
		DeliverPolicy: jetstream.DeliverAllPolicy,  // deliver all messages, even if they were sent before the consumer was created
		AckPolicy:     jetstream.AckExplicitPolicy, // ack messages manually
		AckWait:       5 * time.Second,             // wait for ack for 5 seconds
		MaxAckPending: -1,
	})
	if err != nil {
		s.logger.Error("Failed to create consumer", slog.Any("error", err.Error()), slog.Any("operation", op))
		return nil, fmt.Errorf("%s - jetstream.CreateConsumer: %w", op, err)
	}
	return consumer, nil
}
