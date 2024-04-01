client	:
	go run cmd/client/main.go

server	:
	go run cmd/server/main.go

test	:
	go test -v ./internal/*

volume	:
	docker volume create redis_data

up	:
	docker-compose -f docker-compose.yaml up -d

down	:
	docker-compose -f docker-compose.yaml down
