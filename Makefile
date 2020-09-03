APP_NAME              = khcheck-aws-iam-role
VERSION               = $(shell awk -F= '$$1 ~ /APP_VERSION/ {print $$2}' Dockerfile)
DATE                 ?= $(shell date +%FT%T%z)
CUR_DIR               = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.DEFAULT_GOAL := help

all:

help:
			@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

image: ## Build Local Container Image (amd64)
			docker build -f Dockerfile -t $(APP_NAME):$(VERSION) -t $(APP_NAME):$(VERSION) $(CUR_DIR)

.PHONY: all image
