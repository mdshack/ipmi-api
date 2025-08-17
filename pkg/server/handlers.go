package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mdshack/ipmi-api/pkg/types"
)

// healthHandler handles the basic health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC(),
	}

	s.sendSuccessResponse(w, response, 0)
}

// healthReadyHandler handles the readiness probe
func (s *Server) healthReadyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, we're always ready
	// In a production environment, you might check database connections, etc.
	response := map[string]interface{}{
		"status": "ready",
		"time":   time.Now().UTC(),
	}

	s.sendSuccessResponse(w, response, 0)
}

// healthLiveHandler handles the liveness probe
func (s *Server) healthLiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, we're always alive
	// In a production environment, you might check for deadlocks, etc.
	response := map[string]interface{}{
		"status": "alive",
		"time":   time.Now().UTC(),
	}

	s.sendSuccessResponse(w, response, 0)
}

// systemInfoHandler handles getting system information
func (s *Server) systemInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	info, err := client.GetSystemInfo()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get system info: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, info, executionTime)
}

// systemGUIDHandler handles getting system GUID
func (s *Server) systemGUIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	guid, err := client.GetSystemGUID()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get system GUID: %v", err))
		return
	}

	response := map[string]string{
		"guid": guid,
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// systemBootOptionsHandler handles getting boot device options
func (s *Server) systemBootOptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This needs to be implemented
	s.sendErrorResponse(w, http.StatusNotImplemented, "Not implemented")
	return
}

// systemBootDeviceHandler handles setting boot device
func (s *Server) systemBootDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	// Parse request body
	var req struct {
		Device string `json:"device"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
		return
	}

	if err := client.SetBootDevice(req.Device); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set boot device: %v", err))
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("Boot device set to %s successfully", req.Device),
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// powerStatusHandler handles getting power status
func (s *Server) powerStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	status, err := client.GetPowerStatus()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get power status: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, status, executionTime)
}

// powerOnHandler handles powering on the system
func (s *Server) powerOnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	if err := client.PowerOn(); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to power on: %v", err))
		return
	}

	response := map[string]string{
		"message": "Power on command sent successfully",
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// powerOffHandler handles powering off the system
func (s *Server) powerOffHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	if err := client.PowerOff(); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to power off: %v", err))
		return
	}

	response := map[string]string{
		"message": "Power off command sent successfully",
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// powerCycleHandler handles power cycling the system
func (s *Server) powerCycleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	if err := client.PowerCycle(); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to power cycle: %v", err))
		return
	}

	response := map[string]string{
		"message": "Power cycle command sent successfully",
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// powerResetHandler handles resetting the system
func (s *Server) powerResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	if err := client.PowerReset(); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to reset: %v", err))
		return
	}

	response := map[string]string{
		"message": "Power reset command sent successfully",
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// sensorsHandler handles getting sensor data
func (s *Server) sensorsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	sensors, err := client.GetSensors()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get sensors: %v", err))
		return
	}

	// Convert sensors to a more API-friendly format
	sensorData := make([]map[string]interface{}, len(sensors))
	for i, sensor := range sensors {
		sensorData[i] = map[string]interface{}{
			"id":              sensor.Number,
			"name":            sensor.Name,
			"type":            sensor.SensorType.String(),
			"reading":         sensor.ReadingStr(),
			"unit":            sensor.SensorUnit.BaseUnit.String(),
			"status":          sensor.Status(),
			"has_analog":      sensor.HasAnalogReading,
			"entity_id":       sensor.EntityID,
			"entity_instance": sensor.EntityInstance,
		}
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, sensorData, executionTime)
}

// sdrHandler handles getting sensor data records
func (s *Server) sdrHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This needs to be implemented
	s.sendErrorResponse(w, http.StatusNotImplemented, "Not implemented")
	return
}

// selHandler handles getting system event log
func (s *Server) selHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	entries, err := client.GetSEL()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get SEL: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, entries, executionTime)
}

// selClearHandler handles clearing the system event log
func (s *Server) selClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	if err := client.ClearSEL(); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to clear SEL: %v", err))
		return
	}

	response := map[string]string{
		"message": "System event log cleared successfully",
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// chassisStatusHandler handles getting chassis status
func (s *Server) chassisStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	status, err := client.GetChassisStatus()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get chassis status: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, status, executionTime)
}

// chassisIdentifyHandler handles chassis identify LED
func (s *Server) chassisIdentifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This needs to be implemented
	s.sendErrorResponse(w, http.StatusNotImplemented, "Not implemented")
	return
}

// chassisBootDeviceHandler handles getting chassis boot device
func (s *Server) chassisBootDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	device, err := client.GetBootDevice()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get boot device: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, device, executionTime)
}

// chassisBootDeviceSetHandler handles setting chassis boot device
func (s *Server) chassisBootDeviceSetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	// Parse request body
	var req struct {
		Device string `json:"device"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
		return
	}

	if err := client.SetBootDevice(req.Device); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set boot device: %v", err))
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("Boot device set to %s successfully", req.Device),
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// fruHandler handles getting all FRU data
func (s *Server) fruHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	frus, err := client.GetFRUData()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get FRU data: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, frus, executionTime)
}

