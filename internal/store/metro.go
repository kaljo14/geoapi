package store

import (
	"context"
	"fmt"
)

type MetroStationRow struct {
	StopID         string
	StopName       string
	StopLat        float64
	StopLon        float64
	Line           string
	RouteColor     string
	RouteTextColor string
}

const metroStationsSQL = `
SELECT DISTINCT
    s.stop_id,
    COALESCE(s.stop_name, '') AS stop_name,
    ST_Y(s.geom)::float8 AS stop_lat,
    ST_X(s.geom)::float8 AS stop_lon,
    r.route_id AS line,
    COALESCE(r.route_color, '') AS route_color,
    COALESCE(r.route_text_color, '') AS route_text_color
FROM gtfs_stops s
JOIN gtfs_stop_times st ON s.stop_id = st.stop_id
JOIN gtfs_trips t ON st.trip_id = t.trip_id
JOIN gtfs_routes r ON t.route_id = r.route_id
WHERE r.route_type = '1'
  AND s.geom IS NOT NULL
ORDER BY line, stop_name
`

func (q *Queries) GetMetroStations(ctx context.Context) ([]MetroStationRow, error) {
	rows, err := q.db.Query(ctx, metroStationsSQL)
	if err != nil {
		return nil, fmt.Errorf("get metro stations: %w", err)
	}
	defer rows.Close()

	var items []MetroStationRow
	for rows.Next() {
		var r MetroStationRow
		if err := rows.Scan(&r.StopID, &r.StopName, &r.StopLat, &r.StopLon, &r.Line, &r.RouteColor, &r.RouteTextColor); err != nil {
			return nil, fmt.Errorf("scan metro station: %w", err)
		}
		items = append(items, r)
	}
	return items, rows.Err()
}

type TransitStopRow struct {
	StopID   string
	StopName string
	StopLat  float64
	StopLon  float64
	StopType string // "bus", "tram", "trolleybus", "other"
}

const transitStopsSQL = `
SELECT
    stop_id,
    COALESCE(stop_name, '') AS stop_name,
    ST_Y(geom)::float8 AS stop_lat,
    ST_X(geom)::float8 AS stop_lon,
    CASE
        WHEN stop_id LIKE 'A%' THEN 'bus'
        WHEN stop_id LIKE 'TM%' THEN 'tram'
        WHEN stop_id LIKE 'TB%' THEN 'trolleybus'
        ELSE 'other'
    END AS stop_type
FROM gtfs_stops
WHERE geom IS NOT NULL
ORDER BY stop_name
`

func (q *Queries) GetTransitStops(ctx context.Context) ([]TransitStopRow, error) {
	rows, err := q.db.Query(ctx, transitStopsSQL)
	if err != nil {
		return nil, fmt.Errorf("get transit stops: %w", err)
	}
	defer rows.Close()

	var items []TransitStopRow
	for rows.Next() {
		var r TransitStopRow
		if err := rows.Scan(&r.StopID, &r.StopName, &r.StopLat, &r.StopLon, &r.StopType); err != nil {
			return nil, fmt.Errorf("scan transit stop: %w", err)
		}
		items = append(items, r)
	}
	return items, rows.Err()
}
