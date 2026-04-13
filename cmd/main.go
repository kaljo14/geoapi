package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	geopulseapi "github.com/neofyis/geopulse/api"
	"github.com/neofyis/geopulse/internal/app"
	"github.com/neofyis/geopulse/internal/config"
	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/handler"
	mw "github.com/neofyis/geopulse/internal/middleware"
	"github.com/neofyis/geopulse/internal/store"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var (
		querier store.Store
		pinger  app.Pinger
		pool    *pgxpool.Pool
	)

	if cfg.DatabaseURL != "" {
		var err error
		pool, err = pgxpool.New(ctx, cfg.DatabaseURL)
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

	a := app.New(querier, pinger, cfg, logger)
	h := handler.NewHandler(a, logger)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(mw.ZapLogger(logger))

	generated.HandlerWithOptions(generated.NewStrictHandler(h, nil), generated.ChiServerOptions{
		BaseRouter: r,
	})
	r.Get("/api/sofiaplan/neighborhoods", h.GetNeighborhoods)
	r.Get("/api/places/export", h.ExportPlaces)
	r.Get("/api/places/export-simple", h.ExportPlacesSimple)
	r.Post("/api/places/import", h.ImportPlaces)
	r.Get("/api/metro/transit-stops", h.GetTransitStops)
	r.Get("/api/adres-locations", h.ListAdresLocations)

	// OpenAPI spec + Swagger UI
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(geopulseapi.Spec)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/openapi.yaml"),
	))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server")
		a.Shutdown() // cancel and wait for background goroutines
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
