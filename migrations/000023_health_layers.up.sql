-- Health service datasets from https://api.sofiaplan.bg
-- Dataset IDs: 597, 598
-- Properties stored as JSONB. Geometry is Polygon/MultiPolygon (4326).

-- Health service concentration — ID 597
CREATE TABLE IF NOT EXISTS sofiaplan_health_service_concentration (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_health_service_concentration_geom
    ON sofiaplan_health_service_concentration USING GIST(geom);

-- Health infrastructure concentration by GE — ID 598
CREATE TABLE IF NOT EXISTS sofiaplan_health_infrastructure_concentration (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_health_infrastructure_concentration_geom
    ON sofiaplan_health_infrastructure_concentration USING GIST(geom);
