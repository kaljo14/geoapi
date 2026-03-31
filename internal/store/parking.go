package store

import (
	"context"
	"encoding/json"
	"fmt"
)

// ParkingZone holds the non-geometry columns plus the GeoJSON geometry string.
type ParkingZone struct {
	ZoneID  int    `json:"zone_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	GeoJSON string `json:"geojson"` // ST_AsGeoJSON output
}

const listParkingZonesSQL = `
SELECT zone_id, name, color, ST_AsGeoJSON(geom) AS geojson
FROM parking_zones
ORDER BY zone_id
`

// ListParkingZones returns all parking zones with their geometry as GeoJSON strings.
func (q *Queries) ListParkingZones(ctx context.Context) ([]ParkingZone, error) {
	rows, err := q.db.Query(ctx, listParkingZonesSQL)
	if err != nil {
		return nil, fmt.Errorf("list parking zones: %w", err)
	}
	defer rows.Close()

	var zones []ParkingZone
	for rows.Next() {
		var z ParkingZone
		if err := rows.Scan(&z.ZoneID, &z.Name, &z.Color, &z.GeoJSON); err != nil {
			return nil, fmt.Errorf("scan parking zone: %w", err)
		}
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

// ParkingZonesGeoJSON builds a GeoJSON FeatureCollection from all parking zones.
func (q *Queries) ParkingZonesGeoJSON(ctx context.Context) (json.RawMessage, error) {
	zones, err := q.ListParkingZones(ctx)
	if err != nil {
		return nil, err
	}

	type Feature struct {
		Type       string          `json:"type"`
		Properties map[string]any  `json:"properties"`
		Geometry   json.RawMessage `json:"geometry"`
	}
	type FeatureCollection struct {
		Type     string    `json:"type"`
		Features []Feature `json:"features"`
	}

	fc := FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]Feature, 0, len(zones)),
	}
	for _, z := range zones {
		fc.Features = append(fc.Features, Feature{
			Type: "Feature",
			Properties: map[string]any{
				"zone_id": z.ZoneID,
				"name":    z.Name,
				"color":   z.Color,
			},
			Geometry: json.RawMessage(z.GeoJSON),
		})
	}

	raw, err := json.Marshal(fc)
	if err != nil {
		return nil, fmt.Errorf("marshal parking zones geojson: %w", err)
	}
	return raw, nil
}
