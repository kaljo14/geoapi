package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// BulkInsertSofiaplanFeatures inserts GeoJSON features into the given SofiaПлан table.
// propsJSON and geomJSON are parallel slices of raw JSON strings.
// The sql parameter must be a pre-built INSERT statement using unnest($1::text[], $2::text[]).
func (q *Queries) BulkInsertSofiaplanFeatures(ctx context.Context, sql string, propsJSON, geomJSON []string) error {
	if len(propsJSON) == 0 {
		return nil
	}
	_, err := q.db.Exec(ctx, sql, propsJSON, geomJSON)
	if err != nil {
		return fmt.Errorf("bulk insert sofiaplan features: %w", err)
	}
	return nil
}

// TruncateSofiaplanTable removes all rows from a SofiaПлан table before re-import.
func (q *Queries) TruncateSofiaplanTable(ctx context.Context, tableName string) error {
	// tableName is from our hardcoded list, not user input — safe to interpolate.
	_, err := q.db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	if err != nil {
		return fmt.Errorf("truncate %s: %w", tableName, err)
	}
	return nil
}

// BuildSofiaplanInsertSQL constructs the INSERT SQL for a given table.
// If bgs2005 is true, the geometry is reprojected from BGS2005 (EPSG:7801) to WGS84.
func BuildSofiaplanInsertSQL(tableName string, bgs2005 bool) string {
	geomExpr := `ST_SetSRID(ST_GeomFromGeoJSON(p.geom_json), 4326)`
	if bgs2005 {
		geomExpr = `ST_Transform(ST_SetSRID(ST_GeomFromGeoJSON(p.geom_json), 7801), 4326)`
	}
	return fmt.Sprintf(`
INSERT INTO %s (properties, geom)
SELECT p.props::jsonb, %s
FROM unnest($1::text[], $2::text[]) AS p(props, geom_json)
`, tableName, geomExpr)
}

// --- Context queries ---

const zoningContextSQL = `
SELECT properties
FROM sofiaplan_zoning
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const incomeContextSQL = `
SELECT properties
FROM sofiaplan_income
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const metroCatchmentContextSQL = `
SELECT properties
FROM sofiaplan_metro_catchments
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const pedestrianContextSQL = `
SELECT properties
FROM sofiaplan_pedestrian_syntax
WHERE ST_DWithin(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326), 0.005)
ORDER BY geom <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)
LIMIT 1
`

const businessTurnoverContextSQL = `
SELECT properties
FROM sofiaplan_business_turnover
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const propertyPriceContextSQL = `
SELECT properties
FROM sofiaplan_property_prices
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const populationContextSQL = `
SELECT properties
FROM sofiaplan_population_grid
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const developmentPotentialContextSQL = `
SELECT properties
FROM sofiaplan_development_potential
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

const neighborhoodContextSQL = `
SELECT properties
FROM sofiaplan_neighborhoods
WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
LIMIT 1
`

// scanProperties runs a spatial query and returns the JSONB properties, or nil if no row matches.
func (q *Queries) scanProperties(ctx context.Context, sql string, lng, lat float64) (json.RawMessage, error) {
	var raw []byte
	err := q.db.QueryRow(ctx, sql, lng, lat).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return json.RawMessage(raw), nil
}

func (q *Queries) GetZoningContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, zoningContextSQL, lng, lat)
}

func (q *Queries) GetIncomeContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, incomeContextSQL, lng, lat)
}

func (q *Queries) GetMetroCatchmentContext(ctx context.Context, lng, lat float64) (json.RawMessage, bool, error) {
	raw, err := q.scanProperties(ctx, metroCatchmentContextSQL, lng, lat)
	if err != nil {
		return nil, false, err
	}
	return raw, raw != nil, nil
}

func (q *Queries) GetPedestrianContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, pedestrianContextSQL, lng, lat)
}

func (q *Queries) GetBusinessTurnoverContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, businessTurnoverContextSQL, lng, lat)
}

func (q *Queries) GetPropertyPriceContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, propertyPriceContextSQL, lng, lat)
}

func (q *Queries) GetPopulationContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, populationContextSQL, lng, lat)
}

func (q *Queries) GetDevelopmentPotentialContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, developmentPotentialContextSQL, lng, lat)
}

func (q *Queries) GetNeighborhoodContext(ctx context.Context, lng, lat float64) (json.RawMessage, error) {
	return q.scanProperties(ctx, neighborhoodContextSQL, lng, lat)
}
