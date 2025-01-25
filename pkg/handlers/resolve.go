package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
)

type ResolveResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Error   string `json:"error,omitempty"`
}

// ResolveENSHandler handles ENS name resolution requests
func ResolveENSHandler(validator chain.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		name := chi.URLParam(r, "name")
		if name == "" {
			http.Error(w, "Name parameter is required", http.StatusBadRequest)
			return
		}

		address, err := validator.ResolveENS(name)
		response := ResolveResponse{
			Name: name,
		}

		if err != nil {
			response.Error = err.Error()
			response.Address = "0x0000000000000000000000000000000000000000"
		} else {
			response.Address = address
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
		}
	}
}
