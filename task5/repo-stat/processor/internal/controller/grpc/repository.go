package grpc

import (
	"context"
	"errors"

	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRepositoryInfo(
	ctx context.Context,
	req *processorpb.GetRepositoryInfoRequest,
) (*processorpb.GetRepositoryInfoResponse, error) {
	s.log.Debug("processor repository info request received", "url", req.GetUrl())

	info, err := s.repository.GetInfo(ctx, req.GetUrl())
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRepositoryURL) {
			return nil, status.Error(codes.InvalidArgument, "invalid repository url")
		}
		if errors.Is(err, usecase.ErrRepositoryNotFound) {
			return nil, status.Error(codes.NotFound, "repository not found")
		}
		if errors.Is(err, usecase.ErrRepositoryUnavailable) {
			return nil, status.Error(codes.Unavailable, "repository unavailable")
		}
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		s.log.Error("failed to get repository info", "error", err)
		return nil, status.Error(codes.Internal, "failed to get repository info")
	}

	return &processorpb.GetRepositoryInfoResponse{
		FullName:    info.FullName,
		Description: info.Description,
		Stars:       info.Stars,
		Forks:       info.Forks,
		CreatedAt:   info.CreatedAt,
	}, nil
}
