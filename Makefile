GOCMD=go
GOBUILD=$(GOCMD) build
GIT_SHA := $(shell git describe --always --long --dirty)

.PHONY: list build build-debug

ifndef VERSION
$(error VERSION env variable not set)
endif

list: 
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| column -t  -s ' '
build: ## Build
		$(GOBUILD) -ldflags="-X 'github.com/softpuff/s3commander/cmd.Version=${VERSION}-${GIT_SHA}'"
build-linux:
		GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-X 'github.com/softpuff/s3commander/cmd.Version=${VERSION}-${GIT_SHA}'"
build-debug: ## Build for debugging
		$(GOBUILD) -gcflags=all="-N -l"
