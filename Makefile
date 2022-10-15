.PHONY: docker-build cover test test-integration

APP_NAME=canvas
DOCKER_USERNAME=karanbirsingh
GIT_HASH ?= $(shell git log --format="%h" -n 1)

docker-build:
	docker build --tag ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH} .

docker-push:
	docker push ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH}

docker-release:
	docker tag ${DOCKER_USERNAME}/${APP_NAME}:${GIT_HASH} ${DOCKER_USERNAME}/${APP_NAME}:latest
	docker push ${DOCKER_USERNAME}/${APP_NAME}:latest

cover:
	go tool cover -html=cover.out

run:
	go run -race cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...
