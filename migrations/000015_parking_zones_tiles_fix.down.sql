-- Restore view with properties column (pre-fix state).
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'green' AS color, properties, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'blue'  AS color, properties, geom FROM sofiaplan_parking_blue;
