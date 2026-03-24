package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/store"
)

// Pinger abstracts the DB ping check so it can be mocked in tests.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Service is the business layer interface consumed by the handler.
type Service interface {
	ListPlaces(ctx context.Context, category, tag string) ([]generated.Place, error)
	GetPlace(ctx context.Context, placeID string) (*generated.Place, error)
	CreatePlace(ctx context.Context, req generated.CreatePlaceRequest) (*generated.Place, error)
	UpdatePlace(ctx context.Context, placeID string, req generated.UpdatePlaceRequest) (*generated.Place, error)
	DeletePlace(ctx context.Context, placeID string) error
	GetMetroShapes(ctx context.Context) (*generated.GeoJSONFeatureCollection, error)
	GetMetroStops(ctx context.Context) (*generated.GeoJSONFeatureCollection, error)
	GetSaturation(ctx context.Context, lat, lng, radius float64, category string) (*generated.SaturationResult, error)
	GetHeatmap(ctx context.Context, minLat, minLng, maxLat, maxLng, cellSize float64, category string) ([]generated.HeatmapTile, error)
	StartScraper(ctx context.Context) error
	StartEnricher(ctx context.Context) error
	ExportCSV(ctx context.Context, w io.Writer, simple bool) error
	ImportCSV(ctx context.Context, r io.Reader) (int, error)
	Ready(ctx context.Context) error
}

// App implements Service.
type App struct {
	store  store.Store
	pinger Pinger
	apiKey string
	logger *zap.Logger
}

func New(s store.Store, p Pinger, logger *zap.Logger) *App {
	return &App{
		store:  s,
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
