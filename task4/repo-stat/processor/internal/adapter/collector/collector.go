package collector

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
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
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) GetRepositoryInfo(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error) {
	response, err := c.pb.GetRepositoryInfo(ctx, &collectorpb.GetRepositoryInfoRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		c.log.Error("collector repository info request failed", "error", err)
		return domain.RepositoryInfo{}, err
	}

	return domain.RepositoryInfo{
		FullName:    response.GetFullName(),
		Description: response.GetDescription(),
		Stars:       response.GetStars(),
		Forks:       response.GetForks(),
		CreatedAt:   response.GetCreatedAt(),
	}, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error) {
	response, err := c.pb.GetSubscriptionsInfo(ctx, &collectorpb.GetSubscriptionsInfoRequest{})
	if err != nil {
		c.log.Error("collector subscriptions info request failed", "error", err)
		return nil, err
	}

	repositories := make([]domain.RepositoryInfo, 0, len(response.GetRepositories()))
	for _, repository := range response.GetRepositories() {
		repositories = append(repositories, domain.RepositoryInfo{
			FullName:    repository.GetFullName(),
			Description: repository.GetDescription(),
			Stars:       repository.GetStars(),
			Forks:       repository.GetForks(),
			CreatedAt:   repository.GetCreatedAt(),
		})
	}

	return repositories, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
