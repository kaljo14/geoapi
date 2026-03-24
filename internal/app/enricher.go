package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/store"
)

type placeDetailsResponse struct {
	Result placeDetails `json:"result"`
	Status string       `json:"status"`
}

type placeDetails struct {
	Website                  string   `json:"website"`
	FormattedPhoneNumber     string   `json:"formatted_phone_number"`
	InternationalPhoneNumber string   `json:"international_phone_number"`
	OpeningHours             *struct {
		WeekdayText []string `json:"weekday_text"`
	} `json:"opening_hours"`
	EditorialSummary *struct {
		Overview string `json:"overview"`
	} `json:"editorial_summary"`
	Types            []string `json:"types"`
	PriceLevel       int      `json:"price_level"`
	UserRatingsTotal int      `json:"user_ratings_total"`
	UTCOffset        int      `json:"utc_offset"`
	URL              string   `json:"url"`
	Icon             string   `json:"icon"`
}

func (a *App) StartEnricher(ctx context.Context) error {
	if a.apiKey == "" {
		return fmt.Errorf("GOOGLE_API_KEY not set")
	}

	go func() {
		a.logger.Info("enricher started")
		places, err := a.store.ListPlacesNeedingEnrichment(ctx)
		if err != nil {
			a.logger.Error("enricher: list places failed", zap.Error(err))
			return
		}
		enriched := 0
		for _, p := range places {
			details, err := a.fetchPlaceDetails(p.PlaceID)
			if err != nil {
				a.logger.Warn("enricher: fetch details failed", zap.String("place_id", p.PlaceID), zap.Error(err))
				continue
			}
			if err := a.updateEnrichment(ctx, p.PlaceID, details); err != nil {
				a.logger.Warn("enricher: update failed", zap.String("place_id", p.PlaceID), zap.Error(err))
				continue
			}
			enriched++
			time.Sleep(200 * time.Millisecond)
		}
		a.logger.Info("enricher finished", zap.Int("enriched", enriched))
	}()
	return nil
}

func (a *App) fetchPlaceDetails(placeID string) (placeDetails, error) {
	fields := "website,formatted_phone_number,international_phone_number,opening_hours,editorial_summary,types,price_level,user_ratings_total,utc_offset,url,icon"
	params := url.Values{
		"place_id": {placeID},
		"fields":   {fields},
		"key":      {a.apiKey},
	}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/details/json?" + params.Encode())
	if err != nil {
		return placeDetails{}, fmt.Errorf("details request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp placeDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return placeDetails{}, fmt.Errorf("decode details response: %w", err)
	}
	if apiResp.Status != "OK" {
		return placeDetails{}, fmt.Errorf("place details status: %s", apiResp.Status)
	}
	return apiResp.Result, nil
}

func (a *App) updateEnrichment(ctx context.Context, placeID string, d placeDetails) error {
	var openingHours string
	if d.OpeningHours != nil {
		for i, line := range d.OpeningHours.WeekdayText {
			if i > 0 {
				openingHours += "\n"
			}
			openingHours += line
		}
	}
	var editorial string
	if d.EditorialSummary != nil {
		editorial = d.EditorialSummary.Overview
	}

	return a.store.UpdatePlaceEnrichment(ctx, store.UpdatePlaceEnrichmentParams{
		PlaceID:                  placeID,
		Website:                  pgtype.Text{String: d.Website, Valid: d.Website != ""},
		FormattedPhoneNumber:     pgtype.Text{String: d.FormattedPhoneNumber, Valid: d.FormattedPhoneNumber != ""},
		InternationalPhoneNumber: pgtype.Text{String: d.InternationalPhoneNumber, Valid: d.InternationalPhoneNumber != ""},
		OpeningHours:             pgtype.Text{String: openingHours, Valid: openingHours != ""},
		EditorialSummary:         pgtype.Text{String: editorial, Valid: editorial != ""},
		Types:                    pgtype.Text{String: joinStrings(d.Types), Valid: len(d.Types) > 0},
		PriceLevel:               pgtype.Int4{Int32: int32(d.PriceLevel), Valid: d.PriceLevel != 0},
		UserRatingsTotal:         pgtype.Int4{Int32: int32(d.UserRatingsTotal), Valid: d.UserRatingsTotal != 0},
		UtcOffsetMinutes:         pgtype.Int4{Int32: int32(d.UTCOffset), Valid: d.UTCOffset != 0},
		GoogleMapsUrl:            pgtype.Text{String: d.URL, Valid: d.URL != ""},
		IconUrl:                  pgtype.Text{String: d.Icon, Valid: d.Icon != ""},
	})
}
