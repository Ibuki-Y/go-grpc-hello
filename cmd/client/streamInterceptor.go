package main

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func streamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("[pre] my stream client interceptor", method)
	stream, err := streamer(ctx, desc, cc, method, opts...)
	return &clientStreamWrapper{stream}, err
}

type clientStreamWrapper struct {
	grpc.ClientStream
}

// override
func (s *clientStreamWrapper) SendMsg(m interface{}) error {
	log.Println("[pre message] my stream client interceptor", m)
	return s.ClientStream.SendMsg(m)
}

// override
func (s *clientStreamWrapper) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor", m)
	}
	return err
}

// override
func (s *clientStreamWrapper) CloseSend() error {
	err := s.ClientStream.CloseSend()
	log.Println("[post] my stream client interceptor")
	return err
}
