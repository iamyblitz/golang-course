package grpc

import (
	"context"

	pb "github.com/iamyblitz/golang-course/task2/proto"
)

type Adapter struct {
	client pb.RepoServiceClient
}

func New(client pb.RepoServiceClient) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GetRepo(ctx context.Context, owner, repo string) (*pb.RepoResponse, error) {
	return a.client.GetRepo(
		ctx,
		&pb.RepoRequest{
			Owner: owner,
			Repo:  repo,
		},
	)
}
