proto:
	go generate ./internal/generator/api/rpc/...

run:
	go run cmd/main.go
