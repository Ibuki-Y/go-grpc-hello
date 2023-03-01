package main

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func streamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor: ", info.FullMethod)
	err := handler(srv, &serverStreamWrapper{ss}) // main処理
	log.Println("[post stream] my stream server interceptor finished!")
	return err
}

// grpc.ServerStreamインターフェースを満たす独自構造体
type serverStreamWrapper struct {
	grpc.ServerStream
}

// override
func (s *serverStreamWrapper) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m) // streamからmessageを受信
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor: ", m)
	}
	return err
}

// override
func (s *serverStreamWrapper) SendMsg(m interface{}) error {
	log.Println("[post message] my stream server interceptor: ", m)
	return s.ServerStream.SendMsg(m)
}
