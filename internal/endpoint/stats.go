package endpoint

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Dydhzo/stremthru-proxy/internal/shared"
)

type statsHandler struct{}

// AddBytes implements shared.StatsHandler interface
func (s *statsHandler) AddBytes(bytes int64) {
	AddBytes(bytes)
}

// IncrementConnections implements shared.StatsHandler interface
func (s *statsHandler) IncrementConnections() {
	IncrementConnections()
}

// DecrementConnections implements shared.StatsHandler interface
func (s *statsHandler) DecrementConnections() {
	DecrementConnections()
}

func init() {
	shared.RegisterStatsHandler(&statsHandler{})
}

// StatsData represents bandwidth and connection statistics
type StatsData struct {
	ActiveConnections     int32 `json:"active_connections"`
	SystemNetworkStats    *SystemNetworkStats `json:"system_network,omitempty"`
}

// SystemNetworkStats represents system-wide network statistics
type SystemNetworkStats struct {
	BytesReceivedPerSecond int64 `json:"bytes_received_per_second"`
	BytesSentPerSecond     int64 `json:"bytes_sent_per_second"`
	TotalBytesPerSecond    int64 `json:"total_bytes_per_second"`
}

var (
	activeConnections int32
	lastBytesReceived int64
	lastBytesSent     int64
	lastCheckTime     time.Time
)


// AddBytes tracks bytes transferred for bandwidth statistics
func AddBytes(bytes int64) {
	// Not tracking anymore, using system stats instead
}

// IncrementConnections increases active connection count
func IncrementConnections() {
	atomic.AddInt32(&activeConnections, 1)
}

// DecrementConnections decreases active connection count
func DecrementConnections() {
	atomic.AddInt32(&activeConnections, -1)
}


// getSystemNetworkStats reads network statistics from /proc/net/dev
func getSystemNetworkStats() *SystemNetworkStats {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var totalReceived, totalSent int64

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		iface := strings.TrimSuffix(fields[0], ":")
		// Skip loopback
		if iface == "lo" {
			continue
		}

		received, _ := strconv.ParseInt(fields[1], 10, 64)
		sent, _ := strconv.ParseInt(fields[9], 10, 64)

		totalReceived += received
		totalSent += sent
	}

	now := time.Now()
	timeDiff := now.Sub(lastCheckTime).Seconds()

	if timeDiff < 1 {
		timeDiff = 1 // Avoid division by zero
	}

	var receivedPerSec, sentPerSec int64

	if lastCheckTime.IsZero() {
		// First call, no rate to calculate
		receivedPerSec = 0
		sentPerSec = 0
	} else {
		// Calculate bytes per second
		receivedPerSec = int64(float64(totalReceived - lastBytesReceived) / timeDiff)
		sentPerSec = int64(float64(totalSent - lastBytesSent) / timeDiff)
	}

	// Update last values
	lastBytesReceived = totalReceived
	lastBytesSent = totalSent
	lastCheckTime = now

	if receivedPerSec < 0 {
		receivedPerSec = 0
	}
	if sentPerSec < 0 {
		sentPerSec = 0
	}

	return &SystemNetworkStats{
		BytesReceivedPerSecond: receivedPerSec,
		BytesSentPerSecond:     sentPerSec,
		TotalBytesPerSecond:    receivedPerSec + sentPerSec,
	}
}

// handleStats provides real-time bandwidth statistics endpoint
func handleStats(w http.ResponseWriter, r *http.Request) {
	if !shared.IsMethod(r, http.MethodGet) {
		shared.ErrorMethodNotAllowed(r).Send(w, r)
		return
	}

	// Require auth - block completely if not authorized
	isAuthorized, _, _ := getProxyAuthorization(r, true)
	if !isAuthorized {
		w.Header().Add("WWW-Authenticate", "Basic")
		shared.ErrorForbidden(r).Send(w, r)
		return
	}

	// Calculate stats
	connections := atomic.LoadInt32(&activeConnections)

	stats := &StatsData{
		ActiveConnections:     connections,
		SystemNetworkStats:    getSystemNetworkStats(),
	}

	shared.SendResponse(w, r, 200, stats, nil)
}

// AddStatsEndpoint registers statistics HTTP endpoint
func AddStatsEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/v0/stats", handleStats)
}