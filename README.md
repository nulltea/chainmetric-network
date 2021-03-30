# ChainMetric: Smart Contracts

[![golang badge]][golang]&nbsp;
[![commit activity badge]][repo commit activity]&nbsp;
[![blockchain badge]][hyperledger fabric url]&nbsp;
[![kubernetes badge]][kubernetes url]&nbsp;
[![license badge]][license url]

## Overview

_**Chainmetric Smart Contracts**_ are designed to grand access to blockchain-stored data while implementing such data validation, aggregation, and management functionality.

Being a part of a permissioned blockchain network based on Hyperledger Fabric stack, such Contracts are written in Go and configured to be deployed as Kubernetes services communicating with each other, blockchain peers, and external applications via gRPC protocol and event streaming.

## Requirements

- Kubernetes cluster with previously deployed Chainmetric network ([see network deployment procedure][network deployment section])
- [Helm][helm] binaries must be presented on a local machine from which deployment script will be used.

## Deployment

Bash-script from [network repository][chainmetric network repo] provides a straightforward way of deployment and further upgrading of Smart Contracts, which are essential for a current blockchain solution.

[Chaincodes][chaincode] (alternative to Smart Contracts) in Hyperledger Fabric infrastructure can be deployed both by embedding their source code into the blockchain peers and by deploying them [as external services][chaincode as external service], which is a way more versatile option especially for Kubernetes cluster environment where such Chaincodes can be deployed as a pods.

To deploy Smart Contract for the first time it is required to [pack its configuration][packing chaincode] in an archive and then use it with `peer lifecycle commands`. All required steps are conveniently aggregated in a single command of `network.sh` [script][network.sh script] from the [network repository][chainmetric network repo]. To perform Chaincode's initial deployment use `deploy` command with `cc` action as following:

```
$ ./network.sh deploy cc --cc_name=`chaincode name` --channel='channel name' --peer='peer subdomain name' --org='organization name'
```

After this whenever an upgrade must be performed simply add `--upgrade` flag to the previously executed command. That will rebuild the docker image, send it to the dedicated registry and redeploy the Helm chart.

## Roadmap

- Requirements violation event streaming
- Violations rule engine
- Transaction for assets changes history retrieving
- Devices location management business logic
- Users contract

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
