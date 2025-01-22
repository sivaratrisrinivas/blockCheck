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
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/ens"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator"
	"github.com/sivaratrisrinivas/web3/blockCheck/pkg/models"
)

var (
	ensResolver *ens.Resolver
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize ENS resolver
	ensResolver, err = ens.NewResolver(cfg.ENS.ProviderURL, time.Duration(cfg.Cache.TTL))
	if err != nil {
		log.Fatalf("Failed to initialize ENS resolver: %v", err)
	}
	defer ensResolver.Close()

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(cfg.ENS.TimeoutSeconds) * time.Second))

	// Routes
	r.Get("/health", handleHealth)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/validate/{address}", handleValidateAddress)
		r.Get("/resolveEns/{name}", handleResolveENS)
		r.Get("/addressType/{address}", handleAddressType)
	})

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process lifecycle management
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel() // Ensure context resources are cleaned up

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Printf("Server is starting on %s", addr)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

// Handler functions
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleValidateAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

	response := &models.AddressValidationResponse{
		Address: address,
	}

	// Basic format validation
	response.IsValid = validator.IsValidAddress(address)
	if !response.IsValid {
		response.Error = "Invalid Ethereum address format"
		writeJSON(w, http.StatusBadRequest, response)
		return
	}

	// Checksum validation
	response.HasValidChecksum = validator.IsChecksumAddress(address)

	// Get checksum address
	checksumAddr, err := validator.ToChecksumAddress(address)
	if err != nil {
		response.Error = fmt.Sprintf("Error converting to checksum address: %v", err)
		writeJSON(w, http.StatusInternalServerError, response)
		return
	}
	response.ChecksumAddress = checksumAddr

	writeJSON(w, http.StatusOK, response)
}

func handleResolveENS(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Resolve ENS name
	result, err := ensResolver.Resolve(ctx, name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to resolve ENS name: %v", err),
		})
		return
	}

	if result.Error != "" {
		writeJSON(w, http.StatusNotFound, result)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAddressType(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement in next step
	w.WriteHeader(http.StatusNotImplemented)
}
