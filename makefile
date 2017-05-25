gen-pb:
	@bash -c "protoc --proto_path=./vendor -I ./servers/proto_rpc/ ./servers/proto_rpc/server.proto --gogo_out=plugins=grpc:./servers/proto_rpc"