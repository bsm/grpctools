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
	health *health.HealthServer
}

// NewServer returns a new Server instance.
func NewServer(name string) *Server {
	srv := &Server{
		Server: grpc.NewServer(),
		name:   name,
		health: health.NewHealthServer(),
	}
	healthpb.RegisterHealthServer(srv.Server, srv.health)
	return srv
}

// ListenAndServe starts the server (blocking).
func (s *Server) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	s.health.SetServingStatus(s.name, healthpb.HealthCheckResponse_SERVING)
	err = s.Serve(lis)
	s.health.SetServingStatus(s.name, healthpb.HealthCheckResponse_NOT_SERVING)
	return err
}
