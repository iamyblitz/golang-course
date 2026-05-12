package grpc

import (
	"context"
	"errors"

	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRepositoryInfo(
	ctx context.Context,
	req *collectorpb.GetRepositoryInfoRequest,
) (*collectorpb.GetRepositoryInfoResponse, error) {
	s.log.Debug("collector repository info request received", "owner", req.GetOwner(), "repo", req.GetRepo())

	if req.GetOwner() == "" || req.GetRepo() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner and repo are required")
	}

	info, err := s.repository.GetInfo(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		if errors.Is(err, usecase.ErrRepositoryNotFound) {
			return nil, status.Error(codes.NotFound, "repository not found")
		}

		s.log.Error("failed to get repository info", "error", err)
		return nil, status.Error(codes.Internal, "failed to get repository info")
	}

	return &collectorpb.GetRepositoryInfoResponse{
		FullName:    info.FullName,
		Description: info.Description,
		Stars:       info.Stars,
		Forks:       info.Forks,
		CreatedAt:   info.CreatedAt,
	}, nil
}
