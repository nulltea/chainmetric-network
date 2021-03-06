package model

import "encoding/json"

// Asset defines asset data model
type Asset struct {
	ID       AssetID  `json:"id,omitempty"`
	SKU      string   `json:"sku,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Info     string   `json:"info,omitempty"`
	Cost     float64  `json:"cost,omitempty"`
	Amount   int      `json:"amount,omitempty"`
	Holder   string   `json:"holder,omitempty"`
	State    string   `json:"state,omitempty"`
	Location string   `json:"location,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// AssetID defines composite identifier of Asset
type AssetID string

func (m Asset) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m Asset) Decode(b []byte) (*Asset, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}

func (id AssetID) String() string {
	return string(id)
}
