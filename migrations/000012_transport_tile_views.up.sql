-- Tile-serving views for Accessibility & Transport datasets.
-- Martin auto-discovers these views and serves them as vector tile sources.
-- Each view exposes a "score" column (primary numeric metric) and "label".
--
-- Property keys verified against actual imported data (2026-03-29).

-- Transit accessibility by GE: score = kgr index (0–1.2 float, higher = better access)
CREATE OR REPLACE VIEW sofiaplan_transit_access_ge_tiles AS
SELECT
    id,
    (properties->>'ime_predl')                                                        AS label,
    (properties->>'id')                                                               AS ge_id,
    (properties->>'rajon')                                                            AS rajon,
    COALESCE((properties->>'kgr')::numeric, 0)                                       AS score,
    geom
FROM sofiaplan_transit_access_ge;

-- Transit accessibility by transport district: score = total accessible route length (m, 7–2293)
CREATE OR REPLACE VIEW sofiaplan_transit_access_district_tiles AS
SELECT
    id,
    (properties->>'name')                                                             AS label,
    (properties->>'facilityid')                                                       AS district_id,
    COALESCE((properties->>'total_leng')::numeric, 0)                                AS score,
    geom
FROM sofiaplan_transit_access_district;

-- Metro accessibility 800 m catchment: score = catchment break distance (m)
CREATE OR REPLACE VIEW sofiaplan_metro_access_800m_tiles AS
SELECT
    id,
    (properties->>'name')                                                             AS label,
    COALESCE((properties->>'tobreak')::numeric, 800)                                 AS score,
    COALESCE((properties->>'frombreak')::numeric, 0)                                 AS frombreak,
    (properties->>'facilityid')                                                      AS facilityid,
    geom
FROM sofiaplan_metro_access_800m;

-- Metro accessibility 1200 m+ catchment: score = catchment break distance (m)
CREATE OR REPLACE VIEW sofiaplan_metro_access_1200m_tiles AS
SELECT
    id,
    (properties->>'name')                                                             AS label,
    COALESCE((properties->>'tobreak')::numeric, 1200)                                AS score,
    COALESCE((properties->>'frombreak')::numeric, 0)                                 AS frombreak,
    (properties->>'facilityid')                                                      AS facilityid,
    geom
FROM sofiaplan_metro_access_1200m;

-- Bus lines (primary): label = route number
CREATE OR REPLACE VIEW sofiaplan_bus_lines_tiles AS
SELECT
    id,
    (properties->>'line_bus')                                                        AS label,
    (properties->>'line_bus')                                                        AS route_id,
    geom
FROM sofiaplan_bus_lines;

-- Bus lines (alternate dataset): label = line name (Cyrillic key "ЛИНИЯ")
CREATE OR REPLACE VIEW sofiaplan_bus_lines_alt_tiles AS
SELECT
    id,
    (properties->>'ЛИНИЯ')                                                           AS label,
    (properties->>'ЛИНИЯ')                                                           AS route_id,
    geom
FROM sofiaplan_bus_lines_alt;

-- Trolleybus lines: label = trolleybus line number
CREATE OR REPLACE VIEW sofiaplan_trolleybus_lines_tiles AS
SELECT
    id,
    (properties->>'line_tb')                                                         AS label,
    (properties->>'line_tb')                                                         AS route_id,
    geom
FROM sofiaplan_trolleybus_lines;

-- Tram lines (primary): label = tram line number
CREATE OR REPLACE VIEW sofiaplan_tram_lines_tiles AS
SELECT
    id,
    (properties->>'line_tram')                                                       AS label,
    (properties->>'line_tram')                                                       AS route_id,
    geom
FROM sofiaplan_tram_lines;

-- Tram lines (alternate dataset): label = tram line number
CREATE OR REPLACE VIEW sofiaplan_tram_lines_alt_tiles AS
SELECT
    id,
    (properties->>'line_tram')                                                       AS label,
    (properties->>'line_tram')                                                       AS route_id,
    geom
FROM sofiaplan_tram_lines_alt;

-- Railway stations: score = total annual passengers 2019 (arrivals + departures)
CREATE OR REPLACE VIEW sofiaplan_railway_stations_tiles AS
SELECT
    id,
    (properties->>'tradename')                                                       AS label,
    COALESCE(
        (properties->>'2019_prist')::numeric, 0
    ) + COALESCE(
        (properties->>'2019_zamin')::numeric, 0
    )                                                                                AS score,
    geom
FROM sofiaplan_railway_stations;
