-- Materialized view for Martin auto-discovery.
-- Follows the *_tiles naming convention used by all other tile layers.

CREATE MATERIALIZED VIEW IF NOT EXISTS grocery_desert_tiles AS
SELECT
    h3_id,
    population,
    nearest_m,
    poi_count,
    score,
    geom
FROM grocery_desert_h3
WHERE score IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_grocery_desert_tiles_h3_id
    ON grocery_desert_tiles (h3_id);

CREATE INDEX IF NOT EXISTS idx_grocery_desert_tiles_geom
    ON grocery_desert_tiles USING GIST(geom);
