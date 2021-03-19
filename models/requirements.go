package models

import (
	"encoding/json"
	"errors"
)

// Requirement defines per metric requirement properties limits
type Requirement struct {
	MinLimit float64 `json:"min_limit"`
	MaxLimit float64 `json:"max_limit"`
}

// RequirementsMap represents map with Metric key and corresponding Requirement value
type RequirementsMap map[Metric]Requirement

// Requirement defines requirements data models
type Requirements struct {
	ID      string          `json:"id,omitempty"`
	AssetID string          `json:"asset_id,omitempty"`
	Period  int             `json:"period,omitempty" metadata:",optional"`
	Metrics RequirementsMap `json:"metrics,omitempty"`
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

func (m *Requirements) Validate() error {
	if len(m.ID) == 0 {
		return errors.New("id must be assigned to requirements")
	}

	if len(m.AssetID) == 0 {
		return errors.New("requirements must be assigned to asset (assetID is required)")
	}

	if len(m.Metrics) == 0 {
		return errors.New("requirements must to contain at least one metric")
	}

	return nil
}

func (rm RequirementsMap) Metrics() Metrics {
	var (
		metrics = make(Metrics, len(rm))
		i = 0
	)

	for m := range rm {
		metrics[i] = m
		i++
	}

	return metrics
}
