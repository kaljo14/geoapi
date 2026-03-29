-- Population & Demographics datasets from https://api.sofiaplan.bg

CREATE TABLE IF NOT EXISTS sofiaplan_census_addresses (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_census_addresses_geom ON sofiaplan_census_addresses USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_demographic_forecast (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_demographic_forecast_geom ON sofiaplan_demographic_forecast USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_demographic_forecast_ge (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_demographic_forecast_ge_geom ON sofiaplan_demographic_forecast_ge USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_population_potential (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_population_potential_geom ON sofiaplan_population_potential USING GIST(geom);

CREATE TABLE IF NOT EXISTS sofiaplan_residential_load (
    id         SERIAL PRIMARY KEY,
    properties JSONB,
    geom       GEOMETRY(Geometry, 4326)
);
CREATE INDEX IF NOT EXISTS idx_sofiaplan_residential_load_geom ON sofiaplan_residential_load USING GIST(geom);

-- Tile-serving views

-- Census addresses: score = total people per address
CREATE OR REPLACE VIEW sofiaplan_census_addresses_tiles AS
SELECT
    id,
    (properties->>'nstreetnam')                          AS label,
    (properties->>'nnumber')                             AS street_number,
    (properties->>'ecode_rayon')                         AS rayon,
    COALESCE((properties->>'nbroi_lica')::numeric, 0)   AS score,
    COALESCE((properties->>'nn_jilisht')::numeric, 0)   AS dwellings,
    COALESCE((properties->>'nmale_sum')::numeric, 0)    AS male,
    COALESCE((properties->>'nfemale_su')::numeric, 0)   AS female,
    COALESCE((properties->>'nage0_14')::numeric, 0)     AS age_0_14,
    COALESCE((properties->>'nage15_24')::numeric, 0)    AS age_15_24,
    COALESCE((properties->>'nage25_34')::numeric, 0)    AS age_25_34,
    COALESCE((properties->>'nage35_44')::numeric, 0)    AS age_35_44,
    COALESCE((properties->>'nage45_54')::numeric, 0)    AS age_45_54,
    COALESCE((properties->>'nage55_64')::numeric, 0)    AS age_55_64,
    COALESCE((properties->>'nage65_')::numeric, 0)      AS age_65_plus,
    geom
FROM sofiaplan_census_addresses;

-- Demographic forecast (aggregated): score = current population (2017)
CREATE OR REPLACE VIEW sofiaplan_demographic_forecast_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon') AS label,
    (properties->>'rajon')                                  AS rajon,
    COALESCE((properties->>'nbroi_lica_sum')::numeric, 0)  AS score,
    COALESCE((properties->>'n2017_sum')::numeric, 0)       AS pop_2017,
    COALESCE((properties->>'pp2030_sum')::numeric, 0)      AS forecast_2030,
    COALESCE((properties->>'pp2040_sum')::numeric, 0)      AS forecast_2040,
    COALESCE((properties->>'pp2050_sum')::numeric, 0)      AS forecast_2050,
    COALESCE((properties->>'nmale_sum_sum')::numeric, 0)   AS male,
    COALESCE((properties->>'nfemale_su_sum')::numeric, 0)  AS female,
    COALESCE((properties->>'nage0_14_sum')::numeric, 0)    AS age_0_14,
    COALESCE((properties->>'nage65__sum')::numeric, 0)     AS age_65_plus,
    geom
FROM sofiaplan_demographic_forecast;

-- Demographic forecast per GE: score = population per planning unit
CREATE OR REPLACE VIEW sofiaplan_demographic_forecast_ge_tiles AS
SELECT
    id,
    COALESCE(properties->>'ime_predl', properties->>'regname') AS label,
    (properties->>'rajon')                                      AS rajon,
    (properties->>'ident_gr_u')                                 AS ge_id,
    COALESCE((properties->>'nbroi_lica')::numeric, 0)          AS score,
    COALESCE((properties->>'n2017')::numeric, 0)               AS pop_2017,
    COALESCE((properties->>'pp2030')::numeric, 0)              AS forecast_2030,
    COALESCE((properties->>'pp2040')::numeric, 0)              AS forecast_2040,
    COALESCE((properties->>'pp2050')::numeric, 0)              AS forecast_2050,
    COALESCE((properties->>'nmale_sum')::numeric, 0)           AS male,
    COALESCE((properties->>'nfemale_su')::numeric, 0)          AS female,
    geom
FROM sofiaplan_demographic_forecast_ge;

-- Population potential: score = potential population at 30m²/person
CREATE OR REPLACE VIEW sofiaplan_population_potential_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon') AS label,
    (properties->>'rajon')                                  AS rajon,
    (properties->>'posoka')                                 AS direction,
    COALESCE((properties->>'ppl_30kvm_sum')::numeric, 0)   AS score,
    COALESCE((properties->>'ppl_35kvm_sum')::numeric, 0)   AS ppl_35,
    COALESCE((properties->>'ppl_40kvm_sum')::numeric, 0)   AS ppl_40,
    COALESCE((properties->>'dens30_ppl_kvkm')::numeric, 0) AS density_30,
    COALESCE((properties->>'dens40_ppl_kvkm')::numeric, 0) AS density_40,
    geom
FROM sofiaplan_population_potential;

-- Residential building load: score = current population at 30m²/person
CREATE OR REPLACE VIEW sofiaplan_residential_load_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon') AS label,
    (properties->>'rajon')                                  AS rajon,
    COALESCE((properties->>'ppl_30kvm')::numeric, 0)       AS score,
    COALESCE((properties->>'ppl_35kvm')::numeric, 0)       AS ppl_35,
    COALESCE((properties->>'ppl_40kvm')::numeric, 0)       AS ppl_40,
    COALESCE((properties->>'rzp_jilsgr')::numeric, 0)      AS rzp_residential,
    COALESCE((properties->>'rzp_all')::numeric, 0)         AS rzp_total,
    COALESCE((properties->>'dens30_ppl_kvkm')::numeric, 0) AS density_30,
    COALESCE((properties->>'dens40_ppl_kvkm')::numeric, 0) AS density_40,
    geom
FROM sofiaplan_residential_load;
