package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	subscribepb "repo-stat/proto/subscribe"
	"repo-stat/subscribe/config"
	githubadapter "repo-stat/subscribe/internal/adapter/github"
	"repo-stat/subscribe/internal/adapter/migrations"
	postgresadapter "repo-stat/subscribe/internal/adapter/postgres"
	grpccontroller "repo-stat/subscribe/internal/controller/grpc"
	"repo-stat/subscribe/internal/storage/db"
	"repo-stat/subscribe/internal/usecase"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscribe server...")
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
	subscriptionRepository := postgresadapter.NewSubscriptionRepository(queries)
	repositoryVerifier := githubadapter.NewVerifier()
	subscriptionUseCase := usecase.NewSubscription(subscriptionRepository, repositoryVerifier)
	subscribeServer := grpccontroller.NewServer(log, subscriptionUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscribepb.RegisterSubscribeServer(srv.GRPC(), subscribeServer)

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
