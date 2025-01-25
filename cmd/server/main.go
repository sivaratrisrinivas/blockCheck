package main

import (
	"context"
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
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/auth"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/ethereum"
	"github.com/sivaratrisrinivas/web3/blockCheck/pkg/handlers"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Debugf("Loaded configuration: %+v", cfg)

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

	log.Debugf("Creating Ethereum validator with config: %+v", ethConfig)

	ethValidator, err := factory.Create("ethereum", ethConfig)
	if err != nil {
		log.Fatalf("Failed to create Ethereum validator: %v", err)
	}

	if err := registry.Register(ethValidator); err != nil {
		log.Fatalf("Failed to register Ethereum validator instance: %v", err)
	}

	// Initialize JWT auth
	jwtAuth := auth.NewJWTAuth(cfg.JWT.SecretKey, cfg.JWT.Duration)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(cfg.ENS.TimeoutSeconds) * time.Second))

	// Public routes
	r.Get("/health", handlers.HealthCheckHandler)
	r.Post("/v1/token", handlers.GenerateTokenHandler(jwtAuth))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth.Middleware)
		r.Get("/v1/validate/{address}", handlers.ValidateAddressHandler(ethValidator))
		r.Get("/v1/resolveEns/{name}", handlers.ResolveENSHandler(ethValidator))
		r.Get("/v1/isContract/{address}", handlers.IsContractHandler(ethValidator))
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
