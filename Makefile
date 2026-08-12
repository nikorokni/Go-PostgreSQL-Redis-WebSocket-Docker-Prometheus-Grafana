.PHONY: run test lint compose-up compose-down demo
run:
	go run ./cmd/server
test:
	go test -race ./...
lint:
	go vet ./...
	test -z "$$(gofmt -l .)"
compose-up:
	docker compose up --build
compose-down:
	docker compose down
demo:
	curl -sS -X POST localhost:8080/v1/assess -H 'Content-Type: application/json' -d @examples/depeg.json
