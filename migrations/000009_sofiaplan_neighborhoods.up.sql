CREATE TABLE IF NOT EXISTS sofiaplan_neighborhoods (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_neighborhoods_geom ON sofiaplan_neighborhoods USING GIST(geom);
