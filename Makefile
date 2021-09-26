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
	helm upgrade --install ${org}-org-mongodb \
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

fabric-init:
	kubectl create namespace network || echo "Namespace 'network' already exists"
	kubectl config set-context --current --namespace=network
	fabnctl gen artifacts -a=arm64 -d=chainmetric.network -f ./network-config.yaml \
    	--charts=./deploy/charts

fabric-install:
	fabnctl install orderer -a=arm64 -d=chainmetric.network \
		--charts=./deploy/charts

	kubectl create secret generic couchdb-auth \
		--from-literal=user=${COUCHDB_USERNAME} --from-literal=password=${COUCHDB_PASSWORD} \
		--dry-run=client -o yaml | kubectl apply -f -

	fabnctl install peer -a=arm64 -d=chainmetric.network -o chipa-inu -p peer0 \
		--charts=./deploy/charts
	fabnctl install peer -a=arm64 -d=chainmetric.network -o blueberry-go -p peer0 \
		--charts=./deploy/charts
	fabnctl install peer -a=arm64 -d=chainmetric.network -o moon-lan -p peer0 \
		--charts=./deploy/charts

	fabnctl install channel -a=arm64 -d=chainmetric.network -c=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0

	fabnctl install cc assets -a=arm64 -d=chainmetric.network \
		-C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		--image=chainmetric/assets-contract \
		--source=./smartcontracts/assets \
		--charts=./deploy/charts

	fabnctl install cc devices -a=arm64 -d=chainmetric.network \
		-C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		--image=chainmetric/devices-contract \
		--source=./smartcontracts/devices \
		--charts=./deploy/charts

	fabnctl install cc requirements -a=arm64 -d=chainmetric.network \
		-C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		--image=chainmetric/requirements-contract \
		--source=./smartcontracts/requirements \
		--charts=./deploy/charts

	fabnctl install cc readings -a=arm64 -d=chainmetric.network \
		-C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		--image=chainmetric/readings-contract \
		--source=./smartcontracts/readings \
		--charts=./deploy/charts

	fabnctl update channel -a=arm64 -d=chainmetric.network --setAnchors -c=supply-channel \
			-o=chipa-inu \
			-o=blueberry-go \
			-o=moon-lan

fabric-clear:
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

build-chaincode:
	fabnctl build cc ${cc} . --ssh --host=${NODE_IP} -u=ubuntu \
			--target=smartcontracts/${cc} --ignore="bazel-*" --push

install-chaincode:
	fabnctl install cc ${cc} -a=arm64 -d=chainmetric.network \
		-C=supply-channel \
		-o=chipa-inu -p=peer0 \
		-o=blueberry-go -p=peer0 \
		-o=moon-lan -p=peer0 \
		--image=chainmetric/${cc}-contract \
		--source=./smartcontracts/${cc} \
		--charts=./deploy/charts \

deploy-chaincode: build-chaincode install-chaincode

install-orgservice:
	kubectl create -n network secret tls ${service}.${ORG}.org.${DOMAIN}-tls \
		--key=".data/certs/grpc/${service}/server.key" \
		--cert=".data/certs/grpc/${service}/server.crt" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic ${service}.${ORG}.org.${DOMAIN}-ca \
		--from-file=".data/certs/grpc/${service}/ca.crt" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic ${service}-${ORG}-org-fabric-connection \
		--from-file=".data/config/${service}/connection.yaml" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic jwt-keys \
		--from-file=".data/certs/jwt/jwt-cert.pem" \
		--from-file=".data/certs/jwt/jwt-key.pem" \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic vault-credentials \
		--from-literal=token=${VAULT_TOKEN} \
		--dry-run=client -o yaml | kubectl apply -f -

	helm upgrade --install ${service}-${org} deploy/charts/api-service -f deploy/config/orgservices/${service}.yaml

	kubectl create secret generic vault-credentials \
		--from-literal=token=${VAULT_TOKEN} \
		--dry-run=client -o yaml | kubectl apply -f -

	kubectl create secret generic privileges-config \
		--from-file=orgservices/shared/data/privileges.yaml \
		--dry-run=client -o yaml | kubectl apply -f -

	helm upgrade --install ${service}-${ORG} deploy/charts/api-service

build-orgservice:
	bazel run //:gazelle
	bazel build //orgservices/${service}
	bazel run //orgservices/${service}:multiacrh
	bazel run //orgservices/${service}:multiacrh-push

deploy-orgservice: build-orgservice install-orgservice

grpc-tls-gen:
	mkdir .data/certs/grpc/${service} || echo "directory .data/certs/grpc/${service} exists"
	openssl genrsa \
		-out .data/certs/grpc/${service}/ca.key 2048

	openssl req \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=${service}.${org}.org.${DOMAIN}" \
		-new -x509 -days 365 -key .data/certs/grpc/${service}/ca.key -out .data/certs/grpc/${service}/ca.crt

	openssl req -newkey rsa:2048 \
		-nodes -keyout .data/certs/grpc/${service}/server.key \
		-subj "/C=UA/ST=Kiev/O=Chainmetric, Inc./CN=${service}.${org}.org.${DOMAIN}" \
		-out .data/certs/grpc/${service}/server.csr

	openssl x509 -req \
		-in .data/certs/grpc/${service}/server.csr \
		-CA .data/certs/grpc/${service}/ca.crt -CAkey .data/certs/grpc/${service}/ca.key -CAcreateserial -days 365 \
		-extfile <(printf "subjectAltName=DNS:${service}.${org}.org.${DOMAIN},DNS:localhost,DNS:${service}-${org}-org") \
		-out .data/certs/grpc/${service}/server.crt

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
