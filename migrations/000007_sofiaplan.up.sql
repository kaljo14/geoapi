-- SofiaПлан datasets from https://api.sofiaplan.bg
-- Properties stored as JSONB (schema unknown until download).
-- GEOMETRY(Geometry, 4326) accepts any geometry type (Polygon, LineString, Point)
-- without casting errors. Martin auto-discovers these tables on startup.

CREATE TABLE IF NOT EXISTS sofiaplan_zoning (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_zoning_geom ON sofiaplan_zoning USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_zoning_params (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_zoning_params_geom ON sofiaplan_zoning_params USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_income (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_income_geom ON sofiaplan_income USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_business_turnover (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_business_turnover_geom ON sofiaplan_business_turnover USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_property_prices (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_property_prices_geom ON sofiaplan_property_prices USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_pedestrian_syntax (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_pedestrian_syntax_geom ON sofiaplan_pedestrian_syntax USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_metro_catchments (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_metro_catchments_geom ON sofiaplan_metro_catchments USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_population_grid (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_population_grid_geom ON sofiaplan_population_grid USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_development_potential (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_development_potential_geom ON sofiaplan_development_potential USING GIST(geom);
