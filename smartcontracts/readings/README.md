# Chainmetric: Readings Contract

**Readings Smart Contract** exposes transactions for posting metric readings sourced by sensors of IoT devices.

Allows creating socket connection for receiving updates in real time.

Validates received readings against requirements and emits events on ones violations.

## Flowchart

![flowchart]

[flowchart]: https://github.com/timoth-y/chainmetric-network/blob/main/docs/readings-contract-flowchart.png?raw=true
*Event socket flowchart*

Also see: [Application environment monitoring](https://github.com/timoth-y/chainmetric-app#environment-monitoring).

## Transactions reference:

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

## License

Licensed under the [Apache 2.0](https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE).