// fruByIDHandler handles getting specific FRU data
func (s *Server) fruByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This needs to be implemented
	s.sendErrorResponse(w, http.StatusNotImplemented, "Not implemented")
	return
}

// usersHandler handles listing users
func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	users, err := client.GetUsers()
	if err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get users: %v", err))
		return
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, users, executionTime)
}

// userByIDHandler handles getting, updating, or deleting a specific user
func (s *Server) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	// Extract user ID from URL path
	userIDStr := r.URL.Path[len("/api/v1/ipmi/users/"):]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 1 || userID > 15 {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID (must be between 1 and 15)")
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := client.GetUser(uint8(userID))
		if err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %v", err))
			return
		}

		executionTime := time.Since(startTime)
		s.sendSuccessResponse(w, user, executionTime)

	case http.MethodPut:
		// Parse request body
		var req struct {
			Username  string `json:"username,omitempty"`
			Password  string `json:"password,omitempty"`
			PrivLevel string `json:"priv_level,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
			return
		}

		if err := client.UpdateUser(uint8(userID), req.Username, req.Password, req.PrivLevel); err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
			return
		}

		response := map[string]string{
			"message": fmt.Sprintf("User %d updated successfully", userID),
		}

		executionTime := time.Since(startTime)
		s.sendSuccessResponse(w, response, executionTime)

	case http.MethodDelete:
		if err := client.DeleteUser(uint8(userID)); err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete user: %v", err))
			return
		}

		response := map[string]string{
			"message": fmt.Sprintf("User %d deleted successfully", userID),
		}

		executionTime := time.Since(startTime)
		s.sendSuccessResponse(w, response, executionTime)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// usersCreateHandler handles creating a new user
func (s *Server) usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	// Parse request body
	var req struct {
		UserID    int    `json:"user_id"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		PrivLevel string `json:"priv_level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
		return
	}

	// Validate user ID
	if req.UserID < 1 || req.UserID > 15 {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID (must be between 1 and 15)")
		return
	}

	if err := client.CreateUser(uint8(req.UserID), req.Username, req.Password, req.PrivLevel); err != nil {
		s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("User %d created successfully", req.UserID),
	}

	executionTime := time.Since(startTime)
	s.sendSuccessResponse(w, response, executionTime)
}

// lanConfigHandler handles getting and setting LAN configuration
func (s *Server) lanConfigHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	client, ok := ipmiClientFromContext(r.Context())
	if !ok {
		s.sendErrorResponse(w, http.StatusInternalServerError, "IPMI client not found in context")
		return
	}

	switch r.Method {
	case http.MethodGet:
		config, err := client.GetLANConfig()
		if err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get LAN config: %v", err))
			return
		}

		executionTime := time.Since(startTime)
		s.sendSuccessResponse(w, config, executionTime)

	case http.MethodPost:
		// Parse request body
		var config types.LANConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			s.sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse request: %v", err))
			return
		}

		if err := client.SetLANConfig(&config); err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set LAN config: %v", err))
			return
		}

		response := map[string]string{
			"message": "LAN configuration updated successfully",
		}

		executionTime := time.Since(startTime)
		s.sendSuccessResponse(w, response, executionTime)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// contextWithTimeout creates a context with timeout from the request context
func contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	// Default timeout of 30 seconds
	return context.WithTimeout(ctx, 30*time.Second)
}
