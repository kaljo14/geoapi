package app

import (
	"context"
	"fmt"

	"github.com/neofyis/geopulse/bff/internal/generated"
)

func (a *App) GetMetroShapes(_ context.Context) (*generated.GeoJSONFeatureCollection, error) {
	// Route geometry (LineStrings) requires GTFS shapes data.
	// Returns empty collection until gtfs_shapes table is populated.
	return &generated.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: []interface{}{},
	}, nil
}

func (a *App) GetMetroStops(ctx context.Context) (*generated.GeoJSONFeatureCollection, error) {
	stops, err := a.store.GetMetroStops(ctx)
	if err != nil {
		return nil, fmt.Errorf("get metro stops: %w", err)
	}

	features := make([]interface{}, 0, len(stops))
	for _, s := range stops {
		features = append(features, map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{s.StopLon.Float64, s.StopLat.Float64},
			},
			"properties": map[string]interface{}{
				"stop_id":   s.StopID,
				"stop_name": s.StopName.String,
			},
		})
	}

	return &generated.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}, nil
}
