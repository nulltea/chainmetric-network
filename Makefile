swagger:
	swag init -g ./src/users/api/doc.go -o ./src/users/docs

deploy-identity:
	kubectl create secret generic identity-chipa-inu-org-hlf-connection \
	 --from-file=connection.yaml \
	  --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install identity-chipa-inu deploy/charts/api-service

docker-build:
	sudo docker buildx build \
		--platform linux/arm64 -t chainmetric/api.identity \
		-f ./deploy/docker/users.Dockerfile --push .
