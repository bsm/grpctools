package grpctools

import (
	"reflect"
	"sort"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestVErrors(t *testing.T) {
	var verr VErrors
	verr = verr.Reset()
	verr = verr.Append("name", "is required")
	verr = verr.Append("name", "must be 5 chars")
	verr = verr.Append("external_id", "is taken")

	t.Run("Append", func(t *testing.T) {
		got := verr
		exp := VErrors{
			{Field: "name", Description: "is required"},
			{Field: "name", Description: "must be 5 chars"},
			{Field: "external_id", Description: "is taken"},
		}
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("Messages", func(t *testing.T) {
		got := verr.Messages()
		sort.Strings(got)

		exp := []string{
			"external_id: is taken",
			"name: is required",
			"name: must be 5 chars",
		}
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("Status", func(t *testing.T) {
		sts := verr.Status(codes.InvalidArgument, "custom")
		got := sts.Err().Error()

		if exp := `rpc error: code = InvalidArgument desc = custom`; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("VErrorsFromStatus", func(t *testing.T) {
		sts := verr.Status(codes.InvalidArgument, "custom")
		got := VErrorsFromStatus(sts).Messages()

		if exp := verr.Messages(); !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("VErrorsConvert", func(t *testing.T) {
		sts := verr.Status(codes.InvalidArgument, "custom")
		got := VErrorsConvert(sts.Err()).Messages()

		if exp := verr.Messages(); !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})
}
