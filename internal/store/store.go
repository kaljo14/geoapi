package store

import "context"

// Store extends the sqlc-generated Querier with hand-written raw queries.
type Store interface {
	Querier
	GetSaturation(ctx context.Context, p SaturationParams) (SaturationRow, error)
	GetHeatmap(ctx context.Context, p HeatmapParams) ([]HeatmapRow, error)
	BulkUpsertNodes(ctx context.Context, nodes []OSMNode) error
	BulkUpsertEdges(ctx context.Context, edges []OSMEdge) error
	BulkUpsertPOIs(ctx context.Context, pois []OSMPOI) error
	ComputeWalkScores(ctx context.Context) ([]NodeScore, error)
	UpdateNodeWalkScores(ctx context.Context, scores []NodeScore) error
	PropagateEdgeScores(ctx context.Context) error
}

var _ Store = (*Queries)(nil)
