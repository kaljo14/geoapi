package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DatabaseURL  string
	Port         string
	GoogleAPIKey string

	// Scraper defaults (Sofia city centre)
	ScrapeLat    float64
	ScrapeLng    float64
	ScrapeRadius float64
	ScrapeTypes  string
	ScrapeQuery  string // when set, use Text Search API instead of Nearby Search
	ScrapeTags   string // comma-separated tags to apply to scraped places

}

// Load reads configuration from environment variables and validates required fields.
func Load() (*Config, error) {
	c := &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Port:         envString("PORT", "8080"),
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
		ScrapeLat:    envFloat("SCRAPE_LAT", 42.6977),
		ScrapeLng:    envFloat("SCRAPE_LNG", 23.3219),
		ScrapeRadius: envFloat("SCRAPE_RADIUS", 5000),
		ScrapeTypes:  envString("SCRAPE_TYPES", "restaurant"),
		ScrapeQuery:  os.Getenv("SCRAPE_QUERY"),
		ScrapeTags:   os.Getenv("SCRAPE_TAGS"),
	}
	if c.Port == "" {
		return nil, fmt.Errorf("PORT must not be empty")
	}
	return c, nil
}

func envString(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envFloat(key string, defaultVal float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultVal
}
