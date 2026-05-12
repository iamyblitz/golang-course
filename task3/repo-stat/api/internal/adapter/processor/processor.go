package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
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
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) GetRepositoryInfo(ctx context.Context, url string) (domain.RepositoryInfo, error) {
	response, err := c.pb.GetRepositoryInfo(ctx, &processorpb.GetRepositoryInfoRequest{
		Url: url,
	})
	if err != nil {
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

func (c *Client) Close() error {
	return c.conn.Close()
}
