package grpcutils

import (
	"context"

	"google.golang.org/grpc"
)

func (b *serverBuilder) AddUnaryInterceptor(intc grpc.UnaryServerInterceptor) *serverBuilder {
	b.unaryInts = append(b.unaryInts, intc)
	return b
}

func (b *serverBuilder) AddStreamInterceptor(intc grpc.StreamServerInterceptor) *serverBuilder {
	b.streamInts = append(b.streamInts, intc)
	return b
}

func (b *serverBuilder) chainUnaryServer() grpc.UnaryServerInterceptor {
	n := len(b.unaryInts)

	if n > 1 {
		lastI := n - 1
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			var (
				chainHandler grpc.UnaryHandler
				curI         int
			)

			chainHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				if curI == lastI {
					return handler(currentCtx, currentReq)
				}
				curI++
				resp, err := b.unaryInts[curI](currentCtx, currentReq, info, chainHandler)
				curI--
				return resp, err
			}

			return b.unaryInts[0](ctx, req, info, chainHandler)
		}
	}

	if n == 1 {
		return b.unaryInts[0]
	}

	// n == 0; Dummy interceptor maintained for backward compatibility to avoid returning nil.
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func (b *serverBuilder) chainStreamServer() grpc.StreamServerInterceptor {
	n := len(b.streamInts)

	if n > 1 {
		lastI := n - 1
		return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			var (
				chainHandler grpc.StreamHandler
				curI         int
			)

			chainHandler = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
				if curI == lastI {
					return handler(currentSrv, currentStream)
				}
				curI++
				err := b.streamInts[curI](currentSrv, currentStream, info, chainHandler)
				curI--
				return err
			}

			return b.streamInts[0](srv, stream, info, chainHandler)
		}
	}

	if n == 1 {
		return b.streamInts[0]
	}

	// n == 0; Dummy interceptor maintained for backward compatibility to avoid returning nil.
	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
