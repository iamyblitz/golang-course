package main

import (
	"log"
	"net"

	pb "github.com/iamyblitz/golang-course/task2/proto"
	"google.golang.org/grpc"

	grpcHandler "github.com/iamyblitz/golang-course/task2/collector/handler/grpc"
	"github.com/iamyblitz/golang-course/task2/collector/usecase"
)

func main() {
	list, err := net.Listen("tcp", ":2001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	repoUC := usecase.NewRepoUsecase()
	pb.RegisterRepoServiceServer(s, grpcHandler.NewServer(repoUC))

	log.Println("collector started on :2001")

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
