build:
	docker-compose build

start:
	docker-compose up -d
	sleep 1
	docker-compose ps

up: build start
	@echo "All services are up"

stop:
	docker-compose stop