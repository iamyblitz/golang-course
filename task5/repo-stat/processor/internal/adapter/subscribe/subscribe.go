package subscribe

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
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

func (c *Client) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	response, err := c.pb.ListSubscriptions(ctx, &subscribepb.ListSubscriptionsRequest{})
	if err != nil {
		c.log.Error("subscribe list subscriptions request failed", "error", err)
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(response.GetSubscriptions()))
	for _, subscription := range response.GetSubscriptions() {
		subscriptions = append(subscriptions, domain.Subscription{
			Owner: subscription.GetOwner(),
			Repo:  subscription.GetRepo(),
		})
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
