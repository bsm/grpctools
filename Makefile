SRC=$(shell find . -name '*.go')
PKG=./...

TARGET_PKG=$(patsubst cmd/%/main.go,bin/%,$(wildcard cmd/*/main.go))
TARGET_OS=linux darwin
TARGET_ARCH=amd64
TARGETS=$(foreach pkg,$(TARGET_PKG),$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(pkg)-$(os)-$(arch))))
ZIPS=$(foreach t,$(TARGETS),$(t).zip)
TGZS=$(foreach t,$(TARGETS),$(t).tgz)

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)

build: $(TARGETS)
dist: $(ZIPS) $(TGZS)

.PHONY: test vet errcheck deps build dist

# ---------------------------------------------------------------------

bin/grpc-health-%.zip: bin/grpc-health-%
	zip -j $@ $<

bin/grpc-health-%.tgz: bin/grpc-health-%
	tar -czf $@ -C $(dir $<) $(notdir $<)

bin/grpc-health-%: $(SRC)
	@mkdir -p $(dir $@)
	$(eval os := $(word 3, $(subst -, ,$@)))
	$(eval arch := $(word 4, $(subst -, ,$@)))
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -o $@ $(patsubst bin/%-$(os)-$(arch),cmd/%/main.go,$@)
