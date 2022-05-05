up:
	docker-compose up
	migrate-up
	go run scripts/fixtures/test.go
	go run cmd/main.go
init:
	rm -rf pkg/protobuf/*
	find -L api/proto/ -name "*.proto" -exec protoc -I api/proto --go-grpc_out=require_unimplemented_servers=false:pkg --go_out=pkg {} +

lint:
	golangci-lint run ./internal/...

migrate-up:
	migrate -path scripts/migrations -database "mysql://user:pass@(localhost:3311)/user-platform?charset=utf8&parseTime=true" \
 	-verbose up

migrate-down:
	migrate -path scripts/migrations -database "mysql://user:pass@(localhost:3311)/user-platform?charset=utf8&parseTime=true" \
	-verbose down

migrate-test-up:
	migrate -path scripts/migrations -database "mysql://test:test@(localhost:3312)/user-test?charset=utf8&parseTime=true" \
     	-verbose up

migrate-test-down:
	migrate -path scripts/migrations -database "mysql://test:test@(localhost:3312)/user-test?charset=utf8&parseTime=true" \
      -verbose down

.PHONY: test
test:
	go test -v ./test/...