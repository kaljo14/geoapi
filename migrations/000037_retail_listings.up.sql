CREATE TABLE IF NOT EXISTS retail_listings (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title        TEXT NOT NULL,
    address      TEXT,
    lat          FLOAT NOT NULL,
    lng          FLOAT NOT NULL,
    location     GEOMETRY(Point, 4326),
    size_sqm     FLOAT,
    price_eur    FLOAT,
    listing_url  TEXT,
    is_exact     BOOLEAN NOT NULL DEFAULT TRUE,
    created_by   TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_retail_listings_location ON retail_listings USING GIST(location);
