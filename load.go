package grpctools

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

type LoadReportMeter interface {
	Increment(int64)
}

func UnaryLoadReporter(lrm LoadReportMeter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		lrm.Increment(1)
		return handler(ctx, req)
	}
}

func StreamLoadReporter(lrm LoadReportMeter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, loadReportingStream{
			ServerStream: stream,
			lrm:          lrm,
		})
	}
}

type loadReportingStream struct {
	grpc.ServerStream
	lrm LoadReportMeter
}

func (s loadReportingStream) RecvMsg(m interface{}) error {
	s.lrm.Increment(1)
	return s.ServerStream.RecvMsg(m)
}
