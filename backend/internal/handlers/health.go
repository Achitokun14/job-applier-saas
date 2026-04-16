package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ComponentStatus describes the health of a single component.
type ComponentStatus struct {
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

// HealthResponse is the structured JSON returned by the health endpoint.
type HealthResponse struct {
	Status     string                     `json:"status"`
	Version    string                     `json:"version"`
	Components map[string]ComponentStatus `json:"components"`
}

const healthCheckTimeout = 5 * time.Second

// HealthCheck returns an http.HandlerFunc that probes the database, Redis,
// and Python service and returns a structured health report.
func HealthCheck(db *gorm.DB, redisAddr string, pythonServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Status:  "healthy",
			Version: "1.0.0",
			Components: map[string]ComponentStatus{
				"database":       checkDatabase(db),
				"redis":          checkRedis(redisAddr),
				"python_service": checkPythonService(pythonServiceURL),
			},
		}

		// If any component is down, the overall status is degraded.
		for _, c := range resp.Components {
			if c.Status != "up" {
				resp.Status = "degraded"
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if resp.Status != "healthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func checkDatabase(db *gorm.DB) ComponentStatus {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	start := time.Now()
	sqlDB, err := db.DB()
	if err != nil {
		return ComponentStatus{Status: "down", LatencyMs: time.Since(start).Milliseconds(), Error: err.Error()}
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return ComponentStatus{Status: "down", LatencyMs: time.Since(start).Milliseconds(), Error: err.Error()}
	}
	return ComponentStatus{Status: "up", LatencyMs: time.Since(start).Milliseconds()}
}

func checkRedis(addr string) ComponentStatus {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	start := time.Now()
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return ComponentStatus{Status: "down", LatencyMs: time.Since(start).Milliseconds(), Error: err.Error()}
	}
	return ComponentStatus{Status: "up", LatencyMs: time.Since(start).Milliseconds()}
}

func checkPythonService(baseURL string) ComponentStatus {
	start := time.Now()
	client := &http.Client{Timeout: healthCheckTimeout}

	resp, err := client.Get(fmt.Sprintf("%s/health", baseURL))
	if err != nil {
		return ComponentStatus{Status: "down", LatencyMs: time.Since(start).Milliseconds(), Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ComponentStatus{
			Status:    "down",
			LatencyMs: time.Since(start).Milliseconds(),
			Error:     fmt.Sprintf("unexpected status %d", resp.StatusCode),
		}
	}
	return ComponentStatus{Status: "up", LatencyMs: time.Since(start).Milliseconds()}
}
