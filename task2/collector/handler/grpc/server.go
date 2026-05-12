package grpc

import (
	"context"
	"errors"

	"github.com/iamyblitz/golang-course/task2/collector/usecase"
	pb "github.com/iamyblitz/golang-course/task2/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repoUsecase interface {
	GetRepo(ctx context.Context, owner, repo string) (*pb.RepoResponse, error)
}

type Server struct {
	pb.UnimplementedRepoServiceServer
	repoUC repoUsecase
}

func NewServer(repoUC repoUsecase) *Server {
	return &Server{repoUC: repoUC}
}

// реализация метода из proto
func (s *Server) GetRepo(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {

	// 1. валидация
	if req.Owner == "" || req.Repo == "" {
		return nil, status.Error(codes.InvalidArgument, "owner and repo required")
	}

	// 2. бизнес-логика в usecase
	resp, err := s.repoUC.GetRepo(ctx, req.Owner, req.Repo)
	if err != nil {
		if errors.Is(err, usecase.ErrRepoNotFound) {
			return nil, status.Error(codes.NotFound, "repository not found")
		}

		return nil, status.Error(codes.Internal, "failed to get repo")
	}

	return resp, nil
}
