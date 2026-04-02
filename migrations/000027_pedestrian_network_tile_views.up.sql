-- Tile-serving views for pedestrian network datasets.
-- Martin auto-discovers these views and serves them as vector tile sources.
-- Line network layers use space syntax properties (t1024_inte, t1024_choi, connectivi, segment_le).

-- Pedestrian network — Sofia city (primary): space syntax metrics
CREATE OR REPLACE VIEW sofiaplan_pedestrian_city_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_city;

-- Pedestrian network — Sofia city (alt): same metrics
CREATE OR REPLACE VIEW sofiaplan_pedestrian_city_alt_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_city_alt;

-- Pedestrian network — Sofia municipality: same metrics
CREATE OR REPLACE VIEW sofiaplan_pedestrian_municipality_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_municipality;

-- Pedestrian network — Sofia municipality (alt): same metrics
CREATE OR REPLACE VIEW sofiaplan_pedestrian_municipality_alt_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_municipality_alt;

-- Pedestrian network segmented: same metrics
CREATE OR REPLACE VIEW sofiaplan_pedestrian_segmented_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_segmented;

-- Pedestrian integration near infrastructure dividers (MultiPolygon, aggregated per area)
-- Uses same space syntax metrics aggregated per region
CREATE OR REPLACE VIEW sofiaplan_pedestrian_integration_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', '')                  AS label,
    COALESCE(properties->>'rajon', '')                    AS district,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    geom
FROM sofiaplan_pedestrian_integration;
