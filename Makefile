SHELL=/bin/bash
GITROOT=$(shell git rev-parse --show-toplevel)
GITHUB_TOKEN ?= "SPECIFY_GITHUB_TOKEN_IN_ENVIRONMENT"

.DEFAULT_GOAL := help

.PHONY: help
help: ## show make targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install
install: ## install go modules
	@go version
	@go mod tidy
	@go mod verify

.PHONY: build
build: ## build go binary
	@tag="$(shell git tag -l --points-at HEAD)"; go build -ldflags="-X 'main.version=$$tag'" -o macgyver main.go
	@tar zcvf macgyver.tar.gz macgyver

# Ubuntu Only
.PHONY: setup
setup: ## install tools
	@sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-key C99B11DEB97541F0
	@sudo apt-add-repository https://cli.github.com/packages
	@sudo apt update && sudo apt install -y gh

# Note: You still need to do the final check or edit this release manually in the GitHub (gh-cli is required)
.PHONY: release
release: ## create release package
	@tag="$(shell git tag -l --points-at HEAD)"; gh release create $$tag -t $$tag --draft --prerelease macgyver.tar.gz --generate-notes
	
.PHONY: clean
clean: ## clean build files
	@rm -f macgyver macgyver.tar.gz
