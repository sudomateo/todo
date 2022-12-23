.PHONY: build
build:
	docker compose build

.PHONY: up
up: build
	docker compose up

.PHONY: down
down:
	docker compose down --volumes

.PHONY: test
test:
	go test -v ./... -count 1
