GOCMD=go
GOBUILD=$(GOCMD) build

.PHONY: list build build-debug

list: 
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| column -t  -s ' '
build: ## Build
		$(GOBUILD)
build-debug: ## Build for debugging
		$(GOBUILD) -gcflags=all="-N -l"
