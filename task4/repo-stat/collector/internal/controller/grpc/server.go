package grpc

import (
	"log/slog"

	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"
)

type Server struct {
	collectorpb.UnimplementedCollectorServer
	log        *slog.Logger
	repository *usecase.Repository
}

func NewServer(log *slog.Logger, repository *usecase.Repository) *Server {
	return &Server{
		log:        log,
		repository: repository,
	}
}
