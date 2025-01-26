package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
)

type ValidateResponse struct {
	Address string `json:"address"`
	IsValid bool   `json:"isValid"`
}

// ValidateAddressHandler handles Ethereum address validation requests
func ValidateAddressHandler(validator chain.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := chi.URLParam(r, "address")
		if address == "" {
			http.Error(w, "Address parameter is required", http.StatusBadRequest)
			return
		}

		// Validate using EIP-55 checksum
		isValid := validator.IsChecksumAddress(address)
		logrus.WithFields(logrus.Fields{
			"address": address,
			"isValid": isValid,
		}).Debug("EIP-55 validation result")

		resp := ValidateResponse{
			Address: address,
			IsValid: isValid,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
		}
	}
}
