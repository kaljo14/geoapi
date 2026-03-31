-- sofiaplan_parking_green is the city-centre BLUE zone (синя зона).
-- sofiaplan_parking_blue  is the outer GREEN zone (зелена зона).
-- The color labels were swapped in migration 000015 — correct them here.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'blue'::text  AS color, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'green'::text AS color, geom FROM sofiaplan_parking_blue;
