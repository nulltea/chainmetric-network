package couchdb

import (
	"encoding/json"

	"github.com/timoth-y/chainmetric-core/models"
)

const (
	AssetRecordType           = "asset"
	DeviceRecordType          = "device"
	RequirementsRecordType    = "requirements"
	ReadingsRecordType        = "readings"
	CommandLogEntryRecordType = "device_command"
)

// CouchDB specific records model structure types for each entity: Asset, Device, Requirements, MetricReadings.
type (
	// Asset wraps models.Asset with additional database specific fields.
	Asset struct {
		RecordType string `json:"record_type"` // Constant value of 'record_type': AssetRecordType
		*models.Asset
	}

	// Device wraps models.Device with additional database specific fields.
	Device struct {
		RecordType string `json:"record_type"` // Constant value of 'record_type': DeviceRecordType
		*models.Device
	}

	// Requirements wraps models.Requirements with additional database specific fields.
	Requirements struct {
		RecordType string `json:"record_type"` // Constant value of 'record_type': RequirementsRecordType
		*models.Requirements
	}

	// MetricReadings wraps models.MetricReadings with additional database specific fields.
	MetricReadings struct {
		RecordType string `json:"record_type"` // Constant value of 'record_type': ReadingsRecordType
		*models.MetricReadings
	}

	// DeviceCommandLogEntry wraps models.DeviceCommandLogEntry with additional database specific fields.
	DeviceCommandLogEntry struct {
		RecordType string `json:"record_type"` // Constant value of 'record_type': CommandLogEntryRecordType
		*models.DeviceCommandLogEntry
	}
)

// NewAssetRecord constructs new Asset record based on models.Asset with predefining database specific fields.
func NewAssetRecord(base *models.Asset) *Asset {
	return &Asset{
		RecordType: AssetRecordType,
		Asset:      base,
	}
}

// NewDeviceRecord constructs new Device record based on models.Device with predefining database specific fields.
func NewDeviceRecord(base *models.Device) *Device {
	return &Device{
		RecordType: DeviceRecordType,
		Device:     base,
	}
}

// NewRequirementsRecord constructs new Requirements record based on models.Requirements
// with predefining database specific fields.
func NewRequirementsRecord(base *models.Requirements) *Requirements {
	return &Requirements{
		RecordType:   RequirementsRecordType,
		Requirements: base,
	}
}

// NewMetricReadingsRecord constructs new MetricReadings record based on models.MetricReadings
// with predefining database specific fields.
func NewMetricReadingsRecord(base *models.MetricReadings) *MetricReadings {
	return &MetricReadings{
		RecordType:     ReadingsRecordType,
		MetricReadings: base,
	}
}

// NewDeviceCommandLogEntry constructs new DeviceCommand record based on models.DeviceCommandLogEntry
// with predefining database specific fields.
func NewDeviceCommandLogEntry(base *models.DeviceCommandLogEntry) *DeviceCommandLogEntry {
	return &DeviceCommandLogEntry{
		RecordType:            CommandLogEntryRecordType,
		DeviceCommandLogEntry: base,
	}
}

// Encode serializes the Asset record.
func (m *Asset) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}

	return data
}

// Encode serializes the Device record.
func (m *Device) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}

	return data
}

// Encode serializes the Requirements record.
func (m *Requirements) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}

	return data
}

// Encode serializes the MetricReadings record.
func (m *MetricReadings) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}

	return data
}

// Encode serializes the DeviceCommandLogEntry record.
func (m *DeviceCommandLogEntry) Encode() []byte {
	data, err := json.Marshal(m); if err != nil {
		return nil
	}

	return data
}
