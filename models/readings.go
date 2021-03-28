package models

import (
	"encoding/json"
	"time"
)

// MetricReadings defines data models for readings from sensors
type MetricReadings struct {
	ID        string                 `json:"id"`
	AssetID   string                 `json:"asset_id"`
	DeviceID  string                 `json:"device_id"`
	Location  string                 `json:"location"`
	Timestamp time.Time              `json:"timestamp"`
	Values    map[Metric]interface{} `json:"values"`
}

func (m MetricReadings) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m MetricReadings) Decode(b []byte) (*MetricReadings, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}
