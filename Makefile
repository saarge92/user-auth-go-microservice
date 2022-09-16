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

.PHONY: test
test:
	$(eval TEST_CONTAINER_ID=$(shell docker run --rm -d --env-file .env.test -p 3312:3306 percona/percona-server:latest))
	while [ -n `docker logs ${TEST_CONTAINER_ID} 2>&1 | grep 'ready' | grep '3306'` ]; do sleep 1; done
	go run scripts/util/migration.go --env-file=.env.test || (docker stop ${TEST_CONTAINER_ID}; exit 1)
	go generate -v ./...
	go test -v ./internal/... || (docker stop ${TEST_CONTAINER_ID}; exit 1)
	docker stop ${TEST_CONTAINER_ID}