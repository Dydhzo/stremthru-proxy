package endpoint

import (
	"net/http"
	"time"

	"github.com/Dydhzo/stremthru-proxy/internal/config"
	"github.com/Dydhzo/stremthru-proxy/internal/shared"
)

// HealthData represents basic health check response
type HealthData struct {
	Status string `json:"status"`
}

// HealthDebugData represents detailed health check response
type HealthDebugData struct {
	Time    string                 `json:"time"`
	Version string                 `json:"version"`
	User    *HealthDebugUserData   `json:"user,omitempty"`
	IP      *HealthDebugIPData     `json:"ip,omitempty"`
}

// HealthDebugUserData represents user authentication info in debug response
type HealthDebugUserData struct {
	Name string `json:"name"`
}

// HealthDebugIPData represents client IP information
type HealthDebugIPData struct {
	Machine  string            `json:"machine"`
	Tunnel   map[string]string `json:"tunnel"`
	Exposed  map[string]string `json:"exposed"`
}

// handleHealth provides basic health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := &HealthData{}
	health.Status = "ok"
	shared.SendResponse(w, r, 200, health, nil)
}

// handleHealthDebug provides detailed health information with authentication
func handleHealthDebug(w http.ResponseWriter, r *http.Request) {
	if !shared.IsMethod(r, http.MethodGet) {
		shared.ErrorMethodNotAllowed(r).Send(w, r)
		return
	}

	debug := &HealthDebugData{
		Time:    time.Now().Format(time.RFC3339),
		Version: config.Version,
	}

	// Check auth exactly like original
	isAuthorized, user, _ := getProxyAuthorization(r, true)

	// ONLY provide data if authorized (exactly like original)
	if isAuthorized && user != "" {
		debug.User = &HealthDebugUserData{
			Name: user,
		}

		// Get Machine IP like original (only when authorized)
		machineIP := config.IP.GetMachineIP()

		// Simple IP structure for clean proxy (like original logic)
		debug.IP = &HealthDebugIPData{
			Machine: machineIP,
			Tunnel:  map[string]string{}, // Empty for clean proxy
			Exposed: map[string]string{
				"*": machineIP, // Same as original logic
			},
		}
	}

	// Send directly - shared.SendResponse already wraps in "data"
	shared.SendResponse(w, r, 200, debug, nil)
}

// AddHealthEndpoints registers health check HTTP endpoints
func AddHealthEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/v0/health", handleHealth)
	mux.HandleFunc("/v0/health/__debug__", handleHealthDebug)
}
