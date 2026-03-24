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
SELECT stop_id, stop_name, stop_lat, stop_lon
FROM gtfs_stops
WHERE geom IS NOT NULL
ORDER BY stop_name;

-- Saturation and heatmap queries are executed as raw pgx queries in app.go
-- due to complex PostGIS CTE expressions that sqlc cannot type-check.
