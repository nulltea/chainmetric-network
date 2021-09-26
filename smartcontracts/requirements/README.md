# Chainmetric: Requirements Contract

**Requirements Smart Contract** exposes transactions for assigning, changing and revoking environmental requirements for storing and transferring assets.

## Transactions reference:

| Transaction   | Arguments     | Response                             | Description                                                         |
| :------------ | :------------ | :----------------------------------- | :------------------------------------------------------------------ |
| **Retrieve**  | `id`          | [`Requirements`][requirements model] | Retrieves single record from ledger by a given ID                   |
| **ForAsset**  | `assetID`     | `[]Requirements`                     | Retrieves all requirements from ledger for specific asset           |
| **ForAssets** | `assetIDs`    | `[]Requirements`                     | Retrieves all requirements from ledger for specific multiply assets |
| **Assign**    | `Requiremnts` | `id`                                 | Assigns requirements to an asset and stores it on the ledger        |
| **Revoke**    | `id`          | -                                    | Revokes requirements from an asset and removes it from the ledger   |

[requirements model]: https://github.com/timoth-y/chainmetric-core/blob/main/models/requirements.go#L18

### Example

```go
requirements, err := network.GetContract("requirements").EvaluateTransaction("ForAssets", assetsIDs)
```

## License

Licensed under the [Apache 2.0](https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE).
