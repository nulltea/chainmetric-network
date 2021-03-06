package model

import (
	"encoding/json"
)

// Requirement defines per metric requirement properties limits
type Requirement struct {
	MinLimit float64 `json:"minThreshold,omitempty"`
	MaxLimit float64 `json:"maxThreshold,omitempty"`
}

// Requirement defines requirements data model
type Requirements struct {
	ID      string                 `json:"id,omitempty"`
	Type    string                 `json:"type,omitempty"`
	AssetID AssetID                `json:"asset_id,omitempty"`
	Metrics map[Metric]Requirement `json:"metrics,omitempty"`
}

func (m Requirements) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m Requirements) Decode(b []byte) (*Requirements, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}
