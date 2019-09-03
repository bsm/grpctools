package grpctools

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
)

var _ = Describe("VErrors", func() {
	var subject VErrors

	BeforeEach(func() {
		subject = subject.Reset()
		subject = subject.Append("name", "is required")
		subject = subject.Append("name", "must be 5 chars")
		subject = subject.Append("external_id", "is taken")
	})

	It("should add errors", func() {
		Expect(subject).To(Equal(VErrors{
			{Field: "name", Description: "is required"},
			{Field: "name", Description: "must be 5 chars"},
			{Field: "external_id", Description: "is taken"},
		}))
	})

	It("should build messages", func() {
		Expect(subject.Messages()).To(ConsistOf(
			"name: is required",
			"name: must be 5 chars",
			"external_id: is taken",
		))
	})

	It("should export status", func() {
		sts := subject.Status(codes.InvalidArgument, "custom")
		Expect(sts.Err()).To(MatchError(`rpc error: code = InvalidArgument desc = custom`))
	})

	It("should parse from status", func() {
		sts := subject.Status(codes.InvalidArgument, "custom")
		Expect(VErrorsFromStatus(sts).Messages()).To(Equal(subject.Messages()))
	})

	It("should convert from error", func() {
		sts := subject.Status(codes.InvalidArgument, "custom")
		Expect(VErrorsConvert(sts.Err()).Messages()).To(Equal(subject.Messages()))
	})
})
