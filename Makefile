.PHONY: docker-build cover test test-integration

APP_NAME=canvas
DOCKER_REGISTRY=docker.io
DOCKER_USERNAME=karanbirsingh
GIT_HASH ?= $(shell git log --format="%h" -n 1)
AZ_RESOURCE_GROUP=$(shell az group list | jq '.[].name')
AZ_ACI_APP_NAME=canvas-app

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

.PHONY: aci-deploy
aci-deploy:
	az container create \
		--resource-group ${AZ_RESOURCE_GROUP} \
		--name ${AZ_ACI_APP_NAME} \
		--image ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${APP_NAME} \
		--dns-name-label ${AZ_ACI_APP_NAME} \
		--ports 80 \
		--environment-variables 'PORT'='80'

.PHONY: aci-destroy
aci-destroy:
	az container delete \
		--resource-group ${AZ_RESOURCE_GROUP} \
		--name ${AZ_ACI_APP_NAME} -y

.PHONY: db-start
db-start:
	@docker container run --name postgres-canvas -d -p 5432:5432 -e POSTGRES_USER='canvas' -e POSTGRES_PASSWORD='postgres' -v /tmp/postgres:/var/lib/postgresql/data postgres:12

.PHONY: db-stop
db-stop:
	@docker container rm -f postgres-canvas || true

.PHONY: db-restart
db-restart: db-stop db-start