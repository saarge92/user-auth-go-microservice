
init:
	rm -rf pkg/protobuf/*
	find -L api/proto/ -name "*.proto" -exec protoc -I api/proto --go-grpc_out=require_unimplemented_servers=false:pkg --go_out=pkg {} +

lint:
	golangci-lint run ./internal/...

migrateup:
	migrate -path scripts/migrations -database "mysql://user:pass@(localhost:3311)/user-platform?charset=utf8&parseTime=true" \
 	-verbose up

migratedown:
	migrate -path scripts/migrations -database "mysql://user:pass@(localhost:3311)/user-platform?charset=utf8&parseTime=true" \
	-verbose down

migrateup-test:
	migrate -path scripts/migrations -database "mysql://test:test@(localhost:3312)/user-test?charset=utf8&parseTime=true" \
     	-verbose up
