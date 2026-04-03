package app

import (
	"context"
	"fmt"

	"github.com/neofyis/geopulse/internal/generated"
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
	stations, err := a.store.GetMetroStations(ctx)
	if err != nil {
		return nil, fmt.Errorf("get metro stations: %w", err)
	}

	features := make([]interface{}, 0, len(stations))
	for _, s := range stations {
		features = append(features, map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{s.StopLon, s.StopLat},
			},
			"properties": map[string]interface{}{
				"stop_id":          s.StopID,
				"stop_name":        s.StopName,
				"line":             s.Line,
				"route_color":      s.RouteColor,
				"route_text_color": s.RouteTextColor,
			},
		})
	}

	return &generated.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}, nil
}

func (a *App) GetTransitStops(ctx context.Context) (*generated.GeoJSONFeatureCollection, error) {
	stops, err := a.store.GetTransitStops(ctx)
	if err != nil {
		return nil, fmt.Errorf("get transit stops: %w", err)
	}

	features := make([]interface{}, 0, len(stops))
	for _, s := range stops {
		features = append(features, map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{s.StopLon, s.StopLat},
			},
			"properties": map[string]interface{}{
				"stop_id":   s.StopID,
				"stop_name": s.StopName,
				"stop_type": s.StopType,
			},
		})
	}

	return &generated.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}, nil
}
