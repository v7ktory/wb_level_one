package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/v7ktory/wb_task_one/internal/config"
	v1 "github.com/v7ktory/wb_task_one/internal/controller/http/v1"
	natsjs "github.com/v7ktory/wb_task_one/internal/controller/nats_js"
	"github.com/v7ktory/wb_task_one/internal/entity"
	httpserver "github.com/v7ktory/wb_task_one/internal/http_server"
	"github.com/v7ktory/wb_task_one/internal/repo/cache"
	"github.com/v7ktory/wb_task_one/internal/repo/pgdb"
	"github.com/v7ktory/wb_task_one/pkg/logger"
	natsclient "github.com/v7ktory/wb_task_one/pkg/nats_client"
	"github.com/v7ktory/wb_task_one/pkg/postgres"
)

func Run() {
	// Configuration
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logger := logger.NewLogger(slog.LevelDebug)

	// Postgres
	logger.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize), postgres.ConnAttempts(cfg.PG.ConnAttempts), postgres.ConnTimeout(cfg.PG.ConnTimeout))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// PgRepo
	logger.Info("Initializing pgRepo...")
	pgRepo := pgdb.NewPgRepo(pg)

	// CacheRepo
	logger.Info("Initializing cacheRepo...")
	cacheRepo := cache.NewLRUCache[string, *entity.Order](1_073_741_824) // Cache capacity = 1GB

	// Cache Warmup
	logger.Info("Warming up cache...")
	err = cache.Warmup(context.Background(), pgRepo, cacheRepo)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - cache.Warmup: %w", err))
	}

	// NATS
	logger.Info("Initializing NATS...")
	n, err := natsclient.New(cfg.NATS.URL, natsclient.WithMaxReconnects(cfg.NATS.MaxReconnects), natsclient.WithReconnectWait(cfg.NATS.ReconnectWait), natsclient.WithConnTimeout(cfg.NATS.Timeout))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - nats.New: %w", err))
	}
	defer n.Close()

	// JetStream
	logger.Info("Initializing JetStream...")
	js, err := jetstream.New(n.Conn)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - n.Conn.JetStream: %w", err))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Publisher
	logger.Info("Initializing publisher...")
	pub := natsjs.NewPublisher(js, logger)

	// Create NATS stream
	logger.Info("Creating NATS stream...")
	_, err = pub.CreateStream(ctx, cfg.NATS.StreamName, cfg.NATS.Subject)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - sub.CreateStream: %w", err))
	}

	// Subscriber
	logger.Info("Initializing subscriber...")
	sub := natsjs.NewSubscriber(js, pgRepo, cacheRepo, logger)

	// Create NATS consumer
	logger.Info("Creating NATS consumer...")
	c, err := sub.CreateConsumer(ctx, cfg.NATS.StreamName, cfg.NATS.ConsumerName)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - sub.CreateConsumer: %w", err))
	}

	// Subscribe to NATS stream and consume incoming messages
	go func() {
		err = sub.Subscribe(ctx, c)
		if err != nil {
			log.Fatal(fmt.Errorf("app - Run - sub.Subscribe: %w", err))
		}
	}()

	// Handlers
	mux := http.NewServeMux()
	v1.AddRoutes(mux, cacheRepo, pgRepo, logger)

	// HTTP server
	logger.Info("Starting http server...")
	logger.Debug("Server port", slog.Any("port", cfg.HTTP.Port))
	httpServer := httpserver.New(mux, httpserver.Port(cfg.HTTP.Port), httpserver.ReadTimeout(cfg.HTTP.ReadTimeout), httpserver.WriteTimeout(cfg.HTTP.WriteTimeout))

	// Waiting signal
	logger.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify: ", slog.Any("error", err.Error()))
	}

	// Graceful shutdown
	logger.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown: ", slog.Any("error", err.Error()))
	}
}
