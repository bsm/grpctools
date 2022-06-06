package grpctools

import (
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	srv := NewServer("test", "127.0.0.1:8080", nil)
	defer srv.Stop()

	if sub := srv.Server; sub == nil {
		t.Errorf("expected %v not to be nil", sub)
	}

	exp := reflect.TypeOf(&grpc.Server{})
	if got := reflect.TypeOf(srv.Server); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
