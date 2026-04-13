package app

import (
	"context"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/config"
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
	GetTransitStops(ctx context.Context) (*generated.GeoJSONFeatureCollection, error)
	GetSaturation(ctx context.Context, lat, lng, radius float64, category string) (*generated.SaturationResult, error)
	GetHeatmap(ctx context.Context, minLat, minLng, maxLat, maxLng, cellSize float64, category string) ([]generated.HeatmapTile, error)
	StartScraper(ctx context.Context) error
	StartEnricher(ctx context.Context) error
	ImportOSMNetwork(ctx context.Context) error
	ImportOSMPOIs(ctx context.Context) error
	ImportSofiaplan(ctx context.Context, layer string) error
	GetSofiaplanContext(ctx context.Context, lat, lng float64) (*generated.LocationContext, error)
	ListNeighborhoods(ctx context.Context) ([]string, error)
	GetParkingZones(ctx context.Context) (*generated.GeoJSONFeatureCollection, error)
	ListAdresLocations(ctx context.Context) ([]AdresLocation, error)
	ListRetailListings(ctx context.Context) ([]generated.RetailListing, error)
	GetRetailListing(ctx context.Context, id string) (*generated.RetailListing, error)
	CreateRetailListing(ctx context.Context, req generated.CreateRetailListingRequest) (*generated.RetailListing, error)
	UpdateRetailListing(ctx context.Context, id string, req generated.UpdateRetailListingRequest) (*generated.RetailListing, error)
	DeleteRetailListing(ctx context.Context, id string) error
	ExportCSV(ctx context.Context, w io.Writer, simple bool) error
	ImportCSV(ctx context.Context, r io.Reader) (int, error)
	Ready(ctx context.Context) error
}

// App implements Service.
type App struct {
	store  store.Store
	pinger Pinger
	cfg    *config.Config
	logger *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(s store.Store, p Pinger, cfg *config.Config, logger *zap.Logger) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		store:  s,
		pinger: p,
		cfg:    cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Shutdown cancels all background goroutines and waits for them to finish.
func (a *App) Shutdown() {
	a.cancel()
	a.wg.Wait()
}

func (a *App) Ready(ctx context.Context) error {
	if a.pinger == nil {
		return fmt.Errorf("database not configured")
	}
	return a.pinger.Ping(ctx)
}
