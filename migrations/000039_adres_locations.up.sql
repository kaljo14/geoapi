CREATE TABLE IF NOT EXISTS adres_locations (
    offer_id        BIGINT PRIMARY KEY,
    url             TEXT NOT NULL,
    property_type   TEXT,
    neighborhood    TEXT,
    area_sqm        NUMERIC(10,2),
    price_eur       NUMERIC(10,2),
    price_per_sqm   NUMERIC(10,2),
    floor           SMALLINT,
    address_text    TEXT,
    lat             DOUBLE PRECISION,
    lng             DOUBLE PRECISION,
    geo_source      TEXT,
    scraped_at      TIMESTAMPTZ DEFAULT NOW(),
    location        GEOMETRY(Point, 4326)
);

CREATE INDEX IF NOT EXISTS idx_adres_locations_geom ON adres_locations USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_adres_locations_neighborhood ON adres_locations(neighborhood);
