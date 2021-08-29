
init:
	rm -rf pkg/protobuf/*
	find -L api/proto/ -name "*.proto" -exec protoc -I api/proto --go-grpc_out=require_unimplemented_servers=false:pkg --go_out=pkg {} +