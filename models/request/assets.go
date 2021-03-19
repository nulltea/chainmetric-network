package request

import (
	"encoding/json"

	"github.com/timoth-y/iot-blockchain-contracts/models"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

type AssetsQuery struct {
	IDs      []string `json:"type,omitempty"`
	Type     *string  `json:"type,omitempty"`
	Holder   *string  `json:"holder,omitempty"`
	State    *string  `json:"state,omitempty"`
	Location *string  `json:"location,omitempty"`
	Tag      *string  `json:"tag,omitempty"`
	Limit    int32    `json:"limit,omitempty"`
	ScrollID string   `json:"scroll_id,omitempty"`
}

func (q *AssetsQuery) Satisfies(asset *models.Asset) bool {
	if len(q.IDs) != 0 && !shared.ContainsString(asset.ID, q.IDs) {
		return false
	}
	if q.Type != nil && asset.Type != *q.Type {
		return false
	}
	if q.Holder != nil && asset.Holder != *q.Holder {
		return false
	}
	if q.Location != nil && asset.Location != *q.Location {
		return false
	}
	if q.State != nil && asset.State != *q.State {
		return false
	}
	if q.Tag != nil && !shared.ContainsString(*q.Tag, asset.Tags) {
		return false
	}

	return true
}

func (m AssetsQuery) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m AssetsQuery) Decode(b []byte) (*AssetsQuery, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}
