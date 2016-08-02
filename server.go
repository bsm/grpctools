package grpctools

import (
	"net"
	"time"

	balancepb "github.com/bsm/grpclb/grpclb_backend_v1"
	"github.com/bsm/grpclb/load"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type LoadReportMeter interface {
	Increment(int64)
}

// Server embeds a standard grpc Server with a healthcheck
type Server struct {
	*grpc.Server
	LoadReportMeter

	name   string
	addr   string
	health *health.Server
}

// NewServer returns a new Server instance.
func NewServer(name string, addr string) *Server {
	lrs := load.NewRateReporter(time.Minute)
	srv := &Server{
		Server:          grpc.NewServer(),
		LoadReportMeter: lrs,
		name:            name,
		addr:            addr,
		health:          health.NewServer(),
	}
	healthpb.RegisterHealthServer(srv.Server, srv.health)
	balancepb.RegisterLoadReportServer(srv.Server, lrs)
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
