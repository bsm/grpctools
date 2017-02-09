SRC=$(shell find . -name '*.go' -not -path '*vendor*')
PKG=$(shell glide nv)

TARGET_PKG=$(patsubst cmd/%/main.go,bin/%,$(wildcard cmd/*/main.go))
TARGET_OS=linux darwin
TARGET_ARCH=amd64 386
TARGETS=$(foreach pkg,$(TARGET_PKG),$(foreach os,$(TARGET_OS),$(foreach arch,$(TARGET_ARCH),$(pkg)-$(os)-$(arch))))
ARCHIVES=$(foreach t,$(TARGETS),$(t).zip)

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)

build: $(TARGETS)
dist: $(ARCHIVES)

.PHONY: test vet errcheck deps build dist

# ---------------------------------------------------------------------

bin/grpc-health-%.zip: bin/grpc-health-%
	zip -j $@ $<

bin/grpc-health-%: $(SRC)
	@mkdir -p $(dir $@)
	$(eval os := $(word 3, $(subst -, ,$@)))
	$(eval arch := $(word 4, $(subst -, ,$@)))
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -o $@ $(patsubst bin/%-$(os)-$(arch),cmd/%/main.go,$@)
