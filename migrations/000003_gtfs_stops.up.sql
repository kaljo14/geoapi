CREATE TABLE IF NOT EXISTS gtfs_stops (
    stop_id   TEXT PRIMARY KEY,
    stop_code TEXT,
    stop_name TEXT,
    stop_desc TEXT,
    stop_lat  FLOAT,
    stop_lon  FLOAT,
    geom      GEOMETRY(Point, 4326)
);

CREATE INDEX IF NOT EXISTS idx_gtfs_stops_geom ON gtfs_stops USING GIST(geom);
