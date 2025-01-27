server:
	go run cmd/server/main.go 

client:
	go run cmd/client/main.go

tester:
	go run cmd/tester/main.go

tests:
	go test ./...

build:
	docker compose build 

run:
	docker compose up