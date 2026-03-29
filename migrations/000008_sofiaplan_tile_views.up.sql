-- Tile-serving views that flatten the JSONB properties column into real columns.
-- Martin auto-discovers these views and serves them as vector tile sources.
-- Each view exposes a "score" column (the primary numeric metric) plus
-- a "label" column and other useful properties.

-- Income: score = average monthly income (BGN)
CREATE OR REPLACE VIEW sofiaplan_income_tiles AS
SELECT
    id,
    (properties->>'obns_cyr')                   AS label,
    (properties->>'obns_num')                    AS obns_num,
    COALESCE((properties->>'sr_mes_dohod')::numeric, 0)  AS score,
    COALESCE((properties->>'sr_god_dohod')::numeric, 0)  AS sr_god_dohod,
    COALESCE((properties->>'percent')::numeric, 0)       AS percent,
    COALESCE((properties->>'dohod_ofis')::numeric, 0)    AS dohod_ofis,
    geom
FROM sofiaplan_income;

-- Zoning: categorical — score encodes the zone type string hash for colouring
CREATE OR REPLACE VIEW sofiaplan_zoning_tiles AS
SELECT
    id,
    (properties->>'type_')                       AS label,
    (properties->>'center')                      AS center,
    (properties->>'new_end')                     AS new_end,
    -- Use a stable hash so the frontend can map zone types to colours
    abs(hashtext(COALESCE(properties->>'type_', ''))) % 100 AS score,
    geom
FROM sofiaplan_zoning;

-- Business turnover: score = turnover 2016
CREATE OR REPLACE VIEW sofiaplan_business_turnover_tiles AS
SELECT
    id,
    (properties->>'obns_cyr')                           AS label,
    (properties->>'obns_num')                            AS obns_num,
    COALESCE((properties->>'oborot_2016')::numeric, 0)   AS score,
    COALESCE((properties->>'oborot_2010')::numeric, 0)   AS oborot_2010,
    COALESCE((properties->>'prirast')::numeric, 0)       AS prirast,
    COALESCE((properties->>'prirast_perc')::numeric, 0)  AS prirast_perc,
    COALESCE((properties->>'dial_2016')::numeric, 0)     AS dial_2016,
    COALESCE((properties->>'norma_raztej')::numeric, 0)  AS norma_raztej,
    geom
FROM sofiaplan_business_turnover;

-- Property prices: score = apartment price per sq.m.
CREATE OR REPLACE VIEW sofiaplan_property_prices_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'kvartal')  AS label,
    (properties->>'kvartal')                                   AS kvartal,
    (properties->>'regname')                                   AS regname,
    (properties->>'godina')                                    AS godina,
    COALESCE((properties->>'cena_ap_kv_m')::numeric, 0)       AS score,
    COALESCE((properties->>'cena_ap')::numeric, 0)             AS cena_ap,
    COALESCE((properties->>'cena_ofis')::numeric, 0)           AS cena_ofis,
    COALESCE((properties->>'cena_ofis_kv_m')::numeric, 0)     AS cena_ofis_kv_m,
    COALESCE((properties->>'naem_ap')::numeric, 0)             AS naem_ap,
    COALESCE((properties->>'naem_ap_kv_m')::numeric, 0)       AS naem_ap_kv_m,
    COALESCE((properties->>'naem_ofis')::numeric, 0)           AS naem_ofis,
    COALESCE((properties->>'naem_ofis_kv_m')::numeric, 0)     AS naem_ofis_kv_m,
    COALESCE((properties->>'naem_mag')::numeric, 0)            AS naem_mag,
    COALESCE((properties->>'naem_mag_kv_m')::numeric, 0)      AS naem_mag_kv_m,
    geom
FROM sofiaplan_property_prices;

-- Pedestrian space syntax: score = integration (main walkability metric)
CREATE OR REPLACE VIEW sofiaplan_pedestrian_syntax_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_syntax;

-- Metro catchments: score = catchment radius (tobreak)
CREATE OR REPLACE VIEW sofiaplan_metro_catchments_tiles AS
SELECT
    id,
    (properties->>'name')                                AS label,
    COALESCE((properties->>'tobreak')::numeric, 0)       AS score,
    COALESCE((properties->>'frombreak')::numeric, 0)     AS frombreak,
    (properties->>'facilityid')                          AS facilityid,
    geom
FROM sofiaplan_metro_catchments;

-- Population grid: score = total population
CREATE OR REPLACE VIEW sofiaplan_population_grid_tiles AS
SELECT
    id,
    (properties->>'grd_id')                              AS label,
    COALESCE((properties->>'tot_p')::numeric, 0)         AS score,
    COALESCE((properties->>'tot_m')::numeric, 0)         AS tot_m,
    COALESCE((properties->>'tot_f')::numeric, 0)         AS tot_f,
    COALESCE((properties->>'t_00_14')::numeric, 0)       AS t_00_14,
    COALESCE((properties->>'t_15_64')::numeric, 0)       AS t_15_64,
    COALESCE((properties->>'t_65_')::numeric, 0)         AS t_65_plus,
    geom
FROM sofiaplan_population_grid;

-- Neighborhoods: score = type_kv (neighborhood type code)
CREATE OR REPLACE VIEW sofiaplan_neighborhoods_tiles AS
SELECT
    id,
    (properties->>'kvname')                              AS label,
    (properties->>'prefname')                            AS prefname,
    COALESCE((properties->>'adm_kvid')::numeric, 0)     AS adm_kvid,
    COALESCE((properties->>'type_kv')::numeric, 0)       AS score,
    geom
FROM sofiaplan_neighborhoods;

-- Development potential: score = zastr_potencial (1-3 scale)
CREATE OR REPLACE VIEW sofiaplan_development_potential_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon') AS label,
    (properties->>'rajon')                                  AS rajon,
    (properties->>'zastr_potenc_txt')                       AS description,
    COALESCE((properties->>'zastr_potencial')::numeric, 0)  AS score,
    geom
FROM sofiaplan_development_potential;
