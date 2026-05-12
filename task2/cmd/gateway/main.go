package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gatewaygrpc "github.com/iamyblitz/golang-course/task2/gateway/adapter/grpc"
	httpHandler "github.com/iamyblitz/golang-course/task2/gateway/handler/http"
	pb "github.com/iamyblitz/golang-course/task2/proto"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:2001"
	}

	conn, err := grpc.NewClient(
		collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewRepoServiceClient(conn)

	repoAdapter := gatewaygrpc.New(client)
	handler := httpHandler.New(repoAdapter, 3*time.Second)
	mux := http.NewServeMux()
	handler.Register(mux)

	log.Println("gateway started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
