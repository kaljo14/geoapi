-- Parking zones from SofiaПлан:
--   Dataset 470 — Green parking zone (Зелена зона за паркиране)
--   Dataset 291 — Blue parking zone  (Синя зона за паркиране)

CREATE TABLE IF NOT EXISTS sofiaplan_parking_green (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_parking_green_geom ON sofiaplan_parking_green USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_parking_blue (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_parking_blue_geom ON sofiaplan_parking_blue USING GIST(geom);

-- Replace parking_zones_tiles to union both SofiaПлан tables.
-- DROP first because CREATE OR REPLACE cannot rename existing columns.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'green' AS color, properties, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'blue'  AS color, properties, geom FROM sofiaplan_parking_blue;
