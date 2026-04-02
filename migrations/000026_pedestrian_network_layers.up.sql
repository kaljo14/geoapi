-- Pedestrian network datasets from https://api.sofiaplan.bg
-- Dataset IDs: 318, 309, 332, 361, 284, 603
-- Properties stored as JSONB. Line datasets use MultiLineString, 603 uses MultiPolygon.

-- Pedestrian network — Sofia city (primary) — ID 318
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_city (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_city_geom ON sofiaplan_pedestrian_city USING GIST(geom);

-- Pedestrian network — Sofia city (alt) — ID 309
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_city_alt (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_city_alt_geom ON sofiaplan_pedestrian_city_alt USING GIST(geom);

-- Pedestrian network — Sofia municipality — ID 332
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_municipality (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_municipality_geom ON sofiaplan_pedestrian_municipality USING GIST(geom);

-- Pedestrian network — Sofia municipality (alt) — ID 361
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_municipality_alt (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_municipality_alt_geom ON sofiaplan_pedestrian_municipality_alt USING GIST(geom);

-- Pedestrian network segmented — ID 284
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_segmented (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_segmented_geom ON sofiaplan_pedestrian_segmented USING GIST(geom);

-- Pedestrian integration near infrastructure dividers — ID 603
CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_integration (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_integration_geom ON sofiaplan_pedestrian_integration USING GIST(geom);
