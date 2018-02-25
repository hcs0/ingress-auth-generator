.DEFAULT_GOAL := help
.PHONY: help

OS := $(shell uname -s)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

APP_NAME = htpasswd-gen
REGISTRY = banzaicloud

build: ## Builds binary package
	go build .

build-linux: clean
	CGO_ENABLED=0 GOOS=linux go build .

build-container: build-linux
#	TAG=0.1.${CIRCLE_BUILD_NUM}
	docker build -t ${REGISTRY}/${APP_NAME}:0.1 --build-arg APP_NAME=${APP_NAME} .
	docker push ${REGISTRY}/${APP_NAME}:0.1


