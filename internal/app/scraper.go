package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/store"
)

type placesAPIResponse struct {
	Results       []placeResult `json:"results"`
	NextPageToken string        `json:"next_page_token"`
	Status        string        `json:"status"`
}

type placeResult struct {
	PlaceID        string   `json:"place_id"`
	Name           string   `json:"name"`
	BusinessStatus string   `json:"business_status"`
	Rating         float64  `json:"rating"`
	Types          []string `json:"types"`
	Geometry       struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
	Vicinity string `json:"vicinity"`
}

func (a *App) StartScraper(ctx context.Context) error {
	if a.apiKey == "" {
		return fmt.Errorf("GOOGLE_API_KEY not set")
	}

	lat := envFloat("SCRAPE_LAT", 42.6977)
	lng := envFloat("SCRAPE_LNG", 23.3219)
	radius := envFloat("SCRAPE_RADIUS", 5000)
	types := os.Getenv("SCRAPE_TYPES")
	if types == "" {
		types = "restaurant"
	}

	go func() {
		a.logger.Info("scraper started", zap.Float64("lat", lat), zap.Float64("lng", lng))
		scraped := 0
		pageToken := ""

		for {
			results, nextToken, err := a.fetchNearbyPlaces(lat, lng, radius, types, pageToken)
			if err != nil {
				a.logger.Error("scraper fetch failed", zap.Error(err))
				return
			}
			for _, r := range results {
				if err := a.upsertScrapedPlace(ctx, r, types); err != nil {
					a.logger.Warn("upsert failed", zap.String("place_id", r.PlaceID), zap.Error(err))
				} else {
					scraped++
				}
			}
			if nextToken == "" {
				break
			}
			pageToken = nextToken
			time.Sleep(2 * time.Second) // required by Google API
		}
		a.logger.Info("scraper finished", zap.Int("scraped", scraped))
	}()
	return nil
}

func (a *App) fetchNearbyPlaces(lat, lng, radius float64, placeType, pageToken string) ([]placeResult, string, error) {
	params := url.Values{
		"location": {fmt.Sprintf("%f,%f", lat, lng)},
		"radius":   {fmt.Sprintf("%.0f", radius)},
		"type":     {placeType},
		"key":      {a.apiKey},
	}
	if pageToken != "" {
		params.Set("pagetoken", pageToken)
	}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?" + params.Encode())
	if err != nil {
		return nil, "", fmt.Errorf("nearby search request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp placesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, "", fmt.Errorf("decode nearby search response: %w", err)
	}
	if apiResp.Status != "OK" && apiResp.Status != "ZERO_RESULTS" {
		return nil, "", fmt.Errorf("places API status: %s", apiResp.Status)
	}
	return apiResp.Results, apiResp.NextPageToken, nil
}

func (a *App) upsertScrapedPlace(ctx context.Context, r placeResult, category string) error {
	now := time.Now()
	_, err := a.store.UpsertPlace(ctx, store.UpsertPlaceParams{
		PlaceID:        r.PlaceID,
		Name:           r.Name,
		Address:        pgtype.Text{String: r.Vicinity, Valid: r.Vicinity != ""},
		Lat:            pgtype.Float8{Float64: r.Geometry.Location.Lat, Valid: true},
		Lng:            pgtype.Float8{Float64: r.Geometry.Location.Lng, Valid: true},
		Rating:         pgtype.Float8{Float64: r.Rating, Valid: r.Rating != 0},
		BusinessStatus: pgtype.Text{String: r.BusinessStatus, Valid: r.BusinessStatus != ""},
		Category:       pgtype.Text{String: category, Valid: true},
		ScrapedAt:      pgtype.Timestamptz{Time: now, Valid: true},
	})
	return err
}
