
test:
	go test ./...

mock:
	go generate ./...

run: 
	POSTGRES_PASSWORD=postgres POSTGRES_DB=example \
	 HTTP_ADDRESS=localhost:8000 go run cmd/example/main.go

migrate-up:
	POSTGRES_PASSWORD=postgres POSTGRES_DB=example \
	 go run cmd/migrations/main.go -m up