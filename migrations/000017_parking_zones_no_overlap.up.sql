-- Where blue and green zones overlap, only blue should be shown.
-- Fix at the geometry level: subtract the union of all blue-zone polygons
-- from each green-zone polygon before serving tiles.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
-- Blue zone (city centre) — unchanged
SELECT id, 'blue'::text AS color, geom
FROM sofiaplan_parking_green
UNION ALL
-- Green zone with blue-zone areas punched out
SELECT
    g.id,
    'green'::text AS color,
    ST_Difference(g.geom, blue.geom) AS geom
FROM sofiaplan_parking_blue g
CROSS JOIN (SELECT ST_Union(geom) AS geom FROM sofiaplan_parking_green) blue
WHERE NOT ST_IsEmpty(ST_Difference(g.geom, blue.geom));
