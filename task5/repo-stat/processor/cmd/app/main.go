package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	kafkaadapter "repo-stat/processor/internal/adapter/kafka"
	"repo-stat/processor/internal/adapter/migrations"
	postgresadapter "repo-stat/processor/internal/adapter/postgres"
	subscribeadapter "repo-stat/processor/internal/adapter/subscribe"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/storage/db"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	if err := migrations.Up(cfg.Database.MigrationsPath, cfg.Database.DSN); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("create pgx pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	queries := db.New(pool)
	repositoryStore := postgresadapter.NewRepositoryStore(queries)

	taskProducer := kafkaadapter.NewTaskProducer(cfg.Kafka.Brokers, cfg.Kafka.TaskTopic, log)
	defer taskProducer.Close()

	resultConsumer := kafkaadapter.NewResultConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.ResultTopic,
		cfg.Kafka.ConsumerGroup,
		repositoryStore,
		log,
	)
	defer resultConsumer.Close()
	go func() {
		if err := resultConsumer.Run(ctx); err != nil {
			log.Error("processor result consumer stopped", "error", err)
		}
	}()

	subscribeClient, err := subscribeadapter.NewClient(cfg.Services.Subscribe, log)
	if err != nil {
		return fmt.Errorf("create subscribe client: %w", err)
	}
	defer subscribeClient.Close()

	pingUseCase := usecase.NewPing()
	repositoryUseCase := usecase.NewRepository(repositoryStore, taskProducer, subscribeClient)
	processorServer := grpccontroller.NewServer(log, pingUseCase, repositoryUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
