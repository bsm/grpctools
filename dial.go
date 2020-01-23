package grpctools

import (
	"context"

	"google.golang.org/grpc"
)

// Dial creates a client connection.
func Dial(target string, opts *DialOptions, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, opts, extra...)
}

// DialContext creates a client connection with specified context.
func DialContext(ctx context.Context, target string, opts *DialOptions, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	if opts == nil {
		opts = new(DialOptions)
	}

	full := append(opts.grpcDialOpts(), extra...)
	return grpc.DialContext(ctx, target, full...)
}

// --------------------------------------------------------------------

// DialOptions represent dial options.
type DialOptions struct {
	// Enables transport security.
	SkipInsecure bool

	// Makes Dial non-blocking (Dial won't wait for connection to be up before returning).
	SkipBlock bool

	// Unary client interceptors.
	UnaryInterceptors []grpc.UnaryClientInterceptor

	// Stream client interceptors.
	StreamInterceptors []grpc.StreamClientInterceptor
}

func (o *DialOptions) grpcDialOpts() (opts []grpc.DialOption) {
	if !o.SkipInsecure {
		opts = append(opts, grpc.WithInsecure())
	}

	if !o.SkipBlock {
		opts = append(opts, grpc.WithBlock())
	}

	if chain := append([]grpc.UnaryClientInterceptor{}, o.UnaryInterceptors...); len(chain) != 0 {
		opts = append(opts, grpc.WithUnaryInterceptor(unaryClientInterceptorChain(chain...)))
	}

	if chain := append([]grpc.StreamClientInterceptor{}, o.StreamInterceptors...); len(chain) != 0 {
		opts = append(opts, grpc.WithStreamInterceptor(streamClientInterceptorChain(chain...)))
	}

	return opts
}
