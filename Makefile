default: test

test:
	go test ./...

release:
	goreleaser --rm-dist

lint:
	golangci-lint run

.PHONY: test lint release
