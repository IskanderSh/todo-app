build:
	docker-compose build todo-app

run:
	go run cmd/main.go

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up