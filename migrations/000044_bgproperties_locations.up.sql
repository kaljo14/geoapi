CREATE TABLE IF NOT EXISTS bgproperties_locations (
    url             TEXT PRIMARY KEY,
    property_id     BIGINT,
    title           TEXT,
    status          TEXT,
    property_type   TEXT,
    neighborhood    TEXT,
    area_sqm        NUMERIC(10,2),
    price_eur       NUMERIC(10,2),
    price_per_sqm   NUMERIC(10,2),
    floor           SMALLINT,
    tags            TEXT[],
    subtitle        TEXT,
    description     TEXT,
    image_url       TEXT,
    agent_name      TEXT,
    agent_role      TEXT,
    lat             DOUBLE PRECISION,
    lng             DOUBLE PRECISION,
    geo_source      TEXT,
    scraped_at      TIMESTAMPTZ DEFAULT NOW(),
    location        GEOMETRY(Point, 4326)
);

CREATE INDEX IF NOT EXISTS idx_bgproperties_locations_geom
    ON bgproperties_locations USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_bgproperties_locations_status
    ON bgproperties_locations(status);
CREATE INDEX IF NOT EXISTS idx_bgproperties_locations_neighborhood
    ON bgproperties_locations(neighborhood);
