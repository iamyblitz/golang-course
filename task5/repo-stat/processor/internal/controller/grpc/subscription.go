package grpc

import (
	"context"
	"errors"

	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSubscriptionsInfo(
	ctx context.Context,
	_ *processorpb.GetSubscriptionsInfoRequest,
) (*processorpb.GetSubscriptionsInfoResponse, error) {
	s.log.Debug("processor subscriptions info request received")

	repositories, err := s.repository.GetSubscriptionsInfo(ctx)
	if err != nil {
		if errors.Is(err, usecase.ErrRepositoryNotFound) {
			return nil, status.Error(codes.NotFound, "repository not found")
		}
		if errors.Is(err, usecase.ErrRepositoryUnavailable) {
			return nil, status.Error(codes.Unavailable, "repository unavailable")
		}
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		s.log.Error("failed to get subscriptions info", "error", err)
		return nil, status.Error(codes.Internal, "failed to get subscriptions info")
	}

	response := &processorpb.GetSubscriptionsInfoResponse{
		Repositories: make([]*processorpb.RepositoryInfo, 0, len(repositories)),
	}
	for _, repository := range repositories {
		response.Repositories = append(response.Repositories, &processorpb.RepositoryInfo{
			FullName:    repository.FullName,
			Description: repository.Description,
			Stars:       repository.Stars,
			Forks:       repository.Forks,
			CreatedAt:   repository.CreatedAt,
		})
	}

	return response, nil
}
