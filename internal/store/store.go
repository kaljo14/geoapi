package store

import (
	"context"
	"encoding/json"
)

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
	BulkInsertSofiaplanFeatures(ctx context.Context, sql string, propsJSON, geomJSON []string) error
	TruncateSofiaplanTable(ctx context.Context, tableName string) error
	GetZoningContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetIncomeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetMetroCatchmentContext(ctx context.Context, lng, lat float64) (json.RawMessage, bool, error)
	GetPedestrianContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetBusinessTurnoverContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetPropertyPriceContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetPopulationContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetDevelopmentPotentialContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetNeighborhoodContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetBuildingDensityGeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetBuildingFootprintGeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetResidentialTypologyGeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetUrbanMorphologyGeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetFloodRiskLowContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetFloodRiskMediumContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	GetFloodRiskHighContext(ctx context.Context, lng, lat float64) (json.RawMessage, error)
	ListNeighborhoodNames(ctx context.Context) ([]string, error)
	ListParkingZones(ctx context.Context) ([]ParkingZone, error)
	ParkingZonesGeoJSON(ctx context.Context) (json.RawMessage, error)
}

var _ Store = (*Queries)(nil)
