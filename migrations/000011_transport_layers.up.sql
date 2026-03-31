-- Accessibility & Transport datasets from https://api.sofiaplan.bg
-- Dataset IDs: 96, 279, 289, 282, 268, 333, 223, 472, 254, 602
-- Properties stored as JSONB (schema confirmed after first import).
-- GEOMETRY(Geometry, 4326) accepts any geometry type without casting errors.
-- Martin auto-discovers these tables on startup.

-- Transit accessibility by GE (Градска Единица / urban planning unit) — ID 96
CREATE TABLE IF NOT EXISTS sofiaplan_transit_access_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_transit_access_ge_geom ON sofiaplan_transit_access_ge USING GIST(geom);

-- Transit accessibility by transport district — ID 279
CREATE TABLE IF NOT EXISTS sofiaplan_transit_access_district (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_transit_access_district_geom ON sofiaplan_transit_access_district USING GIST(geom);

-- Metro accessibility catchment zone 800 m — ID 289
CREATE TABLE IF NOT EXISTS sofiaplan_metro_access_800m (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_metro_access_800m_geom ON sofiaplan_metro_access_800m USING GIST(geom);

-- Metro accessibility catchment zone 1200 m+ — ID 282
CREATE TABLE IF NOT EXISTS sofiaplan_metro_access_1200m (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_metro_access_1200m_geom ON sofiaplan_metro_access_1200m USING GIST(geom);

-- Bus lines (primary) — ID 268
CREATE TABLE IF NOT EXISTS sofiaplan_bus_lines (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_bus_lines_geom ON sofiaplan_bus_lines USING GIST(geom);

-- Bus lines (alternate/secondary dataset) — ID 333
CREATE TABLE IF NOT EXISTS sofiaplan_bus_lines_alt (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_bus_lines_alt_geom ON sofiaplan_bus_lines_alt USING GIST(geom);

-- Trolleybus lines — ID 223
CREATE TABLE IF NOT EXISTS sofiaplan_trolleybus_lines (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_trolleybus_lines_geom ON sofiaplan_trolleybus_lines USING GIST(geom);

-- Tram lines (primary) — ID 472
CREATE TABLE IF NOT EXISTS sofiaplan_tram_lines (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_tram_lines_geom ON sofiaplan_tram_lines USING GIST(geom);

-- Tram lines (alternate/secondary dataset) — ID 254
CREATE TABLE IF NOT EXISTS sofiaplan_tram_lines_alt (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_tram_lines_alt_geom ON sofiaplan_tram_lines_alt USING GIST(geom);

-- Railway stations with passenger load — ID 602
CREATE TABLE IF NOT EXISTS sofiaplan_railway_stations (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_railway_stations_geom ON sofiaplan_railway_stations USING GIST(geom);
