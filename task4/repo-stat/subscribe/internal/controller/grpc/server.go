package grpc

import (
	"log/slog"

	subscribepb "repo-stat/proto/subscribe"
	"repo-stat/subscribe/internal/usecase"
)

type Server struct {
	subscribepb.UnimplementedSubscribeServer
	log          *slog.Logger
	subscription *usecase.Subscription
}

func NewServer(log *slog.Logger, subscription *usecase.Subscription) *Server {
	return &Server{
		log:          log,
		subscription: subscription,
	}
}
