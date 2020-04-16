.PHONY: build

BINARY_NAME := tldr
VERSION ?= $(shell cat VERSION)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

build:
	go build -o bin/$(BINARY_NAME) cli/tldr/*.go

clean:
	rm bin/$(BINARY_NAME)*

run:
	go run cli/tldr/*.go

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build -o bin/$(BINARY_NAME)-$(VERSION)-$(os)-amd64 cli/tldr/*.go

.PHONY: release
release: windows linux darwin

