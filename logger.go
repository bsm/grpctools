package grpctools

import (
	"google.golang.org/grpc/grpclog"
)

// SetLogger is a proxy to grpclog.SetLogger
func SetLogger(logger grpclog.LoggerV2) {
	grpclog.SetLoggerV2(logger)
}
