-- Sofia Plan GE building analysis datasets from https://api.sofiaplan.bg

CREATE TABLE IF NOT EXISTS sofiaplan_building_density_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_building_density_ge_geom
    ON sofiaplan_building_density_ge USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_building_footprint_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_building_footprint_ge_geom
    ON sofiaplan_building_footprint_ge USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_residential_typology_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_residential_typology_ge_geom
    ON sofiaplan_residential_typology_ge USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_urban_morphology_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_urban_morphology_ge_geom
    ON sofiaplan_urban_morphology_ge USING GIST(geom);

-- Tile-serving views (Martin auto-discovers `_tiles` suffix)
-- Field names are in Bulgarian as returned by api.sofiaplan.bg

-- Building density GE (dataset 632): zastr_plytnost = plot coverage ratio (0–1)
CREATE OR REPLACE VIEW sofiaplan_building_density_ge_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')                        AS label,
    COALESCE(properties->>'rajon', '')                          AS district,
    COALESCE((properties->>'zastr_plytnost')::numeric, 0)       AS score,
    COALESCE((properties->>'zastr_intenzivnost')::numeric, 0)   AS intensity,
    COALESCE((properties->>'zastr_sklyuchenost')::numeric, 0)   AS enclosure_ratio,
    COALESCE((properties->>'sredna_etajnost')::numeric, 0)      AS avg_floors,
    geom
FROM sofiaplan_building_density_ge;

-- Building footprint GE (dataset 633): rzp = total floor area (m²)
CREATE OR REPLACE VIEW sofiaplan_building_footprint_ge_tiles AS
SELECT
    id,
    COALESCE(properties->>'ge_id', '')                          AS ge_id,
    COALESCE(properties->>'funktyp_gen_txt', '')                AS label,
    ''                                                          AS district,
    COALESCE((properties->>'rzp')::numeric, 0)                  AS score,
    COALESCE((properties->>'zp')::numeric, 0)                   AS zp,
    COALESCE((properties->>'rzp')::numeric, 0)                  AS rzp,
    COALESCE((properties->>'ge_sgradi_broi')::numeric, 0)       AS avg_floors,
    geom
FROM sofiaplan_building_footprint_ge;

-- Residential typology GE (dataset 626): dial_ednfa/dial_mnfam/dial_panel = % share per type
CREATE OR REPLACE VIEW sofiaplan_residential_typology_ge_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')                        AS label,
    COALESCE(properties->>'rajon', '')                          AS district,
    COALESCE((properties->>'dial_ednfa')::numeric, 0)           AS single_pct,
    COALESCE((properties->>'dial_mnfam')::numeric, 0)           AS multi_pct,
    COALESCE((properties->>'dial_panel')::numeric, 0)           AS panel_pct,
    -- score = dominant type: 1=single-family, 2=multi-family, 3=panel
    CASE
        WHEN COALESCE((properties->>'dial_panel')::numeric, 0) > COALESCE((properties->>'dial_ednfa')::numeric, 0)
         AND COALESCE((properties->>'dial_panel')::numeric, 0) > COALESCE((properties->>'dial_mnfam')::numeric, 0) THEN 3
        WHEN COALESCE((properties->>'dial_mnfam')::numeric, 0) > COALESCE((properties->>'dial_ednfa')::numeric, 0) THEN 2
        WHEN COALESCE((properties->>'dial_ednfa')::numeric, 0) > 0 THEN 1
        ELSE 0
    END                                                         AS score,
    geom
FROM sofiaplan_residential_typology_ge;

-- Urban morphology GE (dataset 455): derived from building height distribution + density
-- bui1_perc/bui2_perc = % 1–2 floor buildings; bui4–bui78 = % mid/high-rise; gl_dens = density
CREATE OR REPLACE VIEW sofiaplan_urban_morphology_ge_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')                        AS label,
    COALESCE(properties->>'rajon', '')                          AS district,
    -- score: 1=вилна/ниска, 2=компактна, 3=панелни/високи, 5=смесена
    CASE
        WHEN COALESCE((properties->>'bui1_perc')::numeric, 0) + COALESCE((properties->>'bui2_perc')::numeric, 0) >= 65
          AND COALESCE((properties->>'gl_dens')::numeric, 0) < 8                                               THEN 1
        WHEN COALESCE((properties->>'bui4_perc')::numeric, 0) + COALESCE((properties->>'bui5_perc')::numeric, 0)
           + COALESCE((properties->>'bui6_perc')::numeric, 0) + COALESCE((properties->>'bui78_arep')::numeric, 0) >= 10 THEN 3
        WHEN COALESCE((properties->>'gl_dens')::numeric, 0) >= 15                                              THEN 2
        ELSE 5
    END                                                         AS score,
    COALESCE((properties->>'gl_dens')::numeric, 0)             AS density,
    geom
FROM sofiaplan_urban_morphology_ge;
