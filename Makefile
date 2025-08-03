run-auth:
	@cd services/auth && go run cmd/auth-server/*.go 

run-fileIngestion:
	@cd services/fileIngestion && go run cmd/*.go 

run-worker:
	@cd services/worker && go run cmd/*.go 

run-grpc-client:
	@cd client && go run *.go 

gen:
	@protoc \
		--proto_path=api/proto \
		--go_out=paths=source_relative:services/common/genproto/auth \
		--go-grpc_out=paths=source_relative:services/common/genproto/auth \
		api/proto/auth.proto
	@protoc \
		--proto_path=api/proto \
		--go_out=paths=source_relative:services/common/genproto/fileIngestion \
		--go-grpc_out=paths=source_relative:services/common/genproto/fileIngestion \
		api/proto/fileIngestion.proto
	@protoc \
		--proto_path=api/proto \
		--go_out=paths=source_relative:services/common/genproto/fileUpload \
		--go-grpc_out=paths=source_relative:services/common/genproto/fileUpload \
		api/proto/fileUpload.proto