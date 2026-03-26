CREATE TABLE IF NOT EXISTS osm_pois (
    osm_id   BIGINT PRIMARY KEY,
    lat      FLOAT  NOT NULL,
    lng      FLOAT  NOT NULL,
    geom     GEOMETRY(Point, 4326),
    name     TEXT,
    amenity  TEXT,
    shop     TEXT,
    tourism  TEXT,
    leisure  TEXT,
    category TEXT
);

CREATE INDEX IF NOT EXISTS idx_osm_pois_geom     ON osm_pois USING GIST(geom);
CREATE INDEX IF NOT EXISTS idx_osm_pois_category ON osm_pois(category);
