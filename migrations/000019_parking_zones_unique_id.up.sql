-- Two problems in the previous view:
-- 1. Duplicate IDs: both sofiaplan_parking_green and sofiaplan_parking_blue
--    have id SERIAL starting at 1, so the UNION ALL produces duplicate id values
--    which Martin uses as vector tile feature IDs — this breaks tile generation.
-- 2. ST_Multi(ST_CollectionExtract(...)) strips the SRID, causing Martin to fail
--    on SRID detection for the geometry column.
--
-- Fix: generate a globally unique id with ROW_NUMBER() and preserve SRID with
-- ST_SetSRID on the difference geometry.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT
    ROW_NUMBER() OVER () AS id,
    color,
    geom
FROM (
    -- Blue zone (city centre) — unchanged
    SELECT 'blue'::text AS color, geom
    FROM sofiaplan_parking_green

    UNION ALL

    -- Green zone with blue areas punched out; preserve SRID 4326
    SELECT
        'green'::text AS color,
        ST_SetSRID(
            ST_Multi(
                ST_CollectionExtract(ST_Difference(g.geom, blue_union.geom), 3)
            ),
            4326
        ) AS geom
    FROM sofiaplan_parking_blue g
    CROSS JOIN (
        SELECT COALESCE(ST_Union(geom), ST_GeomFromText('GEOMETRYCOLLECTION EMPTY', 4326)) AS geom
        FROM sofiaplan_parking_green
    ) blue_union
    WHERE ST_Difference(g.geom, blue_union.geom) IS NOT NULL
      AND NOT ST_IsEmpty(ST_Difference(g.geom, blue_union.geom))
) zones;
