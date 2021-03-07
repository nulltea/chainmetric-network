package models

import "encoding/json"

// Device defines device data models
type Device struct {
	ID       string  `json:"id,omitempty"`
	URL      string  `json:"url,omitempty"`
	Name     string  `json:"name,omitempty"`
	Profile  string  `json:"profile,omitempty"`
	Supports Metrics `json:"supports,omitempty"`
	Holder   string  `json:"holder,omitempty"`
	State    string  `json:"state,omitempty"`
	Location string  `json:"location,omitempty"`
}

func (m Device) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m Device) Decode(b []byte) (*Device, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}
