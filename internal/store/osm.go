package store

import (
	"context"
	"fmt"
)

// OSMNode is a parsed OpenStreetMap node ready for bulk insert.
type OSMNode struct {
	OsmID int64
	Lat   float64
	Lng   float64
}

// OSMEdge is a directed segment between two OSM nodes along a way.
type OSMEdge struct {
	OsmWayID   int64
	FromNodeID int64
	ToNodeID   int64
	Highway    string
	FromLat    float64
	FromLng    float64
	ToLat      float64
	ToLng      float64
}

// NodeScore holds a computed walkability score for a node.
type NodeScore struct {
	OsmID     int64
	WalkScore float64
}

// OSMPOI is a parsed OpenStreetMap point of interest ready for bulk insert.
type OSMPOI struct {
	OsmID    int64
	Lat      float64
	Lng      float64
	Name     string
	Amenity  string
	Shop     string
	Tourism  string
	Leisure  string
	Category string
}

const bulkUpsertNodesSQL = `
INSERT INTO osm_nodes (osm_id, lat, lng, geom)
SELECT
    n.osm_id,
    n.lat,
    n.lng,
    ST_SetSRID(ST_MakePoint(n.lng, n.lat), 4326)
FROM unnest($1::bigint[], $2::float8[], $3::float8[]) AS n(osm_id, lat, lng)
ON CONFLICT (osm_id) DO NOTHING
`

func (q *Queries) BulkUpsertNodes(ctx context.Context, nodes []OSMNode) error {
	if len(nodes) == 0 {
		return nil
	}
	ids := make([]int64, len(nodes))
	lats := make([]float64, len(nodes))
	lngs := make([]float64, len(nodes))
	for i, n := range nodes {
		ids[i] = n.OsmID
		lats[i] = n.Lat
		lngs[i] = n.Lng
	}
	_, err := q.db.Exec(ctx, bulkUpsertNodesSQL, ids, lats, lngs)
	if err != nil {
		return fmt.Errorf("bulk upsert nodes: %w", err)
	}
	return nil
}

const bulkUpsertEdgesSQL = `
INSERT INTO osm_edges (osm_way_id, from_node_id, to_node_id, highway, geom)
SELECT
    e.way_id,
    e.from_id,
    e.to_id,
    e.highway,
    ST_MakeLine(
        ST_SetSRID(ST_MakePoint(e.from_lng, e.from_lat), 4326),
        ST_SetSRID(ST_MakePoint(e.to_lng,   e.to_lat),   4326)
    )
FROM unnest(
    $1::bigint[],
    $2::bigint[],
    $3::bigint[],
    $4::text[],
    $5::float8[],
    $6::float8[],
    $7::float8[],
    $8::float8[]
) AS e(way_id, from_id, to_id, highway, from_lat, from_lng, to_lat, to_lng)
ON CONFLICT (osm_way_id, from_node_id, to_node_id) DO NOTHING
`

func (q *Queries) BulkUpsertEdges(ctx context.Context, edges []OSMEdge) error {
	if len(edges) == 0 {
		return nil
	}
	wayIDs := make([]int64, len(edges))
	fromIDs := make([]int64, len(edges))
	toIDs := make([]int64, len(edges))
	highways := make([]string, len(edges))
	fromLats := make([]float64, len(edges))
	fromLngs := make([]float64, len(edges))
	toLats := make([]float64, len(edges))
	toLngs := make([]float64, len(edges))
	for i, e := range edges {
		wayIDs[i] = e.OsmWayID
		fromIDs[i] = e.FromNodeID
		toIDs[i] = e.ToNodeID
		highways[i] = e.Highway
		fromLats[i] = e.FromLat
		fromLngs[i] = e.FromLng
		toLats[i] = e.ToLat
		toLngs[i] = e.ToLng
	}
	_, err := q.db.Exec(ctx, bulkUpsertEdgesSQL,
		wayIDs, fromIDs, toIDs, highways,
		fromLats, fromLngs, toLats, toLngs,
	)
	if err != nil {
		return fmt.Errorf("bulk upsert edges: %w", err)
	}
	return nil
}

