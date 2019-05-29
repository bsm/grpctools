package grpctools

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server embeds a standard grpc Server with a healthcheck
type Server struct {
	*grpc.Server

	name   string
	addr   string
	health *health.Server
}

// NewServer returns a new Server instance.
func NewServer(name string, addr string, opts *Options) *Server {
	if opts == nil {
		opts = new(Options)
	}

	srv := &Server{
		Server: grpc.NewServer(opts.grpcServerOpts()...),
		name:   name,
		addr:   addr,
		health: health.NewServer(),
	}
	healthpb.RegisterHealthServer(srv.Server, srv.health)
	return srv
}

// ListenAndServe starts the server (blocking).
func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	s.health.SetServingStatus(s.name, healthpb.HealthCheckResponse_SERVING)
	err = s.Serve(lis)
	s.health.SetServingStatus(s.name, healthpb.HealthCheckResponse_NOT_SERVING)
	return err
}

// --------------------------------------------------------------------

// Options represent server options
type Options struct {
	// MaxConcurrentStreams will apply a limit on the number of concurrent streams to each ServerTransport.
	MaxConcurrentStreams uint32
	// MaxMessageSize sets the max message size in bytes for inbound mesages.
	// If this is not set, gRPC uses the default 4MB.
	MaxMessageSize int

	SkipCompression     bool
	SkipInstrumentation bool

	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

func (o *Options) grpcServerOpts() []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0)
	uchain := append([]grpc.UnaryServerInterceptor{}, o.UnaryInterceptors...)
	schain := append([]grpc.StreamServerInterceptor{}, o.StreamInterceptors...)

	if o.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(o.MaxConcurrentStreams))
	}

	if o.MaxMessageSize > 0 {
		opts = append(opts, grpc.MaxMsgSize(o.MaxMessageSize))
	}

	if !o.SkipCompression {
		opts = append(opts, grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))
	}

	if !o.SkipInstrumentation {
		uchain = append(uchain, DefaultInstrumenter.UnaryServerInterceptor)
		schain = append(schain, DefaultInstrumenter.StreamServerInterceptor)
	}

	if len(uchain) != 0 {
		opts = append(opts, grpc.UnaryInterceptor(unaryServerInterceptorChain(uchain...)))
	}

	if len(schain) != 0 {
		opts = append(opts, grpc.StreamInterceptor(streamServerInterceptorChain(schain...)))
	}

	return opts
}
