include .env
export

k3s:
	k3sup install \
		--ip=${NODE_IP} \
		--user=${NODE_USERNAME} \
		--k3s-extra-args '--no-deploy traefik'
	kubectl config use-context rpi-${CLUSTER_NAME}-k3s

storage:
	kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml

traefik:
	kubectl create secret generic do-auth-token \
		--from-literal=token=${DO_AUTH_TOKEN} --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install traefik -n=kube-system -f=charts/helm-values/traefik.yaml \
		--set=ports.websecure.tls.domains[0].main=${DOMAIN},ports.websecure.tls.domains[0].sans[0]=*.${DOMAIN} \
		traefik/traefik

infrastructure:
	helm upgrade --install -n=kube-system --set=ingress.routes.host=proxy.${DOMAIN} \
		infrastructure charts/infrastructure

mongodb:
	helm upgrade --install chipa-inu-org-mongodb \
		--set auth.password=${MONGO_PASSWORD} -f ./charts/helm-values/mongodb.yaml bitnami/mongodb
	kubectl create secret generic mongo-auth \
 		--from-literal=username=${MONGO_USERNAME} --from-literal=password=${MONGO_PASSWORD} \
 		--dry-run=client -o yaml | kubectl apply -f -

 vault:
	helm upgrade --install vault -n=kube-system -f=charts/helm-values/vault.yaml \
		hashicorp/vault
	kubectl -n=kube-system exec vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > .secrets/cluster-keys.json
	kubectl -n=kube-system exec vault-0 -- vault operator unseal $(cat .secrets/cluster-keys.json | jq -r ".unseal_keys_b64[]")
	kubectl -n=kube-system exec vault-0 -- vault login $(cat .secrets/cluster-keys.json | jq -r ".root_token")
	kubectl -n=kube-system exec vault-0 -- vault secrets enable -path=/fabric/identity kv-v1
	kubectl -n=kube-system exec vault-0 -- vault auth enable userpass

hyperledger-init:
	kubectl create namespace network || echo "Namespace 'network' already exists"
	kubectl config set-context --current --namespace=network
	fabnctl gen artifacts -a=arm64 -d=chainmetric.network -f ./network-config.yaml

hyperledger-deploy:
	fabnctl deploy orderer -a=arm64 -d=chainmetric.network

	kubectl create secret generic couchdb-auth \
		--from-literal=user=${COUCHDB_USERNAME} --from-literal=password=${COUCHDB_PASSWORD} \
		--dry-run=client -o yaml | kubectl apply -f -

	fabnctl deploy peer -a=arm64 -d=chainmetric.network -o chipa-inu -p peer0
	fabnctl deploy peer -a=arm64 -d=chainmetric.network -o blueberry-go -p peer0
	fabnctl deploy peer -a=arm64 -d=chainmetric.network -o moon-lan -p peer0

	fabnctl deploy channel -a=arm64 -d=chainmetric.network -c=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0

	fabnctl deploy cc -a=arm64 -d=chainmetric.network -c assets -C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		-r iotchainnetwork -f deploy/docker/assets.Dockerfile \
		--rebuild=false \
		../contracts

	fabnctl deploy cc -a=arm64 -d=chainmetric.network -c devices -C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		-r iotchainnetwork -f deploy/docker/devices.Dockerfile \
		--rebuild=false \
		../contracts

	fabnctl deploy cc -a=arm64 -d=chainmetric.network -c requirements -C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		-r iotchainnetwork -f deploy/docker/requirements.Dockerfile \
		--rebuild=false \
		../contracts

	fabnctl deploy cc -a=arm64 -d=chainmetric.network -c readings -C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		-r iotchainnetwork -f deploy/docker/readings.Dockerfile \
		--rebuild=false \
		../contracts

	fabnctl update channel -a=arm64 -d=chainmetric.network --setAnchors -c=supply-channel \
			-o=chipa-inu \
			-o=blueberry-go \
			-o=moon-lan

