CREATE MATERIALIZED VIEW IF NOT EXISTS opportunity_heatmap AS
SELECT
    p.location                                                   AS geom,
    p.category,
    LEAST(100.0, GREATEST(0.0,
        0.6 * LEAST(100.0, COALESCE((
            SELECT SUM(1.0 - ST_Distance(s.geom::geography, p.location::geography) / 1000.0)
            FROM gtfs_stops s
            WHERE s.geom IS NOT NULL
              AND ST_DWithin(s.geom::geography, p.location::geography, 1000)
        ), 0) * 20) +
        0.4 * GREATEST(0.0, 100.0 - COALESCE((
            SELECT COUNT(*) * 10.0
            FROM places p2
            WHERE p2.location IS NOT NULL
              AND p2.category = p.category
              AND p2.business_status = 'OPERATIONAL'
              AND ST_DWithin(p2.location::geography, p.location::geography, 500)
        ), 0))
    ))                                                           AS score
FROM places p
WHERE p.location IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_opportunity_heatmap_geom ON opportunity_heatmap USING GIST(geom);
