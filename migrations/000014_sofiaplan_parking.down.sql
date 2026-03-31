-- Restore parking_zones_tiles to point back at the OSM parking_zones table.
DROP VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT
    zone_id,
    name,
    color,
    geom
FROM parking_zones;

DROP TABLE IF EXISTS sofiaplan_parking_green;
DROP TABLE IF EXISTS sofiaplan_parking_blue;
DROP TABLE IF EXISTS sofiaplan_parking_zones;
