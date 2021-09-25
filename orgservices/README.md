# Chainmetric: Off-Chain Services

In order to support and extend both functionality and availability of Chainmetric's blockchain network, additional per-organization off-chain API services are introduced.

Such services outsource execution of operations that require integration with various infrastructure components (HashiCorp Vault, Mongo DB, Firebase) and cannot be performed in a deterministic way as Fabric's on-chain Smart Contracts require.

## Requirements
- Kubernetes environment with Chainmetric network ([see network deployment procedure][network deployment]) and [HashiCorp Vault][vault] installed
- [Basel build][bazel] command line utility installed on local machine
- Some services require [Firebase][firebase] account for closer integration with user's application.

[network deployment]: https://github.com/timoth-y/chainmetric-network#Deployment
[vault]: https://www.hashicorp.com/products/vault
[firebase]: https://firebase.google.com
[bazel]: https://bazel.build

## Services reference
As mentioned services are deployed per-organization and are therefore accessed only by users of each concrete organization.

Services export their functionality via gRPC API, which have a positive impact both of communication speed and thereafter user experience, as well as integration ease and developer experience.

### Identity service
Manages user's identities: creates new ones on sign-up, authenticates existing users, authorities based on roles and privileges.

See details at: [**orgservices/identity**](https://github.com/timoth-y/chainmetric-network/tree/main/orgservices/identity).

### Notifications service
Allows users to subscribe to changes and events on Blockchain ledger with intend to receive convenient push notifications even when application is closed.

See details at: [**orgservices/notifications**](https://github.com/timoth-y/chainmetric-network/main/orgservices/notifications).

## Deployment
For building image out of source code use `build-orgservice` rule specifying name of the service as following:
```bash
$ make service=identity build-orgservice
```

For installing service with pre-build image use `install-orgservice` rule specifying both name of the service and target organization as following:
```bash
$ make service=identity org=org1 install-orgservice
```

This will create all required secrets in current Kube context and `upgrade --install` Helm chart.

Rule `deploy-orgservice` is a combination of both previous ones:
```bash
$ make service=identity org=org1 deploy-orgservice
```

## Development
For initializing local development environment use `bazel run` command specified gazelle plugin target.

```bash
$ bazel run //:gazelle
```

To paste Protobuf generated files in directories where proto files are defined use following command:

```bash
$ bazel query 'kind("proto_link", //...)'  | xargs -L 1 bazel run
```

## Wrap up

Chainmetric org-services are meant to hide intricacies of deterministic Fabric Blockchain infrastructure and provide batter user and developer experience.

## License
Licensed under the [Apache 2.0][license file].

[license file]: https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE
