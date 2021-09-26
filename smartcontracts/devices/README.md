# Chainmetric: Devices Contract

**Devices Smart Contract** Exposes transactions for registering IoT devices on chain, issue remote commands, transfer ownership, etc.

## Transactions reference:

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

## License

Licensed under the [Apache 2.0](https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE).
