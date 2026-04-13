package app

import (
	"context"
	"fmt"
)

// AdresLocation is the JSON response type for address.bg scraped listings.
type AdresLocation struct {
	OfferID      int64    `json:"offer_id"`
	URL          string   `json:"url"`
	PropertyType *string  `json:"property_type,omitempty"`
	Neighborhood *string  `json:"neighborhood,omitempty"`
	AreaSqm      *float64 `json:"area_sqm,omitempty"`
	PriceEur     *float64 `json:"price_eur,omitempty"`
	PricePerSqm  *float64 `json:"price_per_sqm,omitempty"`
	Floor        *int16   `json:"floor,omitempty"`
	AddressText  *string  `json:"address_text,omitempty"`
	Lat          float64  `json:"lat"`
	Lng          float64  `json:"lng"`
	GeoSource    *string  `json:"geo_source,omitempty"`
	ScrapedAt    *string  `json:"scraped_at,omitempty"`
}

func (a *App) ListAdresLocations(ctx context.Context) ([]AdresLocation, error) {
	rows, err := a.store.ListAdresLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("list adres locations: %w", err)
	}
	listings := make([]AdresLocation, 0, len(rows))
	for _, r := range rows {
		l := AdresLocation{
			OfferID: r.OfferID,
			URL:     r.URL,
			Lat:     r.Lat.Float64,
			Lng:     r.Lng.Float64,
		}
		if r.PropertyType.Valid {
			l.PropertyType = &r.PropertyType.String
		}
		if r.Neighborhood.Valid {
			l.Neighborhood = &r.Neighborhood.String
		}
		l.AreaSqm = pnumeric(r.AreaSqm)
		l.PriceEur = pnumeric(r.PriceEur)
		l.PricePerSqm = pnumeric(r.PricePerSqm)
		if r.Floor.Valid {
			v := r.Floor.Int16
			l.Floor = &v
		}
		if r.AddressText.Valid {
			l.AddressText = &r.AddressText.String
		}
		if r.GeoSource.Valid {
			l.GeoSource = &r.GeoSource.String
		}
		l.ScrapedAt = ptime(r.ScrapedAt)
		listings = append(listings, l)
	}
	return listings, nil
}
