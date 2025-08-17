package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mdshack/ipmi-api/pkg/ipmi"
	"github.com/mdshack/ipmi-api/pkg/types"
)

// withAuth is a middleware that handles IPMI authentication
func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Resolve authentication credentials
		auth, err := s.resolveAuth(r.Header)
		if err != nil {
			s.sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get or create IPMI client
		client, err := s.ipmiPool.GetClient(auth.Host, auth.Port, auth.Username, auth.Password)
		if err != nil {
			s.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to connect to IPMI: %v", err))
			return
		}

		// Wrap the client
		ipmiClient := ipmi.NewClient(client)

		// Add client to request context
		ctx := contextWithIPMIClient(r.Context(), ipmiClient)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// resolveAuth resolves IPMI authentication credentials from headers or environment
func (s *Server) resolveAuth(headers http.Header) (*types.IPMIAuth, error) {
	auth := &types.IPMIAuth{
		Host:     getHeaderOrDefault(headers, "X-IPMI-Host", s.defaultAuth.Host),
		Username: getHeaderOrDefault(headers, "X-IPMI-Username", s.defaultAuth.Username),
		Password: getHeaderOrDefault(headers, "X-IPMI-Password", s.defaultAuth.Password),
		Port:     s.defaultAuth.Port, // Default port
	}

	// Port handling with fallback
	if portHeader := headers.Get("X-IPMI-Port"); portHeader != "" {
		port, err := strconv.Atoi(portHeader)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %v", err)
		}
		auth.Port = port
	} else if s.defaultAuth.Port != 0 {
		auth.Port = s.defaultAuth.Port
	} else {
		auth.Port = 623 // Default IPMI port
	}

	// Validate required fields
	if auth.Host == "" {
		return nil, fmt.Errorf("IPMI host required (header X-IPMI-Host or environment variable DEFAULT_IPMI_HOST)")
	}
	if auth.Username == "" {
		return nil, fmt.Errorf("IPMI username required (header X-IPMI-Username or environment variable DEFAULT_IPMI_USERNAME)")
	}
	if auth.Password == "" {
		return nil, fmt.Errorf("IPMI password required (header X-IPMI-Password or environment variable DEFAULT_IPMI_PASSWORD)")
	}

	return auth, nil
}

// getHeaderOrDefault returns the header value if present, otherwise returns the default
func getHeaderOrDefault(headers http.Header, headerName, defaultValue string) string {
	if value := headers.Get(headerName); value != "" {
		return value
	}
	return defaultValue
}

// sendErrorResponse sends a standardized error response
func (s *Server) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := types.APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// sendSuccessResponse sends a standardized success response
func (s *Server) sendSuccessResponse(w http.ResponseWriter, data interface{}, executionTime time.Duration) {
	response := types.APIResponse{
		Success:       true,
		Data:          data,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		ExecutionTime: executionTime.Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
