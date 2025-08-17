package ipmi

import (
	"context"
	"fmt"
	"time"

	"github.com/bougou/go-ipmi"
	"github.com/mdshack/ipmi-api/pkg/types"
)

// Client wraps the go-ipmi client with convenience methods
type Client struct {
	client  *ipmi.Client
	timeout time.Duration
}

// NewClient creates a new IPMI client wrapper
func NewClient(client *ipmi.Client) *Client {
	return &Client{
		client:  client,
		timeout: 30 * time.Second, // Default timeout
	}
}

// GetSystemInfo retrieves system information
func (c *Client) GetSystemInfo() (*types.SystemInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	resp, err := c.client.GetDeviceID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetDeviceID: %w", err)
	}
	
	return types.ConvertDeviceIDResponse(resp), nil
}

// GetSystemGUID retrieves the system GUID
func (c *Client) GetSystemGUID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	resp, err := c.client.GetSystemGUID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute GetSystemGUID: %w", err)
	}
	
	// Convert [16]byte to string representation
	guid := fmt.Sprintf("%x-%x-%x-%x-%x",
		resp.GUID[0:4], resp.GUID[4:6], resp.GUID[6:8], resp.GUID[8:10], resp.GUID[10:16])
	
	return guid, nil
}

// GetPowerStatus retrieves the current power status
func (c *Client) GetPowerStatus() (*types.PowerStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	resp, err := c.client.GetChassisStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetChassisStatus: %w", err)
	}
	
	status := "unknown"
	if resp.PowerIsOn {
		status = "on"
	} else {
		status = "off"
	}
	
	return &types.PowerStatus{
		Status: status,
	}, nil
}

// PowerOn sends a power on command
func (c *Client) PowerOn() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	_, err := c.client.ChassisControl(ctx, ipmi.ChassisControlPowerUp)
	if err != nil {
		return fmt.Errorf("failed to execute ChassisControl (PowerUp): %w", err)
	}
	
	return nil
}

// PowerOff sends a power off command
func (c *Client) PowerOff() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	_, err := c.client.ChassisControl(ctx, ipmi.ChassisControlPowerDown)
	if err != nil {
		return fmt.Errorf("failed to execute ChassisControl (PowerDown): %w", err)
	}
	
	return nil
}

// PowerCycle sends a power cycle command
func (c *Client) PowerCycle() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	_, err := c.client.ChassisControl(ctx, ipmi.ChassisControlPowerCycle)
	if err != nil {
		return fmt.Errorf("failed to execute ChassisControl (PowerCycle): %w", err)
	}
	
	return nil
}

// PowerReset sends a power reset command
func (c *Client) PowerReset() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	_, err := c.client.ChassisControl(ctx, ipmi.ChassisControlHardReset)
	if err != nil {
		return fmt.Errorf("failed to execute ChassisControl (HardReset): %w", err)
	}
	
	return nil
}

// GetSensors retrieves sensor data
func (c *Client) GetSensors() ([]*ipmi.Sensor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	sensors, err := c.client.GetSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensors: %w", err)
	}
	
	return sensors, nil
}

// GetChassisStatus retrieves chassis status
func (c *Client) GetChassisStatus() (*types.ChassisStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	resp, err := c.client.GetChassisStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetChassisStatus: %w", err)
	}
	
	powerState := "unknown"
	if resp.PowerIsOn {
		powerState = "on"
	} else {
		powerState = "off"
	}
	
	lastPowerEvent := "unknown"
	if resp.LastPowerOnByCommand {
		lastPowerEvent = "power up"
	} else if resp.LastPowerDownByPowerFault {
		lastPowerEvent = "power down"
	} else if resp.LastPowerDownByPowerInterlockActivated {
		lastPowerEvent = "power cycle"
	} else {
		lastPowerEvent = "other"
	}
	
	return &types.ChassisStatus{
		PowerState:     powerState,
		LastPowerEvent: lastPowerEvent,
	}, nil
}

// GetBootDevice retrieves the current boot device
func (c *Client) GetBootDevice() (*types.BootDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	params, err := c.client.GetSystemBootOptionsParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get boot options: %w", err)
	}
	
	bootFlags := params.BootFlags
	if bootFlags == nil {
		return nil, fmt.Errorf("boot flags not available")
	}
	
	device := "unknown"
	// Since we don't have access to the specific constants, we'll just use the raw value
	device = fmt.Sprintf("selector %d", bootFlags.BootDeviceSelector)
	
	return &types.BootDevice{
		Device: device,
	}, nil
}

// SetBootDevice sets the boot device
func (c *Client) SetBootDevice(device string) error {
	// Since we don't have access to the specific constants, we'll just return an error
	return fmt.Errorf("setting boot device is not fully implemented due to missing constants")
}

// GetUsers retrieves all IPMI users
func (c *Client) GetUsers() ([]*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Get users from channel 1 (typically the primary channel)
	users, err := c.client.GetUsers(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	
	var result []*types.User
	for _, user := range users {
		result = append(result, types.ConvertUser(user))
	}
	
	return result, nil
}

// GetUser retrieves a specific IPMI user
func (c *Client) GetUser(userID uint8) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	resp, err := c.client.GetUsername(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %w", err)
	}
	
	accessResp, err := c.client.GetUserAccess(ctx, 1, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user access: %w", err)
	}
	
	// Create a user object with the available information
	user := &ipmi.User{
		ID:                   userID,
		Name:                 resp.Username,
		LinkAuthEnabled:      accessResp.LinkAuthEnabled,
		IPMIMessagingEnabled: accessResp.IPMIMessagingEnabled,
		MaxPrivLevel:         ipmi.PrivilegeLevel(accessResp.MaxPrivLevel),
	}
	
	return types.ConvertUser(user), nil
}

