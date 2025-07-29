run-auth:
	@cd services/auth && go run cmd/auth-server/*.go 


gen:
	@protoc \
		--proto_path=api/proto \
		--go_out=paths=source_relative:services/common/genproto/auth \
		--go-grpc_out=paths=source_relative:services/common/genproto/auth \
		api/proto/auth.proto