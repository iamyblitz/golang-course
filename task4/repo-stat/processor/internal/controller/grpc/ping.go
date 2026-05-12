package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"
)

func (s *Server) Ping(ctx context.Context, _ *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("processor ping request received")

	return &processorpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
