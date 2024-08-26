include .env

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down

docs:
	swag init -g internal/app/app.go --pd
.PHONY: docs

keys:
	openssl genpkey -algorithm RSA -out private.key && \
	openssl rsa -pubout -in private.key -out public.key

pg-tests:
	docker run --name postgres --rm -d \
		-p 6000:6000 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=1234567890 \
		-e POSTGRES_DB=postgres postgres:15 -p 6000

init-test-containers: pg-tests

stop-test-containers:
	docker stop postgres

init-tests:
	go test -v ./...

tests: init-test-containers init-tests stop-test-containers
