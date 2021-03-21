package models

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Device defines device data models
type Device struct {
	ID       string      `json:"id"`
	IP       string      `json:"ip"`
	MAC      string      `json:"mac,omitempty" metadata:",optional"`
	Name     string      `json:"name,omitempty" metadata:",optional"`
	Hostname string      `json:"hostname"`
	Profile  string      `json:"profile,omitempty" metadata:",optional"`
	Supports Metrics     `json:"supports"`
	Holder   string      `json:"holder"`
	State    DeviceState `json:"state,omitempty" metadata:",optional"`
	Location string      `json:"location"`
}

type DeviceState string

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

func (m *Device) Validate() error {
	if len(m.ID) == 0 {
		return errors.New("id must be assigned to device")
	}

	if len(m.Hostname) == 0 {
		return errors.New("hostname must be assigned to device")
	}

	if len(m.Supports) == 0 {
		return errors.New("device must to support at least one metric")
	}

	return nil
}
