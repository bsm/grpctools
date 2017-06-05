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
		opts = new(Options)
	}
	return grpc.Dial(ctx, target, opts.grpcDialOpts()...)
}

// --------------------------------------------------------------------

// DialOptions represent dial options.
type DialOptions struct {
	// SkipInsecure enables transport security.
	SkipInsecure bool
	// SkipBlock makes Dial non-blocking (Dial won't wait for connection to be up before returning).
	SkipBlock bool

	// LBAddr specifies github.com/bsm/grpclb balancer address, defaults to 127.0.0.1:8383.
	LBAddr string
	// SkipLB disables load-balancing.
	SkipLB bool

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

	if !o.SkipLB {
		lbAddr := o.LBAddr
		if lbAddr == "" {
			lbAddr = "127.0.0.1:8383"
		}
		balancer := grpclb.PickFirst(&grpclb.Options{
			Address: lbAddr,
		})
		opts = append(opts, grpc.WithBalancer(balancer))
	}

	if !o.SkipBackoff {
		backoffCfg := o.BackoffConfig
		if backoffCfg.MaxDelay == 0 {
			backoffCfg.MaxDelay = 30 * time.Second
		}
		opts = append(opts, grpc.WithBackoffConfig(backoffCfg))
	}

	return opts
}
