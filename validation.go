package grpctools

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VErrors contains a set of field-violations.
type VErrors []*errdetails.BadRequest_FieldViolation

// VErrorsConvert extracts validation errors from an error.
// This function will return nil if no validation errors are attached.
func VErrorsConvert(err error) VErrors {
	return VErrorsFromStatus(status.Convert(err))
}

// VErrorsFromStatus extracts validation errors from status.
func VErrorsFromStatus(sts *status.Status) VErrors {
	for _, detail := range sts.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			return t.GetFieldViolations()
		}
	}
	return nil
}

// Append appends a field message.
func (e VErrors) Append(field, message string) VErrors {
	return append(e, &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: message,
	})
}

// Len returns the error count.
func (e VErrors) Len() int {
	return len(e)
}

// Reset resets the slice.
func (e VErrors) Reset() VErrors {
	return e[:0]
}

// Messages returns messages.
func (e VErrors) Messages() []string {
	msgs := make([]string, 0, e.Len())
	for _, fv := range e {
		msgs = append(msgs, fv.Field+": "+fv.Description)
	}
	return msgs
}

// Status returns a custom status.
func (e VErrors) Status(code codes.Code, message string) *status.Status {
	// this should not error, if it does it's better panic here to instantly figure out why
	sts, err := status.New(code, message).
		WithDetails(&errdetails.BadRequest{FieldViolations: e})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return sts
}

// InvalidArgument returns an InvalidArgument status.
func (e VErrors) InvalidArgument(message string) *status.Status {
	if message == "" {
		message = "invalid argument"
	}
	return e.Status(codes.InvalidArgument, message)
}
