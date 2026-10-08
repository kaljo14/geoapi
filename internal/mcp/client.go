package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client wraps HTTP calls to the GeoPulse REST API.
type Client struct {
	base   string // e.g. "http://localhost:8080"
	client *http.Client
}

// NewClient creates a Client pointing at the given base URL.
func NewClient(baseURL string) *Client {
	return &Client{
		base: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// get performs a GET request and JSON-decodes the response into dst.
func (c *Client) get(path string, query url.Values, dst any) error {
	u := c.base + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	resp, err := c.client.Get(u)
	if err != nil {
		return fmt.Errorf("http get %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d from %s: %s", resp.StatusCode, path, body)
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}

// ListPlaces calls GET /api/places with optional category and tag filters.
func (c *Client) ListPlaces(category, tag string) (json.RawMessage, error) {
	q := url.Values{}
	if category != "" {
		q.Set("category", category)
	}
	if tag != "" {
		q.Set("tag", tag)
	}
	var result json.RawMessage
	return result, c.get("/api/places", q, &result)
}

// GetSaturation calls GET /api/saturation.
func (c *Client) GetSaturation(lat, lng, radius float64, category string) (json.RawMessage, error) {
	q := url.Values{
		"lat":    {fmt.Sprintf("%f", lat)},
		"lng":    {fmt.Sprintf("%f", lng)},
		"radius": {fmt.Sprintf("%f", radius)},
	}
	if category != "" {
		q.Set("category", category)
	}
	var result json.RawMessage
	return result, c.get("/api/saturation", q, &result)
}

// GetHeatmap calls GET /api/heatmap.
func (c *Client) GetHeatmap(minLat, minLng, maxLat, maxLng, cellSize float64, category string) (json.RawMessage, error) {
	q := url.Values{
		"min_lat":   {fmt.Sprintf("%f", minLat)},
		"min_lng":   {fmt.Sprintf("%f", minLng)},
		"max_lat":   {fmt.Sprintf("%f", maxLat)},
		"max_lng":   {fmt.Sprintf("%f", maxLng)},
		"cell_size": {fmt.Sprintf("%f", cellSize)},
	}
	if category != "" {
		q.Set("category", category)
	}
	var result json.RawMessage
	return result, c.get("/api/heatmap", q, &result)
}

// GetLocationContext calls GET /api/sofiaplan/context.
func (c *Client) GetLocationContext(lat, lng float64) (json.RawMessage, error) {
	q := url.Values{
		"lat": {fmt.Sprintf("%f", lat)},
		"lng": {fmt.Sprintf("%f", lng)},
	}
	var result json.RawMessage
	return result, c.get("/api/sofiaplan/context", q, &result)
}

// ListRetailListings calls GET /api/retail-listings.
func (c *Client) ListRetailListings() (json.RawMessage, error) {
	var result json.RawMessage
	return result, c.get("/api/retail-listings", nil, &result)
}
