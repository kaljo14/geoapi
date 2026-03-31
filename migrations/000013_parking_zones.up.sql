-- Paid parking zones for Sofia municipality (green zone, blue zone).
-- Polygons sourced from OpenStreetMap.

CREATE TABLE IF NOT EXISTS parking_zones (
    zone_id  INTEGER PRIMARY KEY,
    name     TEXT    NOT NULL,
    color    TEXT    NOT NULL,
    geom     GEOMETRY(MultiPolygon, 4326) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_parking_zones_geom ON parking_zones USING GIST(geom);

-- Martin auto-discovers this view and serves it as vector tiles.
CREATE OR REPLACE VIEW parking_zones_tiles AS
SELECT
    zone_id,
    name,
    color,
    geom
FROM parking_zones;
