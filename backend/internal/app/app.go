package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/backend/internal/store"
)

// Pinger abstracts the DB ping check.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Service is the business layer interface consumed by the handler.
type Service interface {
	StartScraper(ctx context.Context) error
	StartEnricher(ctx context.Context) error
	ExportCSV(ctx context.Context, w io.Writer, simple bool) error
	ImportCSV(ctx context.Context, r io.Reader) (int, error)
	Ready(ctx context.Context) error
}

// App implements Service.
type App struct {
	store  store.Querier
	pinger Pinger
	apiKey string
	logger *zap.Logger
}

func New(q store.Querier, p Pinger, logger *zap.Logger) *App {
	return &App{
		store:  q,
		pinger: p,
		apiKey: os.Getenv("GOOGLE_API_KEY"),
		logger: logger,
	}
}

func (a *App) Ready(ctx context.Context) error {
	if a.pinger == nil {
		return fmt.Errorf("database not configured")
	}
	return a.pinger.Ping(ctx)
}
