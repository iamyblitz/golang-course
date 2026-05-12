package grpc

import (
	"log/slog"

	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

type Server struct {
	processorpb.UnimplementedProcessorServer
	log        *slog.Logger
	ping       *usecase.Ping
	repository *usecase.Repository
}

func NewServer(log *slog.Logger, ping *usecase.Ping, repository *usecase.Repository) *Server {
	return &Server{
		log:        log,
		ping:       ping,
		repository: repository,
	}
}
