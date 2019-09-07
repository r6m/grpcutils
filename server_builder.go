package grpcutils

import (
	"google.golang.org/grpc"
)

type serverBuilder struct {
	options    []grpc.ServerOption
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
}

func NewServerBuilder() *serverBuilder {
	return &serverBuilder{
		options:    make([]grpc.ServerOption, 0),
		unaryInts:  make([]grpc.UnaryServerInterceptor, 0),
		streamInts: make([]grpc.StreamServerInterceptor, 0),
	}
}

func (b *serverBuilder) prepare() {
	b.AddOption(grpc.UnaryInterceptor(b.chainUnaryServer()))
	b.AddOption(grpc.StreamInterceptor(b.chainStreamServer()))
}

func (b *serverBuilder) Server() *grpc.Server {
	b.prepare()
	return grpc.NewServer(b.options...)
}
