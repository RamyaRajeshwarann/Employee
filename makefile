run: # Run the server 
	cd cmd && go run . app:serve

swagger: # Generate swagger
	swag init -g cmd/server.go

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable" up

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable" down

test:
	go test ./... -v

.PHONY: run