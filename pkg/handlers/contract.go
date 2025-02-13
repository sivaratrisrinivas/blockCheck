package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
)

type ContractResponse struct {
	Address    string `json:"address"`
	IsContract bool   `json:"isContract"`
	Error      string `json:"error,omitempty"`
}

// IsContractHandler handles contract detection requests
func IsContractHandler(validator chain.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := chi.URLParam(r, "address")
		if address == "" {
			http.Error(w, "Address parameter is required", http.StatusBadRequest)
			return
		}

		isContract, err := validator.IsContract(r.Context(), address)
		response := ContractResponse{
			Address: address,
		}

		if err != nil {
			response.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
		} else {
			response.IsContract = isContract
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
		}
	}
}
