GOCMD=go
GOBUILD=$(GOCMD) build
VERSION := $(shell git describe --tag --always --long)

.PHONY: list build build-debug


list: 
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| column -t  -s ' '
build: ## Build
		$(GOBUILD) -ldflags="-X 'github.com/softpuff/s3commander/cmd.Version=${VERSION}'"
build-linux:
		GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-X 'github.com/softpuff/s3commander/cmd.Version=${VERSION}'"
build-debug: ## Build for debugging
		$(GOBUILD) -gcflags=all="-N -l"
