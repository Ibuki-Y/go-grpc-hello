package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func unaryClientInterceptor(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, ops ...grpc.CallOption) error {
	fmt.Println("[pre] my unary client interceptor", method, req)
	err := invoker(ctx, method, req, res, cc, ops...) // main処理
	fmt.Println("[post] my unary client interceptor", res)
	return err
}
