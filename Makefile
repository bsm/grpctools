default: test

test:
	go test ./...

release:
	goreleaser --rm-dist

staticcheck:
	staticcheck ./...

.PHONY: test  release
