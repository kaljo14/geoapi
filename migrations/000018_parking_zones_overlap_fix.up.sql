-- Fix ST_Difference approach: previous migration filtered all green rows because
-- ST_Difference can return NULL or GEOMETRYCOLLECTION, making ST_IsEmpty(NULL)
-- evaluate to NULL (falsy), dropping every row.
-- Use a LATERAL subquery to compute the difference once and ST_CollectionExtract
-- to guarantee Martin receives a clean MULTIPOLYGON, not a GEOMETRYCOLLECTION.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
-- Blue zone (city centre) — unchanged
SELECT id, 'blue'::text AS color, geom
FROM sofiaplan_parking_green
UNION ALL
-- Green zone with blue areas punched out
SELECT
    g.id,
    'green'::text AS color,
    diff.geom
FROM sofiaplan_parking_blue g
CROSS JOIN (
    SELECT COALESCE(ST_Union(geom), 'GEOMETRYCOLLECTION EMPTY'::geometry) AS geom
    FROM sofiaplan_parking_green
) blue_union
CROSS JOIN LATERAL (
    SELECT ST_Multi(
        ST_CollectionExtract(ST_Difference(g.geom, blue_union.geom), 3)
    ) AS geom
) diff
WHERE diff.geom IS NOT NULL
  AND NOT ST_IsEmpty(diff.geom);
