-- Cycling network datasets from https://api.sofiaplan.bg
-- Dataset IDs: 606, 290, 146
-- Properties stored as JSONB. Geometry is MultiLineString (4326).

-- Built cycling network (primary) — ID 606
-- Properties: id, type (path type), posoka (direction), length (metres)
CREATE TABLE IF NOT EXISTS sofiaplan_cycling_network (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_cycling_network_geom ON sofiaplan_cycling_network USING GIST(geom);

-- Built cycling network (alternate dataset) — ID 290
-- Properties: id, type, posoka (no length field)
CREATE TABLE IF NOT EXISTS sofiaplan_cycling_network_alt (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_cycling_network_alt_geom ON sofiaplan_cycling_network_alt USING GIST(geom);

-- Planned cycling extensions — ID 146
-- Properties: id, name (street name), priority, project, note
CREATE TABLE IF NOT EXISTS sofiaplan_cycling_planned (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_cycling_planned_geom ON sofiaplan_cycling_planned USING GIST(geom);
