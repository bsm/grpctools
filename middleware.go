package grpctools

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryServerInterceptorChain combines multiple grpc.UnaryServerInterceptor
// functions into one chain.
func unaryServerInterceptorChain(fns ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	wrap := func(fn grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return fn(ctx, req, info, next)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		chain := handler
		for i := len(fns) - 1; i >= 0; i-- {
			chain = wrap(fns[i], info, chain)
		}
		return chain(ctx, req)
	}
}

// StreamServerInterceptorChain combines multiple grpc.StreamServerInterceptor
// functions into one chain.
func streamServerInterceptorChain(fns ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	wrap := func(fn grpc.StreamServerInterceptor, info *grpc.StreamServerInfo, next grpc.StreamHandler) grpc.StreamHandler {
		return func(srv interface{}, stream grpc.ServerStream) error {
			return fn(srv, stream, info, next)
		}
	}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		chain := handler
		for i := len(fns) - 1; i >= 0; i-- {
			chain = wrap(fns[i], info, chain)
		}
		return chain(srv, stream)
	}
}
