.PHONY: up clean down lint test

up:
	docker-compose up -d --build

down:
	docker-compose down

clean:
	docker-compose down -v

lint:
	cd order-service && golangci-lint run ./...

test:
	cd order-service && go test ./...

