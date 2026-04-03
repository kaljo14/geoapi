-- H3 hexagon grid with population data.
-- Populated by geo-service/pipeline, used as base grid for density layers.

CREATE TABLE IF NOT EXISTS food_desert_h3 (
    h3_id       TEXT PRIMARY KEY,
    population  DOUBLE PRECISION NOT NULL DEFAULT 0,
    nearest_m   DOUBLE PRECISION,
    poi_count   INTEGER NOT NULL DEFAULT 0,
    score       DOUBLE PRECISION,
    geom        GEOMETRY(Polygon, 4326) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_fd_h3_geom
    ON food_desert_h3 USING GIST(geom);
CREATE INDEX IF NOT EXISTS idx_fd_h3_score
    ON food_desert_h3 (score);
