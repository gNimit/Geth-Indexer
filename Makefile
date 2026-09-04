.PHONY: all build test bench clean docker-build docker-up lint

BINARY_NAME=geth-indexer

all: test build

build:
	go build -v -o bin/$(BINARY_NAME) ./cmd/geth-indexer

test:
	go test -v -race -run=Test ./...

bench:
	go test -v -bench=. -benchmem ./...

docker-build:
	docker build -t $(BINARY_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

clean:
	rm -rf bin/ *.json *.tmp* coverage.out
