package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/auth"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/logger"
	"go.uber.org/zap"
)

type TokenResponse struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}

// GenerateTokenHandler creates a new API key and JWT token
func GenerateTokenHandler(jwtAuth *auth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Generating new API token")

		w.Header().Set("Content-Type", "application/json")

		// Generate a new API key
		apiKey := uuid.New().String()

		// Generate JWT token
		token, err := jwtAuth.GenerateToken(apiKey)
		if err != nil {
			logger.Error("Failed to generate token",
				zap.Error(err))
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		logger.Info("Token generated successfully",
			zap.String("apiKey", apiKey))

		response := TokenResponse{
			APIKey: apiKey,
			Token:  token,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logrus.Errorf("Failed to encode response: %v", err)
		}
	}
}
