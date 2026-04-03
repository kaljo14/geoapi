package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/store"
)

func (a *App) ListRetailListings(ctx context.Context) ([]generated.RetailListing, error) {
	rows, err := a.store.ListRetailListings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list retail listings: %w", err)
	}
	listings := make([]generated.RetailListing, 0, len(rows))
	for _, r := range rows {
		listings = append(listings, listRetailListingRowToListing(r))
	}
	return listings, nil
}

func (a *App) GetRetailListing(ctx context.Context, id string) (*generated.RetailListing, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	row, err := a.store.GetRetailListing(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get retail listing: %w", err)
	}
	l := getRetailListingRowToListing(row)
	return &l, nil
}

func (a *App) CreateRetailListing(ctx context.Context, req generated.CreateRetailListingRequest) (*generated.RetailListing, error) {
	if err := ValidateCreateRetailListing(req.Title, req.Lat, req.Lng); err != nil {
		return nil, fmt.Errorf("validate create retail listing: %w", err)
	}
	row, err := a.store.CreateRetailListing(ctx, store.CreateRetailListingParams{
		Title:         req.Title,
		Address:       pgtype.Text{String: strVal(req.Address), Valid: req.Address != nil},
		Lat:           pgtype.Float8{Float64: req.Lat, Valid: true},
		Lng:           pgtype.Float8{Float64: req.Lng, Valid: true},
		SizeSqm:       pgtype.Float8{Float64: f64Val(req.SizeSqm), Valid: req.SizeSqm != nil},
		PriceEur:      pgtype.Float8{Float64: f64Val(req.PriceEur), Valid: req.PriceEur != nil},
		ListingUrl:    pgtype.Text{String: strVal(req.ListingUrl), Valid: req.ListingUrl != nil},
		GoogleMapsUrl: pgtype.Text{String: strVal(req.GoogleMapsUrl), Valid: req.GoogleMapsUrl != nil},
		IsExact:       req.IsExact,
		CreatedBy:     pgtype.Text{},
	})
	if err != nil {
		return nil, fmt.Errorf("create retail listing: %w", err)
	}
	l := createRetailListingRowToListing(row)
	return &l, nil
}

func (a *App) UpdateRetailListing(ctx context.Context, id string, req generated.UpdateRetailListingRequest) (*generated.RetailListing, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	// Fetch existing to merge partial updates
	existing, err := a.store.GetRetailListing(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get retail listing for update: %w", err)
	}

	title := existing.Title
	if req.Title != nil {
		title = *req.Title
	}
	address := existing.Address
	if req.Address != nil {
		address = pgtype.Text{String: *req.Address, Valid: true}
	}
	lat := existing.Lat
	if req.Lat != nil {
		lat = pgtype.Float8{Float64: *req.Lat, Valid: true}
	}
	lng := existing.Lng
	if req.Lng != nil {
		lng = pgtype.Float8{Float64: *req.Lng, Valid: true}
	}
	sizeSqm := existing.SizeSqm
	if req.SizeSqm != nil {
		sizeSqm = pgtype.Float8{Float64: *req.SizeSqm, Valid: true}
	}
	priceEur := existing.PriceEur
	if req.PriceEur != nil {
		priceEur = pgtype.Float8{Float64: *req.PriceEur, Valid: true}
	}
	listingUrl := existing.ListingUrl
	if req.ListingUrl != nil {
		listingUrl = pgtype.Text{String: *req.ListingUrl, Valid: true}
	}
	googleMapsUrl := existing.GoogleMapsUrl
	if req.GoogleMapsUrl != nil {
		googleMapsUrl = pgtype.Text{String: *req.GoogleMapsUrl, Valid: true}
	}
	isExact := existing.IsExact
	if req.IsExact != nil {
		isExact = *req.IsExact
	}

	if err := ValidateCreateRetailListing(title, lat.Float64, lng.Float64); err != nil {
		return nil, fmt.Errorf("validate update retail listing: %w", err)
	}

	row, err := a.store.UpdateRetailListing(ctx, store.UpdateRetailListingParams{
		Title:         title,
		Address:       address,
		Lat:           lat,
		Lng:           lng,
		SizeSqm:       sizeSqm,
		PriceEur:      priceEur,
		ListingUrl:    listingUrl,
		GoogleMapsUrl: googleMapsUrl,
		IsExact:       isExact,
		ID:            uid,
	})
	if err != nil {
		return nil, fmt.Errorf("update retail listing: %w", err)
	}
	l := updateRetailListingRowToListing(row)
	return &l, nil
}

func (a *App) DeleteRetailListing(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return a.store.DeleteRetailListing(ctx, uid)
}

// --- helpers ---

func parseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return u, fmt.Errorf("parse uuid %q: %w", s, err)
	}
	return u, nil
}

func uuidStr(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// --- converters ---

func listRetailListingRowToListing(r store.ListRetailListingsRow) generated.RetailListing {
	return generated.RetailListing{
		Id:            uuidStr(r.ID),
		Title:         r.Title,
		Address:       pstr(r.Address),
		Lat:           r.Lat.Float64,
		Lng:           r.Lng.Float64,
		SizeSqm:       pf64(r.SizeSqm),
		PriceEur:      pf64(r.PriceEur),
		ListingUrl:    pstr(r.ListingUrl),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl),
		IsExact:       r.IsExact,
		CreatedBy:     pstr(r.CreatedBy),
		CreatedAt:     ptime(r.CreatedAt),
		UpdatedAt:     ptime(r.UpdatedAt),
	}
}

func getRetailListingRowToListing(r store.GetRetailListingRow) generated.RetailListing {
	return generated.RetailListing{
		Id:            uuidStr(r.ID),
		Title:         r.Title,
		Address:       pstr(r.Address),
		Lat:           r.Lat.Float64,
		Lng:           r.Lng.Float64,
		SizeSqm:       pf64(r.SizeSqm),
		PriceEur:      pf64(r.PriceEur),
		ListingUrl:    pstr(r.ListingUrl),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl),
		IsExact:       r.IsExact,
		CreatedBy:     pstr(r.CreatedBy),
		CreatedAt:     ptime(r.CreatedAt),
		UpdatedAt:     ptime(r.UpdatedAt),
	}
}

func createRetailListingRowToListing(r store.CreateRetailListingRow) generated.RetailListing {
	return generated.RetailListing{
		Id:            uuidStr(r.ID),
		Title:         r.Title,
		Address:       pstr(r.Address),
		Lat:           r.Lat.Float64,
		Lng:           r.Lng.Float64,
		SizeSqm:       pf64(r.SizeSqm),
		PriceEur:      pf64(r.PriceEur),
		ListingUrl:    pstr(r.ListingUrl),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl),
		IsExact:       r.IsExact,
		CreatedBy:     pstr(r.CreatedBy),
		CreatedAt:     ptime(r.CreatedAt),
		UpdatedAt:     ptime(r.UpdatedAt),
	}
}

func updateRetailListingRowToListing(r store.UpdateRetailListingRow) generated.RetailListing {
	return generated.RetailListing{
		Id:            uuidStr(r.ID),
		Title:         r.Title,
		Address:       pstr(r.Address),
		Lat:           r.Lat.Float64,
		Lng:           r.Lng.Float64,
		SizeSqm:       pf64(r.SizeSqm),
		PriceEur:      pf64(r.PriceEur),
		ListingUrl:    pstr(r.ListingUrl),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl),
		IsExact:       r.IsExact,
		CreatedBy:     pstr(r.CreatedBy),
		CreatedAt:     ptime(r.CreatedAt),
		UpdatedAt:     ptime(r.UpdatedAt),
	}
}
