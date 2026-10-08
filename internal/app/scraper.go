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
	Vicinity         string `json:"vicinity"`
	FormattedAddress string `json:"formatted_address"`
}

func (r placeResult) address() string {
	if r.FormattedAddress != "" {
		return r.FormattedAddress
	}
	return r.Vicinity
}

func (a *App) StartScraper(ctx context.Context) error {
	if a.cfg.GoogleAPIKey == "" {
		return fmt.Errorf("GOOGLE_API_KEY not set")
	}

	lat := a.cfg.ScrapeLat
	lng := a.cfg.ScrapeLng
	radius := a.cfg.ScrapeRadius
	types := a.cfg.ScrapeTypes
	query := a.cfg.ScrapeQuery
	tags := a.cfg.ScrapeTags

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		bgCtx := a.ctx
		a.logger.Info("scraper started", zap.Float64("lat", lat), zap.Float64("lng", lng), zap.String("query", query))
		scraped := 0
		pageToken := ""

		for {
			select {
			case <-bgCtx.Done():
				a.logger.Info("scraper cancelled", zap.Int("scraped", scraped))
				return
			default:
			}

			var results []placeResult
			var nextToken string
			var err error

			if query != "" {
				results, nextToken, err = a.fetchTextSearchPlaces(query, lat, lng, radius, pageToken)
			} else {
				results, nextToken, err = a.fetchNearbyPlaces(lat, lng, radius, types, pageToken)
			}
			if err != nil {
				a.logger.Error("scraper fetch failed", zap.Error(err))
				return
			}
			for _, r := range results {
				if err := a.upsertScrapedPlace(bgCtx, r, types, tags); err != nil {
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
		"key":      {a.cfg.GoogleAPIKey},
	}
	if pageToken != "" {
		params.Set("pagetoken", pageToken)
	}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?" + params.Encode())
	if err != nil {
		return nil, "", fmt.Errorf("nearby search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var apiResp placesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, "", fmt.Errorf("decode nearby search response: %w", err)
	}
	if apiResp.Status != "OK" && apiResp.Status != "ZERO_RESULTS" {
		return nil, "", fmt.Errorf("places API status: %s", apiResp.Status)
	}
	return apiResp.Results, apiResp.NextPageToken, nil
}

func (a *App) fetchTextSearchPlaces(query string, lat, lng, radius float64, pageToken string) ([]placeResult, string, error) {
	params := url.Values{
		"query":    {query},
		"location": {fmt.Sprintf("%f,%f", lat, lng)},
		"radius":   {fmt.Sprintf("%.0f", radius)},
		"key":      {a.cfg.GoogleAPIKey},
	}
	if pageToken != "" {
		params.Set("pagetoken", pageToken)
	}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/textsearch/json?" + params.Encode())
	if err != nil {
		return nil, "", fmt.Errorf("text search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var apiResp placesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, "", fmt.Errorf("decode text search response: %w", err)
	}
	if apiResp.Status != "OK" && apiResp.Status != "ZERO_RESULTS" {
		return nil, "", fmt.Errorf("places API status: %s", apiResp.Status)
	}
	return apiResp.Results, apiResp.NextPageToken, nil
}

func (a *App) upsertScrapedPlace(ctx context.Context, r placeResult, category, tags string) error {
	now := time.Now()
	_, err := a.store.UpsertPlace(ctx, store.UpsertPlaceParams{
		PlaceID:        r.PlaceID,
		Name:           r.Name,
		Address:        pgtype.Text{String: r.address(), Valid: r.address() != ""},
		Lat:            pgtype.Float8{Float64: r.Geometry.Location.Lat, Valid: true},
		Lng:            pgtype.Float8{Float64: r.Geometry.Location.Lng, Valid: true},
		Rating:         pgtype.Float8{Float64: r.Rating, Valid: r.Rating != 0},
		BusinessStatus: pgtype.Text{String: r.BusinessStatus, Valid: r.BusinessStatus != ""},
		Category:       pgtype.Text{String: category, Valid: true},
		Tags:           pgtype.Text{String: tags, Valid: tags != ""},
		ScrapedAt:      pgtype.Timestamptz{Time: now, Valid: true},
	})
	return err
}
