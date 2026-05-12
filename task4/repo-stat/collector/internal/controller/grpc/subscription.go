package grpc

import (
	"context"

	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSubscriptionsInfo(
	ctx context.Context,
	_ *collectorpb.GetSubscriptionsInfoRequest,
) (*collectorpb.GetSubscriptionsInfoResponse, error) {
	s.log.Debug("collector subscriptions info request received")

	repositories, err := s.repository.GetSubscriptionsInfo(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		s.log.Error("failed to get subscriptions info", "error", err)
		return nil, status.Error(codes.Internal, "failed to get subscriptions info")
	}

	response := &collectorpb.GetSubscriptionsInfoResponse{
		Repositories: make([]*collectorpb.RepositoryInfo, 0, len(repositories)),
	}
	for _, repository := range repositories {
		response.Repositories = append(response.Repositories, &collectorpb.RepositoryInfo{
			FullName:    repository.FullName,
			Description: repository.Description,
			Stars:       repository.Stars,
			Forks:       repository.Forks,
			CreatedAt:   repository.CreatedAt,
		})
	}

	return response, nil
}
