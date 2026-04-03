-- Sofia Plan flood risk datasets from https://api.sofiaplan.bg
-- Dataset 465: Flood risk — low probability
-- Dataset 412: Flood risk — medium probability
-- Dataset 446: Flood risk — high probability

CREATE TABLE IF NOT EXISTS sofiaplan_flood_risk_low (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_flood_risk_low_geom
    ON sofiaplan_flood_risk_low USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_flood_risk_medium (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_flood_risk_medium_geom
    ON sofiaplan_flood_risk_medium USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_flood_risk_high (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_flood_risk_high_geom
    ON sofiaplan_flood_risk_high USING GIST(geom);

-- Tile-serving views (Martin auto-discovers `_tiles` suffix)

-- Flood risk low probability (dataset 465)
CREATE OR REPLACE VIEW sofiaplan_flood_risk_low_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')  AS label,
    COALESCE(properties->>'rajon', '')    AS district,
    1                                     AS risk_level,
    geom
FROM sofiaplan_flood_risk_low;

-- Flood risk medium probability (dataset 412)
CREATE OR REPLACE VIEW sofiaplan_flood_risk_medium_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')  AS label,
    COALESCE(properties->>'rajon', '')    AS district,
    2                                     AS risk_level,
    geom
FROM sofiaplan_flood_risk_medium;

-- Flood risk high probability (dataset 446)
CREATE OR REPLACE VIEW sofiaplan_flood_risk_high_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')  AS label,
    COALESCE(properties->>'rajon', '')    AS district,
    3                                     AS risk_level,
    geom
FROM sofiaplan_flood_risk_high;
