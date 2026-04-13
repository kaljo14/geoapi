package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listAdresLocations = `-- name: ListAdresLocations :many
SELECT offer_id, url, property_type, neighborhood, area_sqm, price_eur, price_per_sqm,
       floor, address_text, lat, lng, geo_source, scraped_at
FROM adres_locations
WHERE lat IS NOT NULL AND lng IS NOT NULL
ORDER BY scraped_at DESC
`

type AdresLocationRow struct {
	OfferID      int64
	URL          string
	PropertyType pgtype.Text
	Neighborhood pgtype.Text
	AreaSqm      pgtype.Numeric
	PriceEur     pgtype.Numeric
	PricePerSqm  pgtype.Numeric
	Floor        pgtype.Int2
	AddressText  pgtype.Text
	Lat          pgtype.Float8
	Lng          pgtype.Float8
	GeoSource    pgtype.Text
	ScrapedAt    pgtype.Timestamptz
}

func (q *Queries) ListAdresLocations(ctx context.Context) ([]AdresLocationRow, error) {
	rows, err := q.db.Query(ctx, listAdresLocations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdresLocationRow
	for rows.Next() {
		var i AdresLocationRow
		if err := rows.Scan(
			&i.OfferID,
			&i.URL,
			&i.PropertyType,
			&i.Neighborhood,
			&i.AreaSqm,
			&i.PriceEur,
			&i.PricePerSqm,
			&i.Floor,
			&i.AddressText,
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
