package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var flags struct {
	Addr, Service string
	Timeout       time.Duration
	EnableTLS     bool
}

func init() {
	flag.StringVar(&flags.Addr, "a", "127.0.0.1:8080", "Address to connect to. Default: 127.0.0.1:8080")
	flag.StringVar(&flags.Service, "s", "", "The service name. REQUIRED.")
	flag.BoolVar(&flags.EnableTLS, "tls", false, "Enable client-side TLS")
	flag.DurationVar(&flags.Timeout, "timeout", 30*time.Second, "The request timeout. Default: 30s")
}

func main() {
	flag.Parse()

	if flags.Service == "" {
		flag.PrintDefaults()
		os.Exit(64)
	}

	os.Exit(run())
}

func run() int {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(flags.Timeout),
	}
	if flags.EnableTLS {
		opts[0] = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	conn, err := grpc.Dial(flags.Addr, opts...)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer conn.Close()

	req := &healthpb.HealthCheckRequest{Service: flags.Service}
	resp, err := healthpb.NewHealthClient(conn).Check(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println("Status:", resp.Status.String())
	switch resp.Status {
	case healthpb.HealthCheckResponse_SERVING:
		return 0
	case healthpb.HealthCheckResponse_NOT_SERVING:
		return 2
	default:
		return 1
	}
}
