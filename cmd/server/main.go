package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	hellopb.UnimplementedGreetingServiceServer
}

func NewMyServer() *server {
	return &server{}
}

func (s *server) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())
	reflection.Register(s) // サーバーリフレクションの設定

	go func() {
		log.Printf("start gPRC server port: %v", port)
		s.Serve(listener)
	}()

	// Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
