package handlers

import (
	"net/http"
)

// HealthCheckHandler returns 200 OK if the service is healthy
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
