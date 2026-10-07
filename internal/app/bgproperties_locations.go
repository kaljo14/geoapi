package app

import (
	"context"
	"fmt"
)

// BgPropertiesLocation is the JSON response type for bulgarianproperties.com scraped listings.
type BgPropertiesLocation struct {
	URL          string   `json:"url"`
	PropertyID   *int64   `json:"property_id,omitempty"`
	Title        *string  `json:"title,omitempty"`
	Status       *string  `json:"status,omitempty"`
	PropertyType *string  `json:"property_type,omitempty"`
	Neighborhood *string  `json:"neighborhood,omitempty"`
	AreaSqm      *float64 `json:"area_sqm,omitempty"`
	PriceEur     *float64 `json:"price_eur,omitempty"`
	PricePerSqm  *float64 `json:"price_per_sqm,omitempty"`
	Floor        *int16   `json:"floor,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Subtitle     *string  `json:"subtitle,omitempty"`
	Description  *string  `json:"description,omitempty"`
	ImageURL     *string  `json:"image_url,omitempty"`
	AgentName    *string  `json:"agent_name,omitempty"`
	AgentRole    *string  `json:"agent_role,omitempty"`
	Lat          float64  `json:"lat"`
	Lng          float64  `json:"lng"`
	GeoSource    *string  `json:"geo_source,omitempty"`
	ScrapedAt    *string  `json:"scraped_at,omitempty"`
}

func (a *App) ListBgPropertiesLocations(ctx context.Context) ([]BgPropertiesLocation, error) {
	rows, err := a.store.ListBgPropertiesLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("list bgproperties locations: %w", err)
	}
	listings := make([]BgPropertiesLocation, 0, len(rows))
	for _, r := range rows {
		l := BgPropertiesLocation{
			URL:  r.URL,
			Tags: r.Tags,
			Lat:  r.Lat.Float64,
			Lng:  r.Lng.Float64,
		}
		if r.PropertyID.Valid {
			v := r.PropertyID.Int64
			l.PropertyID = &v
		}
		if r.Title.Valid {
			l.Title = &r.Title.String
		}
		if r.Status.Valid {
			l.Status = &r.Status.String
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
		if r.Subtitle.Valid {
			l.Subtitle = &r.Subtitle.String
		}
		if r.Description.Valid {
			l.Description = &r.Description.String
		}
		if r.ImageURL.Valid {
			l.ImageURL = &r.ImageURL.String
		}
		if r.AgentName.Valid {
			l.AgentName = &r.AgentName.String
		}
		if r.AgentRole.Valid {
			l.AgentRole = &r.AgentRole.String
		}
		if r.GeoSource.Valid {
			l.GeoSource = &r.GeoSource.String
		}
		l.ScrapedAt = ptime(r.ScrapedAt)
		listings = append(listings, l)
	}
	return listings, nil
}
