APP_SWAGGER_DIR=../app/app/assets/swagger
DOMAIN=chainmetric.network
ORG=chipa-inu

swagger:
	swag init -g ./src/users/api/doc.go -o ./src/users/docs -o ${APP_SWAGGER_DIR}
	rm ${APP_SWAGGER_DIR}/docs.go ${APP_SWAGGER_DIR}/swagger.yaml

deploy-identity:
	kubectl create -n network secret tls identity.${ORG}.org.${DOMAIN}-tls \
		--key="data/certs/grpc/server.key" \
		--cert="data/certs/grpc/server.crt" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic identity.${ORG}.org.${DOMAIN}-ca \
		--from-file="data/certs/grpc/ca.crt" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic identity-${ORG}-org-hlf-connection \
		--from-file=connection.yaml \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic identity-${ORG}-org-jwt-keys \
		--from-file=data/certs/jwt/jwt-cert.pem \
		--from-file=data/certs/jwt/jwt-key.pem \
		--dry-run=client -o yaml | kubectl apply -f -

	helm upgrade --install identity-chipa-inu deploy/charts/api-service

docker-build:
	sudo docker buildx build \
		--platform linux/arm64 -t chainmetric/api.identity \
		-f ./deploy/docker/identity.Dockerfile --push .

deploy-build: docker-build deploy-identity

grpc-gen:
	protoc \
		-I=src \
		-I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
		-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
	    --go_out=paths=source_relative:src \
	    --validate_out=lang=go,paths=source_relative:src \
		./src/identity/api/presenter/*.proto

	protoc \
		-I=src \
		-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
		--go-grpc_out=paths=source_relative:./src \
		./src/identity/api/rpc/*.proto

grpc-ui:
	grpcui \
 		--open-browser \
 		-cert ./data/certs/grpc/server.crt \
 		-key ./data/certs/grpc/server.key \
 		-import-path ./src \
 		-import-path ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
 		-proto ./src/identity/api/rpc/identity.proto \
 		identity.chipa-inu.org.chainmetric.network:443

grpc-tls-gen:
	openssl genrsa \
		-out data/certs/ca.key 2048

	openssl req \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=identity.${ORG}.org.${DOMAIN}" \
		-new -x509 -days 365 -key data/certs/ca.key -out data/certs/ca.crt

	openssl req -newkey rsa:2048 \
		-nodes -keyout data/certs/server.key \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=identity.${ORG}.org.${DOMAIN}" \
		-out data/certs/server.csr

	openssl x509 -req \
		-in data/certs/server.csr \
		-CA data/certs/ca.crt -CAkey data/certs/ca.key -CAcreateserial -days 365 \
		-extfile <(printf "subjectAltName=DNS:identity.${ORG}.org.${DOMAIN},DNS:localhost,DNS:identity-${ORG}-org") \
		-out data/certs/server.crt


cp-proto-app:
	cp ./src/users/api/presenter/users.proto ../app/app/assets/proto/user.proto
	cp ./src/users/api/rpc/identity.proto ../app/app/assets/proto/identity_grpc.proto
