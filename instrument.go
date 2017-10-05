package grpctools

import (
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bsm/rucksack/log"
	"github.com/bsm/rucksack/met"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var rx = regexp.MustCompile(`[\w\-]+\=[\w\-]+`)

var defaultCodeMap = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.Canceled:           http.StatusGone,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusPreconditionFailed,
	codes.AlreadyExists:      http.StatusConflict,
	codes.DeadlineExceeded:   http.StatusRequestTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.OutOfRange:         http.StatusRequestedRangeNotSatisfiable,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.DataLoss:           http.StatusInternalServerError,
}

// Instrumenter instances instrument RPC requests via
// interceptors.
type Instrumenter struct {
	metric  string
	codeMap map[codes.Code]int
}

// DefaultInstrumenter instruments via "rpc.request" metric
var DefaultInstrumenter = NewInstrumenter("rpc.request", nil)

// NewInstrumenter inits a new instrumenter with a metric
func NewInstrumenter(metric string, codeMap map[codes.Code]int) *Instrumenter {
	if metric == "" {
		metric = "rpc.request"
	}
	if codeMap == nil {
		codeMap = make(map[codes.Code]int, len(defaultCodeMap))
	}
	for c, h := range defaultCodeMap {
		if _, ok := codeMap[c]; !ok {
			codeMap[c] = h
		}
	}
	return &Instrumenter{metric: metric, codeMap: codeMap}
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
	status, ok := i.codeMap[grpc.Code(err)]
	if !ok {
		status = http.StatusTeapot
	}

	logger := log.WithField("status", status)
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
		log.WithField("rpc", name).Errorf("panic: %v\n%v", r, debug.Stack())
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

func loggerWithTags(ent *logrus.Entry, tags []string) *logrus.Entry {
	for _, t := range tags {
		if parts := strings.SplitN(t, ":", 2); len(parts) == 2 {
			ent = ent.WithField(parts[0], parts[1])
		}
	}
	return ent
}
