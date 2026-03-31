DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'blue'::text  AS color, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'green'::text AS color, geom FROM sofiaplan_parking_blue;
