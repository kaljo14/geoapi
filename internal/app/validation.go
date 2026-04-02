package app

import (
	"errors"
	"fmt"
)

// Coordinate bounds for Sofia area (generous buffer around the city).
const (
	minLat = 42.0
	maxLat = 43.5
	minLng = 22.5
	maxLng = 24.5
)

// ValidateCoordinates checks that lat/lng are within the Sofia metro area.
func ValidateCoordinates(lat, lng float64) error {
	if lat < minLat || lat > maxLat || lng < minLng || lng > maxLng {
		return fmt.Errorf("coordinates (%.6f, %.6f) outside Sofia area bounds", lat, lng)
	}
	return nil
}

// ValidateCreatePlace checks required fields for place creation.
func ValidateCreatePlace(name string, lat, lng float64) error {
	if name == "" {
		return errors.New("name is required")
	}
	return ValidateCoordinates(lat, lng)
}

// ValidateRadius checks that a radius value is positive and reasonable.
func ValidateRadius(radius float64) error {
	if radius <= 0 {
		return errors.New("radius must be positive")
	}
	if radius > 50000 {
		return errors.New("radius exceeds maximum (50km)")
	}
	return nil
}

// ValidateBBox checks that a bounding box has valid ordering and reasonable size.
func ValidateBBox(minLat, minLng, maxLat, maxLng float64) error {
	if minLat >= maxLat {
		return errors.New("minLat must be less than maxLat")
	}
	if minLng >= maxLng {
		return errors.New("minLng must be less than maxLng")
	}
	return nil
}
