APP_NAME  = khcheck-aws-iam-role
VERSION  ?= $(shell awk -F= '$$1 ~ /ARG VERSION/ {print $$2}' Dockerfile)
DATE     ?= $(shell date +%FT%T%z)
CUR_DIR   = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.DEFAULT_GOAL := help

all:

help:
			@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

arm64: ## Build ARM64 linux container image
			docker build --build-arg TARGET_ARCH=arm64 -f Dockerfile -t $(APP_NAME):$(VERSION) -t $(APP_NAME):$(VERSION) $(CUR_DIR)

amd64: ## Build AMD64 linux container image
			docker build --build-arg TARGET_ARCH=amd64 -f Dockerfile -t $(APP_NAME):$(VERSION) -t $(APP_NAME):$(VERSION) $(CUR_DIR)

test:
			@echo $(VERSION)
.PHONY: all image
