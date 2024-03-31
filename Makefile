client	:
	go run cmd/client/main.go

server	:
	go run cmd/server/main.go

up	:
	docker-compose -f docker-compose.yaml up -d
down	:
	docker-compose -f docker-compose.yaml down
