SHELL := /bin/bash

.PHONY: all
all: docker up

.PHONY: gomod
gomod:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

.PHONY: static
static:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w'

.PHONY: docker
docker: static
	docker build -t local/kerb-common -f ./Dockerfile-common .
	docker build -t local/kerb-server -f ./Dockerfile-server .
	docker build -t local/kerb-client -f ./Dockerfile-client .
	docker build -t local/kerb-demo -f ./Dockerfile-demo .

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down -v --remove-orphans

.PHONY: demo
demo: docker
	docker run --rm -it --net=container:kdc local/kerb-demo
