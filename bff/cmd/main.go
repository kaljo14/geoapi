package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/neofyis/geopulse/bff/internal/app"
	"github.com/neofyis/geopulse/bff/internal/generated"
	"github.com/neofyis/geopulse/bff/internal/handler"
	"github.com/neofyis/geopulse/bff/internal/store"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var (
		querier store.Store
		pinger  app.Pinger
		pool    *pgxpool.Pool
	)

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		var err error
		pool, err = pgxpool.New(ctx, dbURL)
		if err != nil {
			logger.Fatal("failed to create db pool", zap.Error(err))
		}
		defer pool.Close()
		querier = store.New(pool)
		pinger = pool
		logger.Info("database connected")
	} else {
		logger.Warn("DATABASE_URL not set, running without database")
	}

	a := app.New(querier, pinger, logger)
	h := handler.NewHandler(a, logger)
	strict := generated.NewStrictHandler(h, nil)
	router := generated.Handler(strict)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown failed", zap.Error(err))
		}
	}()

	logger.Info("starting server", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed", zap.Error(err))
	}
	logger.Info("server stopped")
}