// computeWalkScoresSQL scores each node by counting OPERATIONAL places AND
// osm_pois within 800 m. Never SELECTs geom — rebuilds points inline from
// lat/lng columns.
const computeWalkScoresSQL = `
SELECT
    n.osm_id,
    LEAST(100.0, COALESCE(
        (COUNT(DISTINCT p.place_id) + COUNT(DISTINCT op.osm_id))::float8 * 5.0, 0
    )) AS walk_score
FROM osm_nodes n
LEFT JOIN places p
    ON p.location IS NOT NULL
    AND p.business_status = 'OPERATIONAL'
    AND ST_DWithin(p.location::geography,
        ST_SetSRID(ST_MakePoint(n.lng, n.lat), 4326)::geography, 800)
LEFT JOIN osm_pois op
    ON ST_DWithin(op.geom::geography,
        ST_SetSRID(ST_MakePoint(n.lng, n.lat), 4326)::geography, 800)
GROUP BY n.osm_id
`

func (q *Queries) ComputeWalkScores(ctx context.Context) ([]NodeScore, error) {
	rows, err := q.db.Query(ctx, computeWalkScoresSQL)
	if err != nil {
		return nil, fmt.Errorf("compute walk scores: %w", err)
	}
	defer rows.Close()
	var result []NodeScore
	for rows.Next() {
		var s NodeScore
		if err := rows.Scan(&s.OsmID, &s.WalkScore); err != nil {
			return nil, fmt.Errorf("scan walk score: %w", err)
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

const updateNodeWalkScoresSQL = `
UPDATE osm_nodes n
SET walk_score = s.score
FROM unnest($1::bigint[], $2::float8[]) AS s(osm_id, score)
WHERE n.osm_id = s.osm_id
`

func (q *Queries) UpdateNodeWalkScores(ctx context.Context, scores []NodeScore) error {
	if len(scores) == 0 {
		return nil
	}
	ids := make([]int64, len(scores))
	vals := make([]float64, len(scores))
	for i, s := range scores {
		ids[i] = s.OsmID
		vals[i] = s.WalkScore
	}
	_, err := q.db.Exec(ctx, updateNodeWalkScoresSQL, ids, vals)
	if err != nil {
		return fmt.Errorf("update node walk scores: %w", err)
	}
	return nil
}

const propagateEdgeScoresSQL = `
UPDATE osm_edges e
SET walk_score = (
    SELECT COALESCE(AVG(n.walk_score), 0)
    FROM osm_nodes n
    WHERE n.osm_id IN (e.from_node_id, e.to_node_id)
)
`

func (q *Queries) PropagateEdgeScores(ctx context.Context) error {
	_, err := q.db.Exec(ctx, propagateEdgeScoresSQL)
	if err != nil {
		return fmt.Errorf("propagate edge scores: %w", err)
	}
	return nil
}

const bulkUpsertPOIsSQL = `
INSERT INTO osm_pois (osm_id, lat, lng, geom, name, amenity, shop, tourism, leisure, category)
SELECT
    p.osm_id, p.lat, p.lng,
    ST_SetSRID(ST_MakePoint(p.lng, p.lat), 4326),
    p.name, p.amenity, p.shop, p.tourism, p.leisure, p.category
FROM unnest($1::bigint[], $2::float8[], $3::float8[], $4::text[], $5::text[], $6::text[], $7::text[], $8::text[], $9::text[])
     AS p(osm_id, lat, lng, name, amenity, shop, tourism, leisure, category)
ON CONFLICT (osm_id) DO NOTHING
`

func (q *Queries) BulkUpsertPOIs(ctx context.Context, pois []OSMPOI) error {
	if len(pois) == 0 {
		return nil
	}
	ids := make([]int64, len(pois))
	lats := make([]float64, len(pois))
	lngs := make([]float64, len(pois))
	names := make([]string, len(pois))
	amenities := make([]string, len(pois))
	shops := make([]string, len(pois))
	tourisms := make([]string, len(pois))
	leisures := make([]string, len(pois))
	categories := make([]string, len(pois))
	for i, p := range pois {
		ids[i] = p.OsmID
		lats[i] = p.Lat
		lngs[i] = p.Lng
		names[i] = p.Name
		amenities[i] = p.Amenity
		shops[i] = p.Shop
		tourisms[i] = p.Tourism
		leisures[i] = p.Leisure
		categories[i] = p.Category
	}
	_, err := q.db.Exec(ctx, bulkUpsertPOIsSQL,
		ids, lats, lngs, names, amenities, shops, tourisms, leisures, categories,
	)
	if err != nil {
		return fmt.Errorf("bulk upsert pois: %w", err)
	}
	return nil
}
