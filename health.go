package main

import (
	"encoding/json"
	"net/http"
)

type HealthStatus struct {
	ApplicationRunning bool
}

func Health(w http.ResponseWriter, r *http.Request) {
	var hs HealthStatus
	hs.ApplicationRunning = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hs)
}
