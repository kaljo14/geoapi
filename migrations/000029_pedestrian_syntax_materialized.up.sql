-- Convert sofiaplan_pedestrian_syntax_tiles from a live view to a materialized view.
--
-- Previously (migration 000022) every tile request re-executed:
--   - JSONB extraction for all 305k rows
--   - ST_Contains spatial join against sofiaplan_neighborhoods (to assign neighbourhood)
--   - PERCENT_RANK() window function over the full dataset
-- All of that happened before Martin could even filter by tile bbox.
--
-- Now the result is pre-computed once and stored as a real table with a GIST index.
-- Martin queries it exactly like a regular table: bbox filter hits the index directly.
-- To update after a re-import: REFRESH MATERIALIZED VIEW sofiaplan_pedestrian_syntax_tiles;

DROP VIEW IF EXISTS sofiaplan_pedestrian_syntax_tiles;

CREATE MATERIALIZED VIEW sofiaplan_pedestrian_syntax_tiles AS
SELECT
    id,
    score,
    choice,
    connectivity,
    segment_length,
    neighborhood,
    PERCENT_RANK() OVER (
        PARTITION BY neighborhood_id
        ORDER BY score
    ) AS local_percentile,
    geom::geometry(Geometry, 4326) AS geom
FROM (
    SELECT
        p.id,
        COALESCE((p.properties->>'t1024_inte')::numeric, 0)  AS score,
        COALESCE((p.properties->>'t1024_choi')::numeric, 0)  AS choice,
        COALESCE((p.properties->>'connectivi')::numeric, 0)  AS connectivity,
        COALESCE((p.properties->>'segment_le')::numeric, 0)  AS segment_length,
        n.id                                                  AS neighborhood_id,
        (n.properties->>'kvname')                            AS neighborhood,
        p.geom
    FROM sofiaplan_pedestrian_syntax p
    LEFT JOIN sofiaplan_neighborhoods n
        ON ST_Contains(n.geom, ST_Centroid(p.geom))
) matched;

-- Spatial index so Martin (and PostGIS) can efficiently clip to tile bbox.
CREATE INDEX idx_pedestrian_syntax_tiles_geom
    ON sofiaplan_pedestrian_syntax_tiles USING GIST(geom);

-- Optional: speed up neighbourhood-filter queries from the frontend.
CREATE INDEX idx_pedestrian_syntax_tiles_neighborhood
    ON sofiaplan_pedestrian_syntax_tiles (neighborhood);
