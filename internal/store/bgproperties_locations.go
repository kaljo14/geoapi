package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listBgPropertiesLocations = `-- name: ListBgPropertiesLocations :many
SELECT url, property_id, title, status, property_type, neighborhood,
       area_sqm, price_eur, price_per_sqm, floor, tags,
       subtitle, description, image_url, agent_name, agent_role,
       lat, lng, geo_source, scraped_at
FROM bgproperties_locations
WHERE lat IS NOT NULL AND lng IS NOT NULL
ORDER BY scraped_at DESC
`

type BgPropertiesLocationRow struct {
	URL          string
	PropertyID   pgtype.Int8
	Title        pgtype.Text
	Status       pgtype.Text
	PropertyType pgtype.Text
	Neighborhood pgtype.Text
	AreaSqm      pgtype.Numeric
	PriceEur     pgtype.Numeric
	PricePerSqm  pgtype.Numeric
	Floor        pgtype.Int2
	Tags         []string
	Subtitle     pgtype.Text
	Description  pgtype.Text
	ImageURL     pgtype.Text
	AgentName    pgtype.Text
	AgentRole    pgtype.Text
	Lat          pgtype.Float8
	Lng          pgtype.Float8
	GeoSource    pgtype.Text
	ScrapedAt    pgtype.Timestamptz
}

func (q *Queries) ListBgPropertiesLocations(ctx context.Context) ([]BgPropertiesLocationRow, error) {
	rows, err := q.db.Query(ctx, listBgPropertiesLocations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BgPropertiesLocationRow
	for rows.Next() {
		var i BgPropertiesLocationRow
		if err := rows.Scan(
			&i.URL,
			&i.PropertyID,
			&i.Title,
			&i.Status,
			&i.PropertyType,
			&i.Neighborhood,
			&i.AreaSqm,
			&i.PriceEur,
			&i.PricePerSqm,
			&i.Floor,
			&i.Tags,
			&i.Subtitle,
			&i.Description,
			&i.ImageURL,
			&i.AgentName,
			&i.AgentRole,
			&i.Lat,
			&i.Lng,
			&i.GeoSource,
			&i.ScrapedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
