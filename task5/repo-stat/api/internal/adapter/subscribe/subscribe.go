package subscribe

import (
	"context"
	"log/slog"

	"repo-stat/api/internal/domain"
	subscribepb "repo-stat/proto/subscribe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscribepb.SubscribeClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscribepb.NewSubscribeClient(conn),
	}, nil
}

func (c *Client) CreateSubscription(ctx context.Context, owner string, repo string) (domain.Subscription, error) {
	response, err := c.pb.CreateSubscription(ctx, &subscribepb.CreateSubscriptionRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return domain.Subscription{}, err
	}

	return subscriptionFromProto(response.GetSubscription()), nil
}

func (c *Client) DeleteSubscription(ctx context.Context, owner string, repo string) error {
	_, err := c.pb.DeleteSubscription(ctx, &subscribepb.DeleteSubscriptionRequest{
		Owner: owner,
		Repo:  repo,
	})
	return err
}

func (c *Client) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	response, err := c.pb.ListSubscriptions(ctx, &subscribepb.ListSubscriptionsRequest{})
	if err != nil {
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(response.GetSubscriptions()))
	for _, subscription := range response.GetSubscriptions() {
		subscriptions = append(subscriptions, subscriptionFromProto(subscription))
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func subscriptionFromProto(subscription *subscribepb.Subscription) domain.Subscription {
	if subscription == nil {
		return domain.Subscription{}
	}

	return domain.Subscription{
		Owner: subscription.GetOwner(),
		Repo:  subscription.GetRepo(),
	}
}
