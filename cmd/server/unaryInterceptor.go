package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("[pre] my unary server interceptor: ", info.FullMethod)
	res, err := handler(ctx, req) // main処理
	log.Println("[post] my unary server interceptor: ", res)
	return res, err
}
