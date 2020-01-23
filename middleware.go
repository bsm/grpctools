package grpctools

import (
	"context"

	"google.golang.org/grpc"
)

// unaryServerInterceptorChain combines multiple grpc.UnaryServerInterceptor
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

// streamServerInterceptorChain combines multiple grpc.StreamServerInterceptor
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

// unaryClientInterceptorChain combines multiple grpc.UnaryClientInterceptor
// functions into one chain.
func unaryClientInterceptorChain(fns ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		wrap := func(fn grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(ctxN context.Context, methodN string, reqN, replyN interface{}, ccN *grpc.ClientConn, optsN ...grpc.CallOption) error {
				return fn(ctxN, methodN, reqN, replyN, ccN, next, optsN...)
			}
		}

		chain := invoker
		for i := len(fns) - 1; i >= 0; i-- {
			chain = wrap(fns[i], chain)
		}
		return chain(ctx, method, req, reply, cc, opts...)
	}
}

// streamClientInterceptorChain combines multiple grpc.StreamClientInterceptor
// functions into one chain.
func streamClientInterceptorChain(fns ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		wrap := func(fn grpc.StreamClientInterceptor, next grpc.Streamer) grpc.Streamer {
			return func(ctxN context.Context, descN *grpc.StreamDesc, ccN *grpc.ClientConn, methodN string, optsN ...grpc.CallOption) (grpc.ClientStream, error) {
				return fn(ctxN, descN, ccN, methodN, next, optsN...)
			}
		}

		chain := streamer
		for i := len(fns) - 1; i >= 0; i-- {
			chain = wrap(fns[i], chain)
		}
		return chain(ctx, desc, cc, method, opts...)
	}
}
