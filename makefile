PORT ?= 9191

run.server:
	go run cmd/server/main.go -port $(PORT)

test:
	go test -v ./...
