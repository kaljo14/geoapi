-- name: ListPlaces :many
SELECT place_id, name, address, lat, lng, rating, business_status,
       website, formatted_phone_number, international_phone_number,
       opening_hours, reviews, editorial_summary, photos, types,
       price_level, user_ratings_total, utc_offset_minutes,
       google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
       reservable, takeout, wheelchair_accessible, category, tags, scraped_at
FROM places
ORDER BY name;

-- name: ListPlacesByCategory :many
SELECT place_id, name, address, lat, lng, rating, business_status,
       website, formatted_phone_number, international_phone_number,
       opening_hours, reviews, editorial_summary, photos, types,
       price_level, user_ratings_total, utc_offset_minutes,
       google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
       reservable, takeout, wheelchair_accessible, category, tags, scraped_at
FROM places
WHERE category = $1
ORDER BY name;

-- name: ListPlacesByCategoryAndTag :many
SELECT place_id, name, address, lat, lng, rating, business_status,
       website, formatted_phone_number, international_phone_number,
       opening_hours, reviews, editorial_summary, photos, types,
       price_level, user_ratings_total, utc_offset_minutes,
       google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
       reservable, takeout, wheelchair_accessible, category, tags, scraped_at
FROM places
WHERE category = $1
  AND (',' || tags || ',') LIKE ('%,' || $2::text || ',%')
ORDER BY name;

-- name: ListPlacesByTag :many
SELECT place_id, name, address, lat, lng, rating, business_status,
       website, formatted_phone_number, international_phone_number,
       opening_hours, reviews, editorial_summary, photos, types,
       price_level, user_ratings_total, utc_offset_minutes,
       google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
       reservable, takeout, wheelchair_accessible, category, tags, scraped_at
FROM places
WHERE (',' || tags || ',') LIKE ('%,' || $1::text || ',%')
ORDER BY name;

-- name: GetPlace :one
SELECT place_id, name, address, lat, lng, rating, business_status,
       website, formatted_phone_number, international_phone_number,
       opening_hours, reviews, editorial_summary, photos, types,
       price_level, user_ratings_total, utc_offset_minutes,
       google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
       reservable, takeout, wheelchair_accessible, category, tags, scraped_at
FROM places
WHERE place_id = $1;

-- name: CreatePlace :one
INSERT INTO places (place_id, name, address, lat, lng, rating, business_status, category, location)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, ST_SetSRID(ST_MakePoint($5, $4), 4326))
RETURNING place_id, name, address, lat, lng, rating, business_status, category,
          website, formatted_phone_number, international_phone_number,
          opening_hours, reviews, editorial_summary, photos, types,
          price_level, user_ratings_total, utc_offset_minutes,
          google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
          reservable, takeout, wheelchair_accessible, tags, scraped_at;

-- name: UpdatePlace :one
UPDATE places
SET name            = $1,
    address         = $2,
    lat             = $3,
    lng             = $4,
    rating          = $5,
    business_status = $6,
    category        = $7,
    location        = ST_SetSRID(ST_MakePoint($4, $3), 4326)
WHERE place_id = $8
RETURNING place_id, name, address, lat, lng, rating, business_status, category,
          website, formatted_phone_number, international_phone_number,
          opening_hours, reviews, editorial_summary, photos, types,
          price_level, user_ratings_total, utc_offset_minutes,
          google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
          reservable, takeout, wheelchair_accessible, tags, scraped_at;

-- name: DeletePlace :exec
DELETE FROM places WHERE place_id = $1;

-- name: GetMetroStops :many
SELECT stop_id, stop_name, ST_Y(geom) AS stop_lat, ST_X(geom) AS stop_lon
FROM gtfs_stops
WHERE geom IS NOT NULL
ORDER BY stop_name;

-- Saturation and heatmap queries are executed as raw pgx queries in analytics.go
-- due to complex PostGIS CTE expressions that sqlc cannot type-check.

-- name: UpsertPlace :one
INSERT INTO places (
    place_id, name, address, lat, lng, rating, business_status, category,
    website, formatted_phone_number, international_phone_number,
    opening_hours, reviews, editorial_summary, photos, types,
    price_level, user_ratings_total, utc_offset_minutes,
    google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
    reservable, takeout, wheelchair_accessible, tags, scraped_at,
    location
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11,
    $12, $13, $14, $15, $16,
    $17, $18, $19,
    $20, $21, $22, $23, $24,
    $25, $26, $27, $28, $29,
    ST_SetSRID(ST_MakePoint($5, $4), 4326)
)
ON CONFLICT (place_id) DO UPDATE SET
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    lat = EXCLUDED.lat,
    lng = EXCLUDED.lng,
    rating = EXCLUDED.rating,
    business_status = EXCLUDED.business_status,
    category = EXCLUDED.category,
    website = EXCLUDED.website,
    formatted_phone_number = EXCLUDED.formatted_phone_number,
    international_phone_number = EXCLUDED.international_phone_number,
    opening_hours = EXCLUDED.opening_hours,
    reviews = EXCLUDED.reviews,
    editorial_summary = EXCLUDED.editorial_summary,
    photos = EXCLUDED.photos,
    types = EXCLUDED.types,
    price_level = EXCLUDED.price_level,
    user_ratings_total = EXCLUDED.user_ratings_total,
    utc_offset_minutes = EXCLUDED.utc_offset_minutes,
    google_maps_url = EXCLUDED.google_maps_url,
    icon_url = EXCLUDED.icon_url,
    curbside_pickup = EXCLUDED.curbside_pickup,
    delivery = EXCLUDED.delivery,
    dine_in = EXCLUDED.dine_in,
    reservable = EXCLUDED.reservable,
    takeout = EXCLUDED.takeout,
    wheelchair_accessible = EXCLUDED.wheelchair_accessible,
    tags = EXCLUDED.tags,
    scraped_at = EXCLUDED.scraped_at,
    location = EXCLUDED.location
RETURNING place_id, name, address, lat, lng, rating, business_status, category,
          website, formatted_phone_number, international_phone_number,
          opening_hours, reviews, editorial_summary, photos, types,
          price_level, user_ratings_total, utc_offset_minutes,
          google_maps_url, icon_url, curbside_pickup, delivery, dine_in,
          reservable, takeout, wheelchair_accessible, tags, scraped_at;

-- name: ListPlacesNeedingEnrichment :many
SELECT place_id, name, lat, lng
FROM places
WHERE website IS NULL
   OR formatted_phone_number IS NULL
ORDER BY scraped_at ASC NULLS FIRST
LIMIT 100;

-- name: UpdatePlaceEnrichment :exec
UPDATE places
SET website                    = $2,
    formatted_phone_number     = $3,
    international_phone_number = $4,
    opening_hours              = $5,
    reviews                    = $6,
    editorial_summary          = $7,
    photos                     = $8,
    types                      = $9,
    price_level                = $10,
    user_ratings_total         = $11,
    utc_offset_minutes         = $12,
    google_maps_url            = $13,
    icon_url                   = $14,
    curbside_pickup            = $15,
    delivery                   = $16,
    dine_in                    = $17,
    reservable                 = $18,
    takeout                    = $19,
    wheelchair_accessible      = $20
WHERE place_id = $1;
