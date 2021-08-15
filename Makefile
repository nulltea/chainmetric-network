APP_SWAGGER_DIR=../app/app/assets/swagger

swagger:
	swag init -g ./src/users/api/doc.go -o ./src/users/docs -o ${APP_SWAGGER_DIR}
	rm ${APP_SWAGGER_DIR}/docs.go ${APP_SWAGGER_DIR}/swagger.yaml

deploy-identity:
	kubectl create secret generic identity-chipa-inu-org-hlf-connection \
	 --from-file=connection.yaml \
	  --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install identity-chipa-inu deploy/charts/api-service

docker-build:
	sudo docker buildx build \
		--platform linux/arm64 -t chainmetric/api.identity \
		-f ./deploy/docker/users.Dockerfile --push .

grpc-gen:
	protoc \
		-I=./src/users/api/presenter \
		-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
	    --gogofast_out=plugins=grpc,paths=source_relative:./src/users/api/presenter \
	    --validate_out=lang=go,paths=source_relative:./src/users/api/presenter \
		./src/users/api/presenter/*.proto

	protoc \
		-I=./src/users/api/rpc \
		-I=./src/users/api/presenter \
		-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
		--gogofast_out=plugins=grpc,paths=source_relative:./src/users/api/rpc \
		./src/users/api/rpc/*.proto
