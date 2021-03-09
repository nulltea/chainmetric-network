package models

import "encoding/json"

// Device defines device data models
type Device struct {
	ID       string  `json:"id,omitempty"`
	IP       string  `json:"ip,omitempty"`
	MAC      string  `json:"mac,omitempty"`
	Name     string  `json:"name,omitempty"`
	Hostname string  `json:"hostname,omitempty"`
	Profile  string  `json:"profile,omitempty"  metadata:",optional"`
	Supports Metrics `json:"supports,omitempty"`
	Holder   string  `json:"holder,omitempty"`
	State    string  `json:"state,omitempty" metadata:",optional"`
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
