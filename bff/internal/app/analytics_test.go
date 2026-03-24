package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/neofyis/geopulse/bff/internal/mocks"
	"github.com/neofyis/geopulse/bff/internal/store"
)

// ---------------------------------------------------------------------------
// GetSaturation — score clamping
// ---------------------------------------------------------------------------

func TestGetSaturation_ScoreClampedAt100(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().GetSaturation(mock.Anything, mock.Anything).
		Return(store.SaturationRow{CompetitorCount: 5, DensityPerKm2: 30}, nil)

	result, err := newApp(ms).GetSaturation(context.Background(), 42.7, 23.3, 500, "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Score != 100 {
		t.Errorf("expected score 100, got %v", result.Score)
	}
}

func TestGetSaturation_ScoreProportionalToDensity(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().GetSaturation(mock.Anything, mock.Anything).
		Return(store.SaturationRow{CompetitorCount: 2, DensityPerKm2: 10}, nil)

	result, err := newApp(ms).GetSaturation(context.Background(), 42.7, 23.3, 500, "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Score != 50 {
		t.Errorf("expected score 50, got %v", result.Score)
	}
}

// ---------------------------------------------------------------------------
// GetHeatmap — score clamping
// ---------------------------------------------------------------------------

func TestGetHeatmap_ScoreClampedToZeroWhenNegative(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().GetHeatmap(mock.Anything, mock.Anything).
		Return([]store.HeatmapRow{{Lat: 42.7, Lng: 23.3, AnchorScore: 10, CompPenalty: 50}}, nil)

	tiles, err := newApp(ms).GetHeatmap(context.Background(), 42.68, 23.30, 42.72, 23.35, 0.01, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(tiles) != 1 || tiles[0].Score != 0 {
		t.Errorf("expected score 0 for negative result, got %v", tiles[0].Score)
	}
}

func TestGetHeatmap_ScoreClampedAt100(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().GetHeatmap(mock.Anything, mock.Anything).
		Return([]store.HeatmapRow{{Lat: 42.7, Lng: 23.3, AnchorScore: 200, CompPenalty: 0}}, nil)

	tiles, err := newApp(ms).GetHeatmap(context.Background(), 42.68, 23.30, 42.72, 23.35, 0.01, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(tiles) != 1 || tiles[0].Score != 100 {
		t.Errorf("expected score 100 when clamped, got %v", tiles[0].Score)
	}
}

func TestGetHeatmap_ScoreCalculated(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().GetHeatmap(mock.Anything, mock.Anything).
		Return([]store.HeatmapRow{{Lat: 42.7, Lng: 23.3, AnchorScore: 80, CompPenalty: 30}}, nil)

	tiles, err := newApp(ms).GetHeatmap(context.Background(), 42.68, 23.30, 42.72, 23.35, 0.01, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(tiles) != 1 || tiles[0].Score != 50 {
		t.Errorf("expected score 50, got %v", tiles[0].Score)
	}
}
