# Chainmetric: Assets Contract

**Assets Smart Contract** exposes transactions for registering assets on Blockchain ledger, log changes of their state, transfer ownership, etc.

## Transactions reference:

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

## License

Licensed under the [Apache 2.0](https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE).
