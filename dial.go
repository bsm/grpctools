package grpctools

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// Dial creates a client connection.
func Dial(target string, opts *DialOptions) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), target, opts)
}

// DialContext creates a client connection with specified context.
func DialContext(ctx context.Context, target string, opts *DialOptions) (*grpc.ClientConn, error) {
	if opts == nil {
		opts = new(DialOptions)
	}
	return grpc.DialContext(ctx, target, opts.grpcDialOpts()...)
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

	// Uses a load balancer, if provided.
	Balancer grpc.Balancer

	// Specifies backoff config, MaxDelay defaults to 30 seconds.
	BackoffConfig grpc.BackoffConfig

	// Disables backoff.
	SkipBackoff bool
}

func (o *DialOptions) grpcDialOpts() (opts []grpc.DialOption) {
	if !o.SkipInsecure {
		opts = append(opts, grpc.WithInsecure())
	}

	if !o.SkipBlock {
		opts = append(opts, grpc.WithBlock())
	}

	for _, cis := range o.UnaryInterceptors {
		opts = append(opts, grpc.WithUnaryInterceptor(cis))
	}
	for _, cis := range o.StreamInterceptors {
		opts = append(opts, grpc.WithStreamInterceptor(cis))
	}

	if o.Balancer != nil {
		opts = append(opts, grpc.WithBalancer(o.Balancer))
	}

	if !o.SkipBackoff {
		opts = append(opts, grpc.WithBackoffConfig(o.getBackoffConfig()))
	}

	return opts
}

func (o *DialOptions) getBackoffConfig() grpc.BackoffConfig {
	c := o.BackoffConfig
	if c.MaxDelay == 0 {
		c.MaxDelay = 30 * time.Second
	}
	return c
}
