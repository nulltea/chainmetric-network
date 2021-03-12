package request

import (
	"encoding/json"

	"github.com/timoth-y/iot-blockchain-contracts/models"
)

// DeviceUpdateRequest defines update request for models.Device
type DeviceUpdateRequest struct {
	Name     *string             `json:"name,omitempty"`
	Profile  *string             `json:"profile,omitempty"`
	Supports models.Metrics     `json:"supports,omitempty"`
	Holder   *string             `json:"holder,omitempty"`
	State    *models.DeviceState `json:"state,omitempty"`
	Location *string             `json:"location,omitempty"`
}

// Update updates models.Device
func (u *DeviceUpdateRequest) Update(device *models.Device) {
	if u.Name != nil {
		device.Name = *u.Name
	}

	if u.Profile != nil {
		device.Profile = *u.Profile
	}

	if u.Supports != nil {
		device.Supports = u.Supports
	}

	if u.Holder != nil {
		device.Holder = *u.Holder
	}

	if u.State != nil {
		device.State = *u.State
	}

	if u.Location != nil {
		device.Location = *u.Location
	}
}

func (m DeviceUpdateRequest) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}
	return data
}

func (m DeviceUpdateRequest) Decode(b []byte) (*DeviceUpdateRequest, error) {
	err := json.Unmarshal(b, &m)
	return &m, err
}
