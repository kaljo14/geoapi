package app

import (
	"context"
	"fmt"

	"github.com/neofyis/geopulse/bff/internal/generated"
	"github.com/neofyis/geopulse/bff/internal/store"
)

func (a *App) GetSaturation(ctx context.Context, lat, lng, radius float64, category string) (*generated.SaturationResult, error) {
	row, err := a.store.GetSaturation(ctx, store.SaturationParams{
		Lat:      lat,
		Lng:      lng,
		Radius:   radius,
		Category: category,
	})
	if err != nil {
		return nil, fmt.Errorf("get saturation: %w", err)
	}

	score := row.DensityPerKm2 * 5
	if score > 100 {
		score = 100
	}
	return &generated.SaturationResult{
		Lat:             lat,
		Lng:             lng,
		RadiusMeters:    radius,
		Category:        &category,
		CompetitorCount: row.CompetitorCount,
		DensityPerKm2:   row.DensityPerKm2,
		Score:           score,
	}, nil
}

func (a *App) GetHeatmap(ctx context.Context, minLat, minLng, maxLat, maxLng, cellSize float64, category string) ([]generated.HeatmapTile, error) {
	rows, err := a.store.GetHeatmap(ctx, store.HeatmapParams{
		MinLat:   minLat,
		MinLng:   minLng,
		MaxLat:   maxLat,
		MaxLng:   maxLng,
		CellSize: cellSize,
		Category: category,
	})
	if err != nil {
		return nil, fmt.Errorf("get heatmap: %w", err)
	}

	tiles := make([]generated.HeatmapTile, 0, len(rows))
	for _, r := range rows {
		score := r.AnchorScore - r.CompPenalty
		if score < 0 {
			score = 0
		}
		if score > 100 {
			score = 100
		}
		tiles = append(tiles, generated.HeatmapTile{
			Lat:         r.Lat,
			Lng:         r.Lng,
			AnchorScore: r.AnchorScore,
			CompPenalty: r.CompPenalty,
			Score:       score,
		})
	}
	return tiles, nil
}