hyperledger-clear:
	helm uninstall peer0-chipa-inu || echo "Chart 'peer/peer0-chipa-inu' already uninstalled"
	helm uninstall peer0-blueberry-go || echo "Chart 'peer/peer0-blueberry-go' already uninstalled"
	helm uninstall peer0-moon-lan || echo "Chart 'peer/peer0-moon-lan' already uninstalled"

	helm uninstall assets-cc-peer0-chipa-inu || echo "Chart 'peer/assets-cc-peer0-chipa-inu' already uninstalled"
	helm uninstall assets-cc-peer0-blueberry-go || echo "Chart 'peer/assets-cc-peer0-blueberry-go' already uninstalled"
	helm uninstall assets-cc-peer0-moon-lan || echo "Chart 'peer/assets-cc-peer0-moon-lan' already uninstalled"

	helm uninstall devices-cc-peer0-chipa-inu || echo "Chart 'peer/devices-cc-peer0-chipa-inu' already uninstalled"
	helm uninstall devices-cc-peer0-blueberry-go || echo "Chart 'peer/devices-cc-peer0-blueberry-go' already uninstalled"
	helm uninstall devices-cc-peer0-moon-lan || echo "Chart 'peer/devices-cc-peer0-moon-lan' already uninstalled"

	helm uninstall requirements-cc-peer0-chipa-inu || echo "Chart 'peer/requirements-cc-peer0-chipa-inu' already uninstalled"
	helm uninstall requirements-cc-peer0-blueberry-go || echo "Chart 'peer/requirements-cc-peer0-blueberry-go' already uninstalled"
	helm uninstall requirements-cc-peer0-moon-lan || echo "Chart 'peer/requirements-cc-peer0-chipa-inu' already uninstalled"

	helm uninstall readings-cc-peer0-chipa-inu || echo "Chart 'peer/readings-cc-peer0-chipa-inu' already uninstalled"
	helm uninstall readings-cc-peer0-blueberry-go || echo "Chart 'peer/readings-cc-peer0-blueberry-go' already uninstalled"
	helm uninstall readings-cc-peer0-moon-lan || echo "Chart 'peer/readings-cc-peer0-moon-lan' already uninstalled"

	helm uninstall orderer || echo "Chart 'orderer/orderer' already uninstalled"

	# helm uninstall artifacts || echo "Chart 'artifacts/artifacts' already uninstalled"

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

	kubectl create secret generic vault-credentials \
		--from-literal=token=${VAULT_TOKEN} \
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
	    --go_out=paths=source_relative:orgservices \
	    --validate_out=lang=go,paths=source_relative:orgservices \
		./orgservices/identity/api/presenter/*.proto

	protoc \
		-I=src \
		-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
		--go-grpc_out=paths=source_relative:./orgservices \
		./orgservices/identity/api/rpc/*.proto

	ls ./orgservices/identity/api/rpc/*_grpc_grpc.pb.go | sed -E "p;s/(.*)_grpc_grpc\.pb\.go/\1_grpc\.pb.\go/" | xargs -n2 mv

grpc-ui:
	grpcui \
 		--open-browser \
 		-cert .data/certs/grpc/server.crt \
 		-key .data/certs/grpc/server.key \
 		-import-path ./src \
 		-import-path ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
 		-proto ./orgservices/identity/api/rpc/identity.proto \
 		-proto ./orgservices/identity/api/rpc/access.proto \
 		identity.chipa-inu.org.chainmetric.network:443

grpc-tls-gen:
	openssl genrsa \
		-out .data/certs/ca.key 2048

	openssl req \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=identity.${ORG}.org.${DOMAIN}" \
		-new -x509 -days 365 -key .data/certs/ca.key -out .data/certs/ca.crt

	openssl req -newkey rsa:2048 \
		-nodes -keyout .data/certs/server.key \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=identity.${ORG}.org.${DOMAIN}" \
		-out .data/certs/server.csr

	openssl x509 -req \
		-in .data/certs/server.csr \
		-CA .data/certs/ca.crt -CAkey .data/certs/ca.key -CAcreateserial -days 365 \
		-extfile <(printf "subjectAltName=DNS:identity.${ORG}.org.${DOMAIN},DNS:localhost,DNS:identity-${ORG}-org") \
		-out .data/certs/server.crt

cp-proto-app:
	cp ./orgservices/users/api/presenter/users.proto ../app/app/assets/proto/user.proto
	cp ./orgservices/users/api/rpc/identity.proto ../app/app/assets/proto/identity_grpc.proto
