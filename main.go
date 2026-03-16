package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/elect0/voxium/internal/api"
	"github.com/elect0/voxium/internal/config"
	"github.com/elect0/voxium/internal/database"
	"github.com/elect0/voxium/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config")
	}

	logger, err := logger.New(cfg.App.Env)

	if err != nil {
		log.Fatal("Error initializing logger")
	}
	defer logger.Sync()

	pool, err := database.Init(context.Background(), cfg.Database, logger)
	if err != nil {
		logger.Fatal("Database initialization failed", zap.Error(err))
	}

	defer pool.Close()

	logger.Info("Database ready (Migrated & Connected)")

	logger.Info("Connecting to NATS...", zap.String("url", cfg.NATS.URL))

	// TOOD: Set up NATS

	r := api.NewRouter(api.Config{
		Log: logger,
		DB:  pool,
	})

	addr := fmt.Sprintf(":%d", cfg.App.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(r, &http2.Server{}),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Listen for interrupt signal
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

}
