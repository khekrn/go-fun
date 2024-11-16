package server

import (
	"coding2fun.in/url-shortner/internal/config"
	"coding2fun.in/url-shortner/internal/log"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	db     *gorm.DB
	server *http.Server
}

func NewServer(config *config.Config, db *gorm.DB) *Server {
	gin.SetMode(config.Server.Mode)

	router := gin.Default()

	server := &Server{
		router: router,
		config: config,
		db:     db,
		server: &http.Server{
			Addr:    config.Server.Port,
			Handler: router,
			// Good practice: enforce timeouts for servers you create
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
	}

	// Setup routes
	server.setUp()

	return server
}

func (s *Server) setUp() {
	s.router.GET("/health", s.defaultHandler)
}

func (s *Server) defaultHandler(ctx *gin.Context) {
	sqlDB, err := s.db.DB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to connect to database",
		})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully connected to database",
	})
}

func (s *Server) Run() error {
	log.Info("Starting server",
		zap.String("port", s.config.Server.Port),
		zap.String("mode", s.config.Server.Mode),
	)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info("Initiating graceful shutdown...")

	// Shutdown the HTTP server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Info("Server stopped accepting new requests")

	// Close database connection if needed
	if sqlDB, err := s.db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Error("Error closing database connection", zap.Error(err))
		}
	}

	return nil
}

// StartWithGracefulShutdown starts the server and handles graceful shutdown
func (s *Server) StartWithGracefulShutdown() error {
	// Create channel for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Error channel to collect server errors
	serverErrors := make(chan error, 1)

	// Start server in background
	go func() {
		if err := s.Run(); err != nil {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for quit signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info("Shutdown signal received",
			zap.String("signal", sig.String()),
		)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Perform graceful shutdown
		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
	}

	return nil
}
