package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/sivaratrisrinivas/web3/blockCheck/config"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/ethereum"
)

var log = logrus.New()

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize validator factory and registry
	factory := chain.NewFactory()
	registry := chain.NewRegistry()

	// Register Ethereum validator
	if err := factory.Register("ethereum", ethereum.NewValidator); err != nil {
		log.Fatalf("Failed to register Ethereum validator: %v", err)
	}

	// Create and register Ethereum validator instance
	ethConfig := map[string]interface{}{
		"provider_url":   cfg.ENS.ProviderURL,
		"cache_duration": int64(cfg.Cache.TTL.Seconds()),
	}

	ethValidator, err := factory.Create("ethereum", ethConfig)
	if err != nil {
		log.Fatalf("Failed to create Ethereum validator: %v", err)
	}

	if err := registry.Register(ethValidator); err != nil {
		log.Fatalf("Failed to register Ethereum validator instance: %v", err)
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(cfg.ENS.TimeoutSeconds) * time.Second))

	// Routes
	r.Get("/health", handleHealth)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/validate/{address}", handleValidateAddress(registry))
		r.Get("/resolveEns/{name}", handleResolveENS(registry))
		r.Get("/isContract/{address}", handleIsContract(registry))
	})

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Infof("Server started on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited properly")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleValidateAddress(registry *chain.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		validator, err := registry.Get("ethereum") // Default to Ethereum for now
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isValid := validator.ValidateAddress(address)
		response := map[string]interface{}{
			"address": address,
			"isValid": isValid,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleResolveENS(registry *chain.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		validator, err := registry.Get("ethereum")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		address, err := validator.ResolveName(r.Context(), name)
		response := map[string]interface{}{
			"name": name,
		}

		if err != nil {
			response["error"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
		} else {
			response["address"] = address
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleIsContract(registry *chain.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		validator, err := registry.Get("ethereum")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isContract, err := validator.IsContract(r.Context(), address)
		response := map[string]interface{}{
			"address": address,
		}

		if err != nil {
			response["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
		} else {
			response["isContract"] = isContract
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
