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



.PHONY: deploy-lambda
deploy-lambda:
	@echo "🔨 Building lambda binary…"
	@GOOS=linux go build -o bin/${LAMBDA_FUNCTION_NAME} services/s3UploadListener/main.go

	@echo "📦 Zipping…"
	@zip -j bin/${LAMBDA_FUNCTION_NAME}.zip bin/${LAMBDA_FUNCTION_NAME}

	@echo "🚀 Ensuring Lambda exists…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    lambda get-function \
	      --function-name ${LAMBDA_FUNCTION_NAME} \
	  > /dev/null 2>&1 \
	  || ( \
	       echo "🆕 Creating Lambda function…"; \
	       awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	         lambda create-function \
	           --function-name ${LAMBDA_FUNCTION_NAME} \
	           --runtime go1.x \
	           --handler ${LAMBDA_FUNCTION_NAME} \
	           --role arn:aws:iam::000000000000:role/${LAMBDA_FUNCTION_NAME}-exec \
	           --zip-file fileb://bin/${LAMBDA_FUNCTION_NAME}.zip \
	     )
	@echo "⏳ Waiting for function to become Active…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    lambda wait function-active-v2 \
	      --function-name ${LAMBDA_FUNCTION_NAME}

	@echo "🚀 Updating function code…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    lambda update-function-code \
	      --function-name ${LAMBDA_FUNCTION_NAME} \
	      --zip-file fileb://bin/${LAMBDA_FUNCTION_NAME}.zip
		  

	@echo "⚙️  Updating env var API_ENDPOINT…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    lambda update-function-configuration \
	      --function-name ${LAMBDA_FUNCTION_NAME} \
		--timeout 120 \
	      --environment "Variables={\
		KAFKA_BROKERS=${KAFKA_BROKERS},\
		API_ENDPOINT=${API_ENDPOINT},\
		KAFKA_TOPIC_FILE_UPLOAD=${KAFKA_TOPIC_FILE_UPLOAD}\
		}"
		



.PHONY: enable-s3-notifications
enable-s3-notifications:

	@echo "🔑 Granting S3 permission to invoke Lambda…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    lambda add-permission \
	      --function-name ${LAMBDA_FUNCTION_NAME} \
	      --statement-id s3invoke-${AWS_S3_BUCKET_NAME} \
	      --principal s3.amazonaws.com \
	      --action lambda:InvokeFunction \
	      --source-arn arn:aws:s3:::${AWS_S3_BUCKET_NAME} \
	  || true

	@echo "🔔 Wiring S3 → Lambda notification…"
	@awslocal --endpoint-url=${AWS_BASE_ENDPOINT} --region=${AWS_DEFAULT_REGION} \
	    s3api put-bucket-notification-configuration \
	      --bucket ${AWS_S3_BUCKET_NAME} \
	      --notification-configuration '{\
"LambdaFunctionConfigurations":[{\
  "LambdaFunctionArn":"arn:aws:lambda:'${AWS_DEFAULT_REGION}':000000000000:function:'${LAMBDA_FUNCTION_NAME}'",\
  "Events":["s3:ObjectCreated:Put"]\
}]\
}'
