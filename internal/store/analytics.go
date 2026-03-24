package store

import (
	"context"
	"fmt"
)

type SaturationParams struct {
	Lat      float64
	Lng      float64
	Radius   float64
	Category string
}

type SaturationRow struct {
	CompetitorCount int32
	DensityPerKm2   float64
}

type HeatmapParams struct {
	MinLat   float64
	MinLng   float64
	MaxLat   float64
	MaxLng   float64
	CellSize float64
	Category string
}

type HeatmapRow struct {
	Lat         float64
	Lng         float64
	AnchorScore float64
	CompPenalty float64
}

const saturationSQL = `
SELECT
    COUNT(*)::int AS competitor_count,
    COALESCE(
        COUNT(*)::float8 / NULLIF(PI() * POWER($3::float8 / 1000.0, 2), 0),
        0
    ) AS density_per_km2
FROM places
WHERE location IS NOT NULL
  AND ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography, $3)
  AND business_status = 'OPERATIONAL'
  AND ($4 = '' OR category = $4)
`

const heatmapSQL = `
SELECT
    ST_Y(g.geom)::float8 AS lat,
    ST_X(g.geom)::float8 AS lng,
    COALESCE((
        SELECT SUM(100.0 * (1.0 - ST_Distance(s.geom::geography, g.geom::geography) / 1000.0))
        FROM gtfs_stops s
        WHERE s.geom IS NOT NULL
          AND ST_DWithin(s.geom::geography, g.geom::geography, 1000)
    ), 0)::float8 AS anchor_score,
    COALESCE((
        SELECT COUNT(*) * 50.0
        FROM places p
        WHERE p.location IS NOT NULL
          AND p.business_status = 'OPERATIONAL'
          AND ($6::text = '' OR p.category = $6::text)
          AND ST_DWithin(p.location::geography, g.geom::geography, 800)
    ), 0)::float8 AS comp_penalty
FROM (
    SELECT ST_SetSRID(ST_MakePoint(
        $2::float8 + j * $5::float8 + $5::float8 / 2.0,
        $1::float8 + i * $5::float8 + $5::float8 / 2.0
    ), 4326) AS geom
    FROM generate_series(0, CEIL(($3::float8 - $1::float8) / $5::float8)::int - 1) AS i,
         generate_series(0, CEIL(($4::float8 - $2::float8) / $5::float8)::int - 1) AS j
) g
`

func (q *Queries) GetSaturation(ctx context.Context, p SaturationParams) (SaturationRow, error) {
	var row SaturationRow
	err := q.db.QueryRow(ctx, saturationSQL, p.Lat, p.Lng, p.Radius, p.Category).
		Scan(&row.CompetitorCount, &row.DensityPerKm2)
	if err != nil {
		return SaturationRow{}, fmt.Errorf("get saturation: %w", err)
	}
	return row, nil
}

func (q *Queries) GetHeatmap(ctx context.Context, p HeatmapParams) ([]HeatmapRow, error) {
	rows, err := q.db.Query(ctx, heatmapSQL, p.MinLat, p.MinLng, p.MaxLat, p.MaxLng, p.CellSize, p.Category)
	if err != nil {
		return nil, fmt.Errorf("get heatmap: %w", err)
	}
	defer rows.Close()

	var result []HeatmapRow
	for rows.Next() {
		var r HeatmapRow
		if err := rows.Scan(&r.Lat, &r.Lng, &r.AnchorScore, &r.CompPenalty); err != nil {
			return nil, fmt.Errorf("scan heatmap row: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}
