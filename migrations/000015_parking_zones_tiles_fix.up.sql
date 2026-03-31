-- Drop the raw JSONB column from the view — Martin cannot encode JSONB into MVT.
-- The only property needed for rendering is `color` ('green' | 'blue').
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'green'::text AS color, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'blue'::text  AS color, geom FROM sofiaplan_parking_blue;
