PKG=$(shell glide nv)
VERSION=0.1.2

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)

pkg-deps:
	go get github.com/mitchellh/gox

pkg: pkg-deps \
	pkg/release/grpc-health_${VERSION}_linux_amd64.zip

.PHONY: test vet errcheck deps pkg

# ---------------------------------------------------------------------

pkg/build/linux/amd64/grpc-health: cmd/grpc-health/main.go
	mkdir -p $(dir $@)
	gox -osarch="linux/amd64" -output $@ ./$(dir $<)

pkg/release/grpc-health_${VERSION}_linux_amd64.zip: pkg/build/linux/amd64/grpc-health
	mkdir -p $(dir $@)
	zip -9 -j $@ $<
