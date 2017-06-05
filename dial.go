package grpctools

import (
	"time"

	"github.com/bsm/grpclb"
	"golang.org/x/net/context"
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
	// SkipInsecure enables transport security.
	SkipInsecure bool
	// SkipBlock makes Dial non-blocking (Dial won't wait for connection to be up before returning).
	SkipBlock bool

	// LBAddr specifies github.com/bsm/grpclb balancer address, optional (no load-balancing unless provided).
	LBAddr string

	// BackoffConfig specifies backoff config, MaxDelay defaults to 30 seconds.
	BackoffConfig grpc.BackoffConfig
	// SkipBackoff disables backoff.
	SkipBackoff bool
}

func (o *DialOptions) grpcDialOpts() (opts []grpc.DialOption) {
	if !o.SkipInsecure {
		opts = append(opts, grpc.WithInsecure())
	}

	if !o.SkipBlock {
		opts = append(opts, grpc.WithBlock())
	}

	if o.LBAddr != "" {
		balancer := grpclb.PickFirst(&grpclb.Options{
			Address: o.LBAddr,
		})
		opts = append(opts, grpc.WithBalancer(balancer))
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
