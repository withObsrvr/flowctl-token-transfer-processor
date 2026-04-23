BINARY=token-transfer-processor

.PHONY: proto build test tidy

proto:
	protoc \
	  -I . \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  proto/stellar/v1/token_transfers.proto

build:
	go build -o bin/$(BINARY) ./cmd/$(BINARY)

test:
	go test ./...

tidy:
	go mod tidy
