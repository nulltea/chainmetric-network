# ChainMetric: Smart Contracts

[![golang badge]][golang]&nbsp;
[![commit activity badge]][repo commit activity]&nbsp;
[![blockchain badge]][hyperledger fabric url]&nbsp;
[![kubernetes badge]][kubernetes url]&nbsp;
[![license badge]][license url]

## Overview

_**Chainmetric Smart Contracts**_ are designed to grand access to blockchain-stored data while implementing such data validation, aggregation, and management functionality.

Being a part of a permissioned blockchain network based on Hyperledger Fabric stack,
such Contracts are written in Go and configured to be deployed as Kubernetes services communicating with each other, blockchain peers, and external applications via gRPC protocol and event streaming.

## Requirements

- Deployed Chainmetric network ([see network deployment procedure][network deployment section]) preferably in Kubernetes environment
- Image registry where chaincode docker images will be hosted (Docker Hub account will do)
- [`fabnctl`][fabnctl] command line utility installed on local machine

## Contracts reference

### Assets contract

#### Retrieve single asset by ID

**`Retrieve`** retrieves single `models.Asset` record from blockchain ledger by a given `id`.

```go
asset, err := network.GetContract("assets").EvaluateTransaction("Retrieve", id)
```
Response model:

```json
{
  "id": "string",
  "sku": "string",
  "name": "string",
  "type": "string",
  "info": "string",
  "holder": "string",
  "state": "string",
  "location": {
    "name": "string",
    "latitude": 0,
    "longitude": 0
  },
  "tags": [
    "string"
  ]
}
```

#### Get all assets

**`All`** retrieves all `models.Asset` records from blockchain ledger.

```go
assets, err := network.GetContract("assets").EvaluateTransaction("All")
```

#### Query assets

**`Query`** performs rich query against blockchain ledger in search of specific `models.Asset` records.

To support pagination it returns results wrapped in `response.AssetsResponse`,
where 'scroll_id' will contain special key for continuing from where the previous request ended.

Request model:

```json
{
  "holder": "string",
  "state": "string",
  "location": {
    "location": {
      "name": "",
      "latitude": 0,
      "longitude": 0
    },
    "distance": 0
  },
  "tag": [
    "string"
  ],
  "scroll_id": "string"
}
```

```go
assets, err := network.GetContract("assets").EvaluateTransaction("Query", query)
```

Response model is same as `[]models.Asset` but with addition of `scroll_id` string value.

#### Insert or update asset

**`Upsert`** inserts new `models.Asset` record into the blockchain ledger or updates existing one.

```go
id, err := network.GetContract("assets").EvaluateTransaction("Update", asset)
```

#### Transfer asset ownership

**`Transfer`** changes holder of the specific `models.Asset`.

```go
err := network.GetContract("assets").EvaluateTransaction("Transfer", id, newHolder)
```

#### Delete assets from ledger

**`Remove`** removes `models.Asset` from the blockchain ledger.


```go
err := network.GetContract("assets").EvaluateTransaction("Remove", id)
```

### Devices contract

#### Retrieve single device by ID

**`Retrieve`** retrieves single `models.Device` record from blockchain ledger  by a given `id`.

```go
device, err := network.GetContract("devices").EvaluateTransaction("Retrieve", id)
```
Response model:

```json
{
  "id": "string",
  "ip": "string",
  "mac": "string",
  "name": "string",
  "hostname": "string",
  "profile": "string",
  "supports": [
    "string"
  ],
  "holder": "string",
  "state": "string",
  "battery": {
    "level": 0,
    "plugged": true
  },
  "location": {
    "name": "string",
    "latitude": 0,
    "longitude": 0
  }
}
```

#### Retrieve all devices

**`All`** retrieves all `models.Device` records from blockchain ledger.

```go
devices, err := network.GetContract("devices").EvaluateTransaction("All")
```

#### Register device on network

**`Register`** creates and registers new device in the blockchain ledger.

```go
id, err := network.GetContract("devices").EvaluateTransaction("Register", device)
```
#### Update device

**`Update`** updates `models.Device` state in blockchain ledger with requested properties.

```go
id, err := network.GetContract("devices").EvaluateTransaction("Update", id, device)
```

#### Unbind device

**`Unbind`** removes `models.Device` from the blockchain ledger.

```go
id, err := network.GetContract("devices").EvaluateTransaction("Unbind", id)
```

## Deployment

[Chaincodes][chaincode] (alternative to Smart Contracts) in Hyperledger Fabric infrastructure can be deployed both by embedding their source code into the blockchain peers and by deploying them [as external services][chaincode as external service], which is a way more versatile option especially for Kubernetes cluster environment where such Chaincodes can be deployed as a pods.

To deploy Smart Contract for the first time it is required to [pack its configuration][packing chaincode] in an archive and then use it with `peer lifecycle commands`. All required steps are conveniently aggregated in a single command of `network.sh` [script][network.sh script] from the [network repository][chainmetric network repo]. To perform Chaincode's initial deployment use `deploy` command with `cc` action as following:

```shell
fabnctl deploy cc --arch=arm64 --domain=example.network --chaincode=example -C=example-channel \
   -o=org1 -p=peer0 \
   -o=org2 -p=peer0 \
   -o=org3 -p=peer0 \
   --registry=dockerhubuser ./chaincodes/example
```

After this whenever an upgrade must be performed simply add `--upgrade` flag to the previously executed command. That will rebuild the docker image, send it to the dedicated registry and redeploy the Helm chart.

## Roadmap

- [x] Device remote commands over blockchain events [(#1)](https://github.com/timoth-y/chainmetric-contracts/pull/1)
- [x] Cache layer for storing contracts operational data e.g. `EventSocketSubscriptionTicket` [(#4)](https://github.com/timoth-y/chainmetric-contracts/pull/4)
- [x] Devices location management business logic [(#2)](https://github.com/timoth-y/chainmetric-contracts/pull/2)
- [ ] Requirements violation notification
- [ ] Requirements violations rule engine
- [ ] Transaction for assets changes history retrieving
- [ ] Users contract

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

[this repo]: https://github.com/timoth-y/chainmetric-contracts
[golang]: https://golang.org
[repo commit activity]: https://github.com/timoth-y/kicksware-api/graphs/commit-activity
[hyperledger fabric url]: https://www.hyperledger.org/use/fabric
[kubernetes url]: https://kubernetes.io
[license url]: https://www.apache.org/licenses/LICENSE-2.0

[network deployment section]: https://github.com/timoth-y/chainmetric-network#Deployment
[network.sh script]: https://github.com/timoth-y/chainmetric-network/blob/main/network.sh
[helm]: https://helm.sh

[chaincode]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/glossary.html#smart-contract
[chaincode as external service]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html
[packing chaincode]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html#packaging-chaincode

[chainmetric network repo]: https://github.com/timoth-y/chainmetric-network
[chainmetric sensorsys repo]: https://github.com/timoth-y/chainmetric-sensorsys
[chainmetric app repo]: https://github.com/timoth-y/chainmetric-app



[license file]: https://github.com/timoth-y/chainmetric-contracts/blob/main/LICENSE
