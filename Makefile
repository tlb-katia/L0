# Установите имя образа и контейнера
PUBLISHER_IMAGE = my-publisher
PUBLISHER_CONTAINER = publisher-container
COMPOSE_FILE = docker-compose.yml
CONTAINER_NAME = app

# Команды для сборки и запуска
.PHONY: build-publisher run-publisher up down

# Сборка образа паблишера
build-publisher:
	docker build -t $(PUBLISHER_IMAGE) -f ./server/order/nats_prod/Dockerfile .

# Запуск контейнера паблишера
run-publisher:
	docker run --rm --name publisher-container --network l0_my-network my-publisher

# Запуск Docker Compose
up:
	docker-compose up --build $(CONTAINER_NAME)

# Остановка и удаление контейнеров, сетей, томов
down:
	docker-compose -f $(COMPOSE_FILE) down