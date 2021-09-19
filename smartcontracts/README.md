# Chainmetric: Smart Contracts

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

#### Transactions

| Transaction  | Arguments                    | Response                           | Description                                                      |
| :----------- | :--------------------------- | :--------------------------------- | :--------------------------------------------------------------- |
| **Retrieve** | `id`                         | [`Asset`][asset model]             | Retrieves single record from ledger by a given ID                |
| **All**      | -                            | `[]Asset`                          | Retrieves all records from ledger                                |
| **Query**    | [`AssetsQuery`][asset query] | [`AssetsResponse`][asset response] | Performs rich query against ledger in search of specific records |
| **Upsert**   | `Asset`                      | `id`                               | Inserts new  record into the ledger or updates existing one      |
| **Transfer** | `id`, `holder`               | -                                  | Changes holder of the specific asset                             |
| **Remove**   | `Asset`                      | -                                  | Removes record from the ledger by given ID                       |

#### Example

```go
asset, err := network.GetContract("assets").EvaluateTransaction("Retrieve", id)
```

[asset model]: https://github.com/timoth-y/chainmetric-core/blob/main/models/asset.go#L6
[asset query]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requests/assets.go#L11
[asset response]: https://github.com/timoth-y/chainmetric-network/blob/main/model/response/assets.go#L12

### Devices contract

#### Transactions

| Transaction              | Arguments                                           | Response                                     | Description                                               |
| :----------------------- | :-------------------------------------------------- | :------------------------------------------- | :-------------------------------------------------------- |
| **Retrieve**             | `id`                                                | [`Device`][device model]                     | Retrieves single record from ledger by a given ID         |
| **All**                  | -                                                   | `[]Device`                                   | Retrieves all records from ledger                         |
| **Register**             | `[]Device`                                          | `id`                                         | Creates and registers new device in the blockchain ledger |
| **Update**               | `id`, `Device`                                      | `Device`                                     | Updates device state in ledger with requested properties  |
| **Unbind**               | `id`                                                | -                                            | Removes record from the ledger by given ID                |
| **Command**              | [`DeviceCommandRequest`][command request]           | -                                            | Handles execution requests for devices                    |
| **SubmitCommandResults** | `entryID`, [`CommandResultsSubmit`][command result] | -                                            | Updates command log entry in the ledger                   |
| **CommandsLog**          | `deviceID`                                          | [`DeviceCommandLogEntry`][command log entry] | Retrieves entire commands log from the blockchain ledger  | 

[device model]: https://github.com/timoth-y/chainmetric-core/blob/main/models/device.go
[command request]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requests/commands.go#L13
[command result]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requests/commands.go#L74
[command log entry]: https://github.com/timoth-y/chainmetric-core/blob/main/models/command.go#L21

#### Example

```go
id, err := network.GetContract("devices").EvaluateTransaction("Register", device)
```

### Requirements contact

#### Transactions

| Transaction   | Arguments     | Response                             | Description                                                         |
| :------------ | :------------ | :----------------------------------- | :------------------------------------------------------------------ |
| **Retrieve**  | `id`          | [`Requirements`][requirements model] | Retrieves single record from ledger by a given ID                   |
| **ForAsset**  | `assetID`     | `[]Requirements`                     | Retrieves all requirements from ledger for specific asset           |
| **ForAssets** | `assetIDs`    | `[]Requirements`                     | Retrieves all requirements from ledger for specific multiply assets |
| **Assign**    | `Requiremnts` | `id`                                 | Assigns requirements to an asset and stores it on the ledger        |
| **Revoke**    | `id`          | -                                    | Revokes requirements from an asset and removes it from the ledger   |

[requirements model]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requirements.go#L18

#### Example

```go
requirements, err := network.GetContract("requirements").EvaluateTransaction("ForAssets", assetsIDs)
```

### Readings contact

#### Transactions

| Transaction           | Arguments                          | Response                                      | Description                                                             |
| :-------------------- | :--------------------------------- | :-------------------------------------------- | :---------------------------------------------------------------------- |
| **ForAsset**          | `assetID`                          | [`MetricReadingsResponse`][readings response] | Retrieves all records records from ledger for specific asset            |
| **ForMetric**         | `assetID`, `metricID`              | `MetricReadingsResponse`                      | Retrieves all records records from ledger for specific asset and metric |
| **Post**              | [`MetricReadings`][readings model] | `id`                                          | Inserts new metric readings record into the ledger                      |
| **BindToEventSocket** | `assetID`, `metricID`              | `eventToken`                                  | Creates event socket subscription ticket for connected party, so that each posted readings record, which satisfies client request for given asset ID and metric, would be send directly to it via event streaming |
| **CloseEventSocket**  | `eventToken`                       | -                                             | Revokes event socket subscription ticket for connected party            |

#### Example

```go
readings, err := network.GetContract("readings").EvaluateTransaction("ForMetric", assetsID, metricID)
```

[readings model]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requirements.go#L18
[readings response]: https://github.com/timoth-y/chainmetric-core/blob/main/models/readings.go#L9


## Deployment

[Chaincodes][chaincode] (alternative to Smart Contracts) in Hyperledger Fabric infrastructure can be deployed both by embedding their source code into the blockchain peers and by deploying them [as external services][chaincode as external service], which is a way more versatile option especially for Kubernetes cluster environment where such Chaincodes can be deployed as a pods.

For chaincodes initial deployment or further updates use `deploy cc` command as following: 

```shell
fabnctl install cc --arch=arm64 --domain=chainmetric.network --chaincode=assets -C=supply-channel \
   -o=org1 -p=peer0 \
   -o=org2 -p=peer0 \
   -o=org3 -p=peer0 \
   --registry=dockerhubuser ./chaincodes
```

For more detailed instructions please refer to `fabnctl` [documentation](https://github.com/timoth-y/fabnctl#deploy-chaincodes).

## Roadmap

- [x] Device remote commands over blockchain events [(#1)](https://github.com/timoth-y/chainmetric-network/pull/1)
- [x] Cache layer for storing contracts operational data e.g. `EventSocketSubscriptionTicket` [(#4)](https://github.com/timoth-y/chainmetric-network/pull/4)
- [x] Devices location management business logic [(#2)](https://github.com/timoth-y/chainmetric-network/pull/2)
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
