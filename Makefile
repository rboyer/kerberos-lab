.PHONY: all
all: docker up

.PHONY: docker
docker:
	docker build -t local/kerb-common -f ./Dockerfile-common .
	docker build -t local/kerb-client -f ./Dockerfile-client .
	docker build -t local/kerb-server -f ./Dockerfile-server .

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down -v --remove-orphans
