default: vet test

test:
	go test ./...

vet:
	go vet ./...

release:
	goreleaser --rm-dist

.PHONY: test vet release
