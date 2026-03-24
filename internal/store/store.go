package store

import "context"

// Store extends the sqlc-generated Querier with hand-written raw queries.
type Store interface {
	Querier
	GetSaturation(ctx context.Context, p SaturationParams) (SaturationRow, error)
	GetHeatmap(ctx context.Context, p HeatmapParams) ([]HeatmapRow, error)
}

var _ Store = (*Queries)(nil)