// CreateUser creates a new IPMI user
func (c *Client) CreateUser(userID uint8, username, password string, privLevel string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Set username
	_, err := c.client.SetUsername(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("failed to set username: %w", err)
	}
	
	// Set password
	_, err = c.client.SetUserPassword(ctx, userID, password, false)
	if err != nil {
		return fmt.Errorf("failed to set user password: %w", err)
	}
	
	// Enable user
	err = c.client.EnableUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to enable user: %w", err)
	}
	
	// Set privilege level
	var priv uint8
	switch privLevel {
	case "callback":
		priv = 0x01 // Callback
	case "user":
		priv = 0x02 // User
	case "operator":
		priv = 0x03 // Operator
	case "admin":
		priv = 0x04 // Administrator
	case "oem":
		priv = 0x05 // OEM
	default:
		priv = 0x02 // Default to User
	}
	
	req := &ipmi.SetUserAccessRequest{
		ChannelNumber:        1,
		UserID:               userID,
		MaxPrivLevel:         priv,
		EnableLinkAuth:       true,
		EnableIPMIMessaging:  true,
	}
	
	_, err = c.client.SetUserAccess(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to set user access: %w", err)
	}
	
	return nil
}

// UpdateUser updates an existing IPMI user
func (c *Client) UpdateUser(userID uint8, username, password string, privLevel string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Update username if provided
	if username != "" {
		_, err := c.client.SetUsername(ctx, userID, username)
		if err != nil {
			return fmt.Errorf("failed to set username: %w", err)
		}
	}
	
	// Update password if provided
	if password != "" {
		_, err := c.client.SetUserPassword(ctx, userID, password, false)
		if err != nil {
			return fmt.Errorf("failed to set user password: %w", err)
		}
	}
	
	// Update privilege level if provided
	if privLevel != "" {
		var priv uint8
		switch privLevel {
		case "callback":
			priv = 0x01 // Callback
		case "user":
			priv = 0x02 // User
		case "operator":
			priv = 0x03 // Operator
		case "admin":
			priv = 0x04 // Administrator
		case "oem":
			priv = 0x05 // OEM
		default:
			priv = 0x02 // Default to User
		}
		
		req := &ipmi.SetUserAccessRequest{
			ChannelNumber:        1,
			UserID:               userID,
			MaxPrivLevel:         priv,
			EnableLinkAuth:       true,
			EnableIPMIMessaging:  true,
		}
		
		_, err := c.client.SetUserAccess(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to set user access: %w", err)
		}
	}
	
	return nil
}

// DeleteUser deletes an IPMI user
func (c *Client) DeleteUser(userID uint8) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Disable user
	err := c.client.DisableUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to disable user: %w", err)
	}
	
	// Clear username
	_, err = c.client.SetUsername(ctx, userID, "")
	if err != nil {
		return fmt.Errorf("failed to clear username: %w", err)
	}
	
	// Clear password
	_, err = c.client.SetUserPassword(ctx, userID, "", false)
	if err != nil {
		return fmt.Errorf("failed to clear user password: %w", err)
	}
	
	return nil
}

// GetFRUData retrieves FRU (Field Replaceable Unit) data
func (c *Client) GetFRUData() ([]*types.FRUData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	frus, err := c.client.GetFRUs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get FRU data: %w", err)
	}
	
	var result []*types.FRUData
	for _, fru := range frus {
		result = append(result, types.ConvertFRU(fru))
	}
	
	return result, nil
}

// GetSEL retrieves System Event Log entries
func (c *Client) GetSEL() ([]*types.SELEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	entries, err := c.client.GetSELEntries(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get SEL entries: %w", err)
	}
	
	var result []*types.SELEntry
	for _, entry := range entries {
		// For now, we'll just use the record ID since we need to check the structure more carefully
		result = append(result, &types.SELEntry{
			RecordID:  entry.RecordID,
			Type:      fmt.Sprintf("Type %d", entry.RecordType),
			Timestamp: "N/A", // We'll need to extract this from the specific record type
		})
	}
	
	return result, nil
}

// ClearSEL clears the System Event Log
func (c *Client) ClearSEL() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Reserve SEL first
	reservation, err := c.client.ReserveSEL(ctx)
	if err != nil {
		return fmt.Errorf("failed to reserve SEL: %w", err)
	}
	
	// Clear SEL
	_, err = c.client.ClearSEL(ctx, reservation.ReservationID)
	if err != nil {
		return fmt.Errorf("failed to clear SEL: %w", err)
	}
	
	return nil
}

// GetLANConfig retrieves LAN configuration
func (c *Client) GetLANConfig() (*types.LANConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	config, err := c.client.GetLanConfig(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get LAN config: %w", err)
	}
	
	return &types.LANConfig{
		IPAddress:      config.IP.String(),
		MACAddress:     config.MAC.String(),
		SubnetMask:     config.SubnetMask.String(),
		GatewayIP:      config.DefaultGatewayIP.String(),
		GatewayMAC:     config.DefaultGatewayMAC.String(),
		PrimaryDNS:     "", // Not directly available in LanConfig
		SecondaryDNS:   "", // Not directly available in LanConfig
		VLANID:         config.VLANID,
		VLANPriority:   config.VLANPriority,
	}, nil
}

// SetLANConfig sets LAN configuration
func (c *Client) SetLANConfig(config *types.LANConfig) error {
	// This is a complex operation that would require setting multiple parameters
	// For now, we'll just return an error indicating it's not fully implemented
	return fmt.Errorf("setting LAN configuration is not fully implemented")
}