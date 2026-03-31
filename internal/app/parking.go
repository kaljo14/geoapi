package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/neofyis/geopulse/internal/generated"
)

func (a *App) GetParkingZones(ctx context.Context) (*generated.GeoJSONFeatureCollection, error) {
	raw, err := a.store.ParkingZonesGeoJSON(ctx)
	if err != nil {
		return nil, fmt.Errorf("get parking zones: %w", err)
	}

	var fc generated.GeoJSONFeatureCollection
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, fmt.Errorf("unmarshal parking zones geojson: %w", err)
	}
	return &fc, nil
}
