build:
	docker-compose build

start:
	docker-compose up -d
	sleep 1
	docker-compose ps

stop:
	docker-compose stop