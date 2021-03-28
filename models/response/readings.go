package response

import (
	"time"

	"github.com/timoth-y/iot-blockchain-contracts/models"
)

// MetricReadingsResponse aggregates readings for asset into per metric reading streams
type MetricReadingsResponse struct {
	AssetID string                                 `json:"asset_id"`
	Streams map[models.Metric]MetricReadingsStream `json:"streams"`
}

// MetricReadingsStream defines a single metric readings stream
type MetricReadingsStream []MetricReadingsPoint

// MetricReadingsPoint defines a single point in time of a metric reading stream
type MetricReadingsPoint struct {
	DeviceID  string      `json:"device_id"`
	Location  string      `json:"location"`
	Timestamp time.Time   `json:"timestamp"`
	Value     interface{} `json:"value"`
}
