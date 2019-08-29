package grpcutils

import (
	"google.golang.org/grpc"
)

func (b *serverBuilder) AddOption(o grpc.ServerOption) *serverBuilder {
	b.options = append(b.options, o)
	return b
}
