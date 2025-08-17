package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mdshack/ipmi-api/pkg/ipmi"
	"github.com/mdshack/ipmi-api/pkg/types"
)

// Server holds the server configuration and components
type Server struct {
	Router      *http.ServeMux
	ipmiPool    *ipmi.ClientPool
	defaultAuth types.IPMIAuth
}

// New creates a new server instance
func New() *Server {
	// Create IPMI client pool
	poolSize, err := strconv.Atoi(getConfig("IPMI_POOL_SIZE", "100"))
	if err != nil {
		poolSize = 100
	}

	timeout, err := time.ParseDuration(getConfig("IPMI_TIMEOUT", "30s"))
	if err != nil {
		timeout = 30 * time.Second
	}

	pool := ipmi.NewClientPool(poolSize, timeout)

	// Set default authentication from environment
	defaultAuth := types.IPMIAuth{
		Host:     os.Getenv("DEFAULT_IPMI_HOST"),
		Port:     623, // Default IPMI port
		Username: os.Getenv("DEFAULT_IPMI_USERNAME"),
		Password: os.Getenv("DEFAULT_IPMI_PASSWORD"),
	}

	// Override port if set in environment
	if portStr := os.Getenv("DEFAULT_IPMI_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			defaultAuth.Port = port
		}
	}

	server := &Server{
		Router:      http.NewServeMux(),
		ipmiPool:    pool,
		defaultAuth: defaultAuth,
	}

	server.setupRoutes()
	return server
}

// setupRoutes registers all the API routes
func (s *Server) setupRoutes() {
	// Health check endpoints
	s.Router.HandleFunc("/health", s.healthHandler)
	s.Router.HandleFunc("/health/ready", s.healthReadyHandler)
	s.Router.HandleFunc("/health/live", s.healthLiveHandler)

	// System information endpoints
	s.Router.HandleFunc("/api/v1/ipmi/system/info", s.withAuth(s.systemInfoHandler))
	s.Router.HandleFunc("/api/v1/ipmi/system/guid", s.withAuth(s.systemGUIDHandler))
	s.Router.HandleFunc("/api/v1/ipmi/system/boot-options", s.withAuth(s.systemBootOptionsHandler))
	s.Router.HandleFunc("/api/v1/ipmi/system/boot-device", s.withAuth(s.systemBootDeviceHandler))

	// Power management endpoints
	s.Router.HandleFunc("/api/v1/ipmi/power/status", s.withAuth(s.powerStatusHandler))
	s.Router.HandleFunc("/api/v1/ipmi/power/on", s.withAuth(s.powerOnHandler))
	s.Router.HandleFunc("/api/v1/ipmi/power/off", s.withAuth(s.powerOffHandler))
	s.Router.HandleFunc("/api/v1/ipmi/power/cycle", s.withAuth(s.powerCycleHandler))
	s.Router.HandleFunc("/api/v1/ipmi/power/reset", s.withAuth(s.powerResetHandler))

	// Sensors & Monitoring endpoints
	s.Router.HandleFunc("/api/v1/ipmi/sensors", s.withAuth(s.sensorsHandler))
	s.Router.HandleFunc("/api/v1/ipmi/sdr", s.withAuth(s.sdrHandler))
	s.Router.HandleFunc("/api/v1/ipmi/sel", s.withAuth(s.selHandler))
	s.Router.HandleFunc("/api/v1/ipmi/sel/clear", s.withAuth(s.selClearHandler))

	// Chassis management endpoints
	s.Router.HandleFunc("/api/v1/ipmi/chassis/status", s.withAuth(s.chassisStatusHandler))
	s.Router.HandleFunc("/api/v1/ipmi/chassis/identify", s.withAuth(s.chassisIdentifyHandler))
	s.Router.HandleFunc("/api/v1/ipmi/chassis/bootdev", s.withAuth(s.chassisBootDeviceHandler))
	s.Router.HandleFunc("/api/v1/ipmi/chassis/bootdev/set", s.withAuth(s.chassisBootDeviceSetHandler))

	// FRU (Field Replaceable Unit) endpoints
	s.Router.HandleFunc("/api/v1/ipmi/fru", s.withAuth(s.fruHandler))
	s.Router.HandleFunc("/api/v1/ipmi/fru/", s.withAuth(s.fruByIDHandler))

	// User management endpoints
	s.Router.HandleFunc("/api/v1/ipmi/users", s.withAuth(s.usersHandler))
	s.Router.HandleFunc("/api/v1/ipmi/users/", s.withAuth(s.userByIDHandler))
	s.Router.HandleFunc("/api/v1/ipmi/users/create", s.withAuth(s.usersCreateHandler))

	// LAN configuration endpoints
	s.Router.HandleFunc("/api/v1/ipmi/lan", s.withAuth(s.lanConfigHandler))
}

// Close closes all IPMI connections
func (s *Server) Close() {
	s.ipmiPool.CloseAll()
}

// getConfig retrieves configuration from environment variables with a default value
func getConfig(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
