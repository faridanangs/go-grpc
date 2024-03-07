protoc:
	protoc --proto_path=protos protos/*.proto --go_out=. --go-grpc_out=.
	@echo "protoc compile selesai"