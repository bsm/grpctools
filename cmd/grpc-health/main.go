package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var flags struct {
	Addr, Service string
	Timeout       time.Duration
}

func init() {
	flag.StringVar(&flags.Addr, "a", "127.0.0.1:8080", "Address to connect to. Default: 127.0.0.1:8080")
	flag.StringVar(&flags.Service, "s", "", "The service name. REQUIRED.")
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
	conn, err := grpc.Dial(flags.Addr, grpc.WithInsecure(), grpc.WithTimeout(flags.Timeout))
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
