build:
	go build ./...
	go build ./cmd/...

test:
	go test -v ./...

fmt:
	go fmt ./...

install:
	go install ./cmd/...

release: build test

clean:
	go clean
	$(RM) godot_web

.PHONY: build test fmt install release clean
