package grpc

import (
	"context"
	"errors"

	subscribepb "repo-stat/proto/subscribe"
	"repo-stat/subscribe/internal/domain"
	"repo-stat/subscribe/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateSubscription(
	ctx context.Context,
	req *subscribepb.CreateSubscriptionRequest,
) (*subscribepb.SubscriptionResponse, error) {
	s.log.Debug("create subscription request received", "owner", req.GetOwner(), "repo", req.GetRepo())

	subscription, err := s.subscription.Create(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		s.log.Error("failed to create subscription", "error", err)
		return nil, subscriptionError(err)
	}

	return &subscribepb.SubscriptionResponse{
		Subscription: subscriptionToProto(subscription),
	}, nil
}

func (s *Server) DeleteSubscription(
	ctx context.Context,
	req *subscribepb.DeleteSubscriptionRequest,
) (*subscribepb.DeleteSubscriptionResponse, error) {
	s.log.Debug("delete subscription request received", "owner", req.GetOwner(), "repo", req.GetRepo())

	if err := s.subscription.Delete(ctx, req.GetOwner(), req.GetRepo()); err != nil {
		return nil, subscriptionError(err)
	}

	return &subscribepb.DeleteSubscriptionResponse{}, nil
}

func (s *Server) ListSubscriptions(
	ctx context.Context,
	_ *subscribepb.ListSubscriptionsRequest,
) (*subscribepb.ListSubscriptionsResponse, error) {
	s.log.Debug("list subscriptions request received")

	subscriptions, err := s.subscription.List(ctx)
	if err != nil {
		return nil, subscriptionError(err)
	}

	response := &subscribepb.ListSubscriptionsResponse{
		Subscriptions: make([]*subscribepb.Subscription, 0, len(subscriptions)),
	}

	for _, subscription := range subscriptions {
		response.Subscriptions = append(response.Subscriptions, subscriptionToProto(subscription))
	}

	return response, nil
}

func subscriptionToProto(subscription domain.Subscription) *subscribepb.Subscription {
	return &subscribepb.Subscription{
		Owner: subscription.Owner,
		Repo:  subscription.Repo,
	}
}

func subscriptionError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrInvalidSubscription):
		return status.Error(codes.InvalidArgument, "owner and repo are required")
	case errors.Is(err, usecase.ErrSubscriptionExists):
		return status.Error(codes.AlreadyExists, "subscription already exists")
	case errors.Is(err, usecase.ErrSubscriptionNotFound):
		return status.Error(codes.NotFound, "subscription not found")
	case errors.Is(err, usecase.ErrRepositoryNotFound):
		return status.Error(codes.NotFound, "repository not found")
	case errors.Is(err, usecase.ErrRepositoryUnavailable):
		return status.Error(codes.Unavailable, "github api unavailable")
	default:
		return status.Error(codes.Internal, "failed to process subscription")
	}
}
