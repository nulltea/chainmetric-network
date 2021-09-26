# Chainmetric: Smart Contracts

_**Chainmetric Smart Contracts**_ are designed to grand access to blockchain-stored data while implementing such data validation, aggregation, and management functionality.

Being a part of a permissioned blockchain network based on Hyperledger Fabric stack,
such Contracts are written in Go and configured to be deployed as Kubernetes services communicating with each other, blockchain peers, and external applications via gRPC protocol and event streaming.

## Requirements

- Deployed Chainmetric network ([see network deployment procedure][network deployment section]) preferably in Kubernetes environment
- Image registry where chaincode docker images will be hosted (Docker Hub account will do)
- [`fabnctl`][fabnctl] command line utility installed on local machine

## Contracts reference

### Assets contract

Exposes transactions for registering assets on Blockchain ledger, log changes of their state, transfer ownership, etc.

See details at: [**smartcontracts/assets**](https://github.com/timoth-y/chainmetric-network/main/smartcontracts/assets).

### Devices contract

Exposes transactions for registering IoT devices on chain, issue remote commands, transfer ownership, etc.

See details at: [**smartcontracts/devices**](https://github.com/timoth-y/chainmetric-network/main/smartcontracts/devices).

### Requirements contact

Exposes transactions for assigning, changing and revoking environmental requirements for storing and transferring assets.

See details at: [**smartcontracts/requirements**](https://github.com/timoth-y/chainmetric-network/main/smartcontracts/requirements).

### Readings contact

Exposes transactions for posting metric readings sourced by sensors of IoT devices.
Allows creating socket connection for receiving updates in real time.
Validates received readings against requirements and emits events on ones violations.

See details at: [**smartcontracts/readings**](https://github.com/timoth-y/chainmetric-network/main/smartcontracts/readings).

## Deployment

[Chaincodes][chaincode] (alternative to Smart Contracts) in Hyperledger Fabric infrastructure can be deployed both by embedding their source code into the blockchain peers and by deploying them [as external services][chaincode as external service], which is a way more versatile option especially for Kubernetes cluster environment where such Chaincodes can be deployed as a pods.

### Building from source

For chaincodes building or further updates use `install cc` command as following:

```bash
# Build over SSH
$ fabnctl build cc assets . \
   --ssh --host=${NODE_IP} --user=node \
   --target=smartcontracts/assets --ignore="bazel-*" --push

# Local Docker build
$ fabnctl build cc assets . \
   --dockerfile=./deploy/docker/assets.Deockerfile \
   --target=chainmetric-assets --push
```

Or use make rule `build-chaincode` specifying both name of the contract as following:

```bash
make cc=assets build-chaincode
```

### Installation

For chaincodes installation or further updates use `install cc` command as following: 

```bash
fabnctl install cc assets --arch=arm64 --domain=chainmetric.network -C=supply-channel \
   -o=org1 -p=peer0 \
   -o=org2 -p=peer0 \
   -o=org3 -p=peer0 \
   --image=chainmetric/assets-contract \
   --source=./smartcontracts/assets \
   --charts=./deploy/charts
```

Or use make rule `install-chaincode` specifying both name of the contract as following:

```bash
make cc=assets install-chaincode
```

Rule `deploy-chaincode` is a combination of both previous ones:
```bash
make service=identity deploy-chaincode
```

For more detailed instructions please refer to `fabnctl` [documentation](https://github.com/timoth-y/fabnctl#deploy-chaincodes).

## Development
For initializing local development environment use `bazel run` command specified gazelle plugin target.

```bash
bazel run //:gazelle
```

## Roadmap

- [x] Device remote commands over blockchain events [(#1)](https://github.com/timoth-y/chainmetric-network/pull/1)
- [x] Cache layer for storing contracts operational data e.g. `EventSocketSubscriptionTicket` [(#4)](https://github.com/timoth-y/chainmetric-network/pull/4)
- [x] Devices location management business logic [(#2)](https://github.com/timoth-y/chainmetric-network/pull/2)
- [x] Requirements violation notification
- [ ] Requirements violations rule engine

## Wrap up

Chainmetric's Smart Contracts are the accumulation of this project's business logic in the form of distributed, atomically granulated on-chain services, where each is responsible only for its entities, use cases, and transactions handling.

They are exposing access to the data stored on [blockchain immutable ledger][chainmetric network repo] by publishing their contracts as gRPC remote procedures and continuously emitted event streams.

Other parts of the Chainmetric project, such as embedded [sensor-equipped IoT][chainmetric sensorsys repo] devices and cross-platform [mobile application][chainmetric app repo] utilizes exposed contacts by remote procedure calling and subscribing to event streams.

## License

Licensed under the [Apache 2.0][license file].



[golang badge]: https://img.shields.io/badge/Code-Golang-informational?style=flat&logo=go&logoColor=white&color=6AD7E5
[lines counter]: https://img.shields.io/tokei/lines/github/timoth-y/chainmetric-contracts?color=teal&label=Lines
[commit activity badge]: https://img.shields.io/github/commit-activity/m/timoth-y/chainmetric-contracts?label=Commit%20activity&color=teal
[blockchain badge]: https://img.shields.io/badge/Blockchain-Hyperledger%20Fabric-informational?style=flat&logo=hyperledger&logoColor=white&labelColor=0A1F1F&color=teal
[kubernetes badge]: https://img.shields.io/badge/Infrastructure-Kubernetes-informational?style=flat&logo=kubernetes&logoColor=white&color=316DE6
[license badge]: https://img.shields.io/badge/License-Apache%202.0-informational?style=flat&color=blue

[this repo]: https://github.com/timoth-y/chainmetric-network
[golang]: https://golang.org
[repo commit activity]: https://github.com/timoth-y/kicksware-api/graphs/commit-activity
[hyperledger fabric url]: https://www.hyperledger.org/use/fabric
[kubernetes url]: https://kubernetes.io
[license url]: https://www.apache.org/licenses/LICENSE-2.0

[fabnctl]: https://github.com/timoth-y/fabnctl

[network deployment section]: https://github.com/timoth-y/chainmetric-network#Deployment
[network.sh script]: https://github.com/timoth-y/chainmetric-network/blob/main/network.sh
[helm]: https://helm.sh

[chaincode]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/glossary.html#smart-contract
[chaincode as external service]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html
[packing chaincode]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html#packaging-chaincode

[chainmetric network repo]: https://github.com/timoth-y/chainmetric-network
[chainmetric sensorsys repo]: https://github.com/timoth-y/chainmetric-sensorsys
[chainmetric app repo]: https://github.com/timoth-y/chainmetric-app

[license file]: https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE
