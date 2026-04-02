-- Fix parking_zones_tiles SRID=0: cast geom to geometry(Geometry,4326)
-- so PostGIS registers the column with SRID=4326 and Martin can serve it.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT
    ROW_NUMBER() OVER () AS id,
    color,
    geom::geometry(Geometry, 4326) AS geom
FROM (
    SELECT 'blue'::text AS color, geom
    FROM sofiaplan_parking_green

    UNION ALL

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
