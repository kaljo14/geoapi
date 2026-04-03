-- Tables for the grocery desert analysis pipeline.
-- Scores H3 hexagons against big-chain supermarkets from osm_pois.
-- Populated by geo-service/pipeline, served as tiles via the matview in 000035.

CREATE TABLE IF NOT EXISTS grocery_desert_h3 (
    h3_id       TEXT PRIMARY KEY,
    population  DOUBLE PRECISION NOT NULL DEFAULT 0,
    nearest_m   DOUBLE PRECISION,
    poi_count   INTEGER NOT NULL DEFAULT 0,
    score       DOUBLE PRECISION,
    geom        GEOMETRY(Polygon, 4326) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_gd_h3_geom
    ON grocery_desert_h3 USING GIST(geom);
CREATE INDEX IF NOT EXISTS idx_gd_h3_score
    ON grocery_desert_h3 (score);
