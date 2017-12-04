package grpctools

import (
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/bsm/rucksack/log"
	"github.com/bsm/rucksack/met"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var rx = regexp.MustCompile(`[\w\-]+\=[\w\-]+`)

// Instrumenter instances instrument RPC requests via
// interceptors.
type Instrumenter struct{ metric string }

// DefaultInstrumenter instruments via "rpc.request" metric
var DefaultInstrumenter = NewInstrumenter("rpc.request")

// NewInstrumenter inits a new instrumenter with a metric
func NewInstrumenter(metric string) *Instrumenter {
	if metric == "" {
		metric = "rpc.request"
	}
	return &Instrumenter{metric: metric}
}

// UnaryInterceptor implements an grpc.UnaryServerInterceptor
func (i *Instrumenter) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer i.recover(info.FullMethod)

	start := time.Now()
	resp, err := handler(ctx, req)
	i.instrument(info.FullMethod, err, time.Since(start))
	return resp, err
}

// StreamServerInterceptor implements an grpc.StreamServerInterceptor
func (i *Instrumenter) StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer i.recover(info.FullMethod)

	start := time.Now()
	err := handler(srv, stream)
	i.instrument(info.FullMethod, err, time.Since(start))
	return err
}

func (i *Instrumenter) instrument(name string, err error, elapsed time.Duration) {
	errtags := extractErrorTags(err)
	status := HTTPStatusFromError(err)
	logger := log.L().With(zap.Int("status", status)).Sugar()
	mtags := []string{
		"rpc:" + name,
		"status:" + strconv.Itoa(status),
	}
	met.RatePerMin(i.metric, mtags).Update(1)

	if status < 500 {
		logger.Infof("%s in %.3fs", name, elapsed.Seconds())
		met.Timer(i.metric+".time", mtags).Update(elapsed)
	} else if err != nil {
		loggerWithTags(logger, errtags).Errorf("%s failed with %s", name, err.Error())
		met.RatePerMin(i.metric+".error", append(mtags, errtags...)).Update(1)
	}
}

func (i *Instrumenter) recover(name string) {
	if r := recover(); r != nil {
		log.L().With(zap.String("rpc", name)).Sugar().Errorf("panic: %v\n%v", r, debug.Stack())
		met.RatePerMin(i.metric+".error", []string{
			"rpc:" + name,
			"status:500",
			"panic:true",
		}).Update(1)
	}
}

func extractErrorTags(err error) []string {
	if err == nil {
		return nil
	}

	tags := rx.FindAllString(err.Error(), -1)
	for i, tag := range tags {
		tags[i] = strings.Replace(tag, "=", ":", 1)
	}
	return tags
}

func loggerWithTags(l *zap.SugaredLogger, tags []string) *zap.SugaredLogger {
	for _, t := range tags {
		if parts := strings.SplitN(t, ":", 2); len(parts) == 2 {
			l = l.With(zap.String(parts[0], parts[1]))
		}
	}
	return l
}
