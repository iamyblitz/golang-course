package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	subscribeadapter "repo-stat/collector/internal/adapter/subscribe"
	grpccontroller "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	collectorpb "repo-stat/proto/collector"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector server...")
	log.Debug("debug messages are enabled")

	subscribeClient, err := subscribeadapter.NewClient(cfg.Services.Subscribe, log)
	if err != nil {
		return fmt.Errorf("create subscribe client: %w", err)
	}
	defer subscribeClient.Close()

	repositoryUseCase := usecase.NewRepository(subscribeClient)
	collectorServer := grpccontroller.NewServer(log, repositoryUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorpb.RegisterCollectorServer(srv.GRPC(), collectorServer)

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
