package main

import (
	"encoding/json"
	"net/http"
)

// HealthStatus contains information regarding the healthiness of the application
type HealthStatus struct {
	ApplicationRunning bool
}

// Health route
func Health(w http.ResponseWriter, r *http.Request) {
	var hs HealthStatus
	hs.ApplicationRunning = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hs)
}
