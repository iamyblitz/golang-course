package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	kafkaadapter "repo-stat/collector/internal/adapter/kafka"
	subscribeadapter "repo-stat/collector/internal/adapter/subscribe"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
	"time"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector kafka worker...")
	log.Debug("debug messages are enabled")

	subscribeClient, err := subscribeadapter.NewClient(cfg.Services.Subscribe, log)
	if err != nil {
		return fmt.Errorf("create subscribe client: %w", err)
	}
	defer subscribeClient.Close()

	repositoryUseCase := usecase.NewRepository(subscribeClient)

	worker := kafkaadapter.NewWorker(
		cfg.Kafka.Brokers,
		cfg.Kafka.TaskTopic,
		cfg.Kafka.ResultTopic,
		cfg.Kafka.ConsumerGroup,
		repositoryUseCase,
		log,
	)
	defer worker.Close()

	taskProducer := kafkaadapter.NewTaskProducer(cfg.Kafka.Brokers, cfg.Kafka.TaskTopic)
	defer taskProducer.Close()

	refresher := usecase.NewSubscriptionRefresher(log, subscribeClient, taskProducer, 15*time.Second)
	go func() {
		if err := refresher.Run(ctx); err != nil {
			log.Error("subscription refresher stopped", "error", err)
		}
	}()

	if err := worker.Run(ctx); err != nil {
		return fmt.Errorf("run kafka worker: %w", err)
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
