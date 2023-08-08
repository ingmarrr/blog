
build:
	go build -o build/main main.go

run:
	go run main.go

fmt:
	go fmt github.com/...

test:
	go test ./...

tw:
	npx tailwindcss -i ./style/input.css -o ./build/output.css --watch

surreal:
	docker compose -f surrealdb.Dockerfile up

