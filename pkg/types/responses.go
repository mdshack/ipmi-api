package types

import (
	"fmt"
	"github.com/bougou/go-ipmi"
)

// IPMIAuth holds the authentication information for IPMI connections
type IPMIAuth struct {
	Host     string
	Port     int
	Username string
	Password string
}

// APIResponse is the standard response format for all API endpoints
type APIResponse struct {
	Success        bool        `json:"success"`
	Data           interface{} `json:"data,omitempty"`
	Error          string      `json:"error,omitempty"`
	Timestamp      string      `json:"timestamp"`
	ExecutionTime  int64       `json:"execution_time_ms,omitempty"`
}

// SystemInfo represents system information from IPMI
type SystemInfo struct {
	DeviceID          uint8  `json:"device_id"`
	DeviceRevision    uint8  `json:"device_revision"`
	FirmwareRev       string `json:"firmware_revision"`
	IPMIVersion       string `json:"ipmi_version"`
	ManufacturerID    uint32 `json:"manufacturer_id"`
	ProductID         uint16 `json:"product_id"`
	ProvideDeviceSDRs bool   `json:"provide_device_sdrs"`
	DeviceAvailable   bool   `json:"device_available"`
}

// PowerStatus represents the power status of a system
type PowerStatus struct {
	Status string `json:"status"`
}

// SensorData represents sensor information
// We'll use the library's Sensor type directly in responses where possible
type SensorData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
	State string `json:"state"`
}

// ChassisStatus represents chassis status information
type ChassisStatus struct {
	PowerState     string `json:"power_state"`
	LastPowerEvent string `json:"last_power_event"`
}

// BootDevice represents boot device configuration
type BootDevice struct {
	Device string `json:"device"`
}

// User represents IPMI user information
// We'll use the library's User type directly in responses where possible
type User struct {
	ID                   uint8  `json:"id"`
	Name                 string `json:"name"`
	Callin               bool   `json:"callin"`
	LinkAuthEnabled      bool   `json:"link_auth_enabled"`
	IPMIMessagingEnabled bool   `json:"ipmi_messaging_enabled"`
	MaxPrivLevel         string `json:"max_priv_level"`
}

// FRUData represents FRU (Field Replaceable Unit) information
// We'll use the library's FRU type directly in responses where possible
type FRUData struct {
	DeviceID   uint8  `json:"device_id"`
	DeviceName string `json:"device_name"`
	Present    bool   `json:"present"`
}

// SELEntry represents a System Event Log entry
type SELEntry struct {
	RecordID  uint16 `json:"record_id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

// LANConfig represents LAN configuration
type LANConfig struct {
	IPAddress      string `json:"ip_address"`
	MACAddress     string `json:"mac_address"`
	SubnetMask     string `json:"subnet_mask"`
	GatewayIP      string `json:"gateway_ip"`
	GatewayMAC     string `json:"gateway_mac"`
	PrimaryDNS     string `json:"primary_dns"`
	SecondaryDNS   string `json:"secondary_dns"`
	VLANID         uint16 `json:"vlan_id"`
	VLANPriority   uint8  `json:"vlan_priority"`
}

// ConvertDeviceIDResponse converts a go-ipmi GetDeviceIDResponse to our SystemInfo type
func ConvertDeviceIDResponse(resp *ipmi.GetDeviceIDResponse) *SystemInfo {
	return &SystemInfo{
		DeviceID:          resp.DeviceID,
		DeviceRevision:    resp.DeviceRevision,
		FirmwareRev:       resp.FirmwareVersionStr(),
		IPMIVersion:       fmt.Sprintf("%d.%d", resp.MajorIPMIVersion, resp.MinorIPMIVersion),
		ManufacturerID:    resp.ManufacturerID,
		ProductID:         resp.ProductID,
		ProvideDeviceSDRs: resp.ProvideDeviceSDRs,
		DeviceAvailable:   resp.DeviceAvailable,
	}
}

// ConvertUser converts a go-ipmi User to our User type
func ConvertUser(user *ipmi.User) *User {
	return &User{
		ID:                   user.ID,
		Name:                 user.Name,
		Callin:               user.Callin,
		LinkAuthEnabled:      user.LinkAuthEnabled,
		IPMIMessagingEnabled: user.IPMIMessagingEnabled,
		MaxPrivLevel:         user.MaxPrivLevel.String(),
	}
}

// ConvertFRU converts a go-ipmi FRU to our FRUData type
func ConvertFRU(fru *ipmi.FRU) *FRUData {
	return &FRUData{
		DeviceID:   fru.DeviceID(),
		DeviceName: fru.DeviceName(),
		Present:    fru.Present(),
	}
}