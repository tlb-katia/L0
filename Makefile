COMPOSE_FILE = docker-compose.yml
CONTAINER_NAME = app
PROJECT_DIR = /home/katia/GolandProjects/L0
NATS_DIR = $(PROJECT_DIR)/internal/server/order/nats_prod

.PHONY: publisher

all: up

up: permission
	docker-compose up --build $(CONTAINER_NAME)

publisher:
	cd $(NATS_DIR) && go run main.go


clean-data: permission
	rm -rf pkg/repository/db/pgdata
	rm -rf pkg/storage/redis/data

down:
	docker-compose -f $(COMPOSE_FILE) down

server-logs:
	docker logs $(CONTAINER_NAME)


permission:
	sudo chmod -R 755 $(PROJECT_DIR)/pkg/repository/db/pgdata
