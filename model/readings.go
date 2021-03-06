package model

import (
	"encoding/json"
	"time"
)

// MetricReadings defines data model for readings from sensors
type MetricReadings struct {
	ID        string                 `json:"id,omitempty"`
	AssetID   string                 `json:"asset_id,omitempty"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Values    map[Metric]interface{} `json:"values,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
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
