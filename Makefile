.PHONY: docker-build cover test test-integration

APP_NAME=canvas
DOCKER_USERNAME=karanbirsingh
GIT_HASH ?= $(shell git log --format="%h" -n 1)

.PHONY: docker-build
docker-build:
	docker build --tag ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH} .

.PHONY: docker-push
docker-push:
	docker push ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH}

.PHONY: docker-release
docker-release:
	docker tag ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH} ${DOCKER_USERNAME}/${APP_NAME}:latest
	docker push ${DOCKER_USERNAME}/${APP_NAME}:latest

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: run
run:
	go run -race cmd/server/*.go

.PHONY: test
test:
	go test -coverprofile=cover.out -short ./...

.PHONY: test-integration
test-integration:
	go test -coverprofile=cover.out -p 1 ./...
