package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/store"
)

func (a *App) ListPlaces(ctx context.Context, category, tag string) ([]generated.Place, error) {
	switch {
	case category != "" && tag != "":
		rows, err := a.store.ListPlacesByCategoryAndTag(ctx, store.ListPlacesByCategoryAndTagParams{
			Category: pgtype.Text{String: category, Valid: true},
			Column2:  tag,
		})
		if err != nil {
			return nil, fmt.Errorf("list places by category and tag: %w", err)
		}
		places := make([]generated.Place, 0, len(rows))
		for _, r := range rows {
			places = append(places, listPlacesByCategoryAndTagRowToPlace(r))
		}
		return places, nil
	case category != "":
		rows, err := a.store.ListPlacesByCategory(ctx, pgtype.Text{String: category, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("list places by category: %w", err)
		}
		places := make([]generated.Place, 0, len(rows))
		for _, r := range rows {
			places = append(places, listPlacesByCategoryRowToPlace(r))
		}
		return places, nil
	case tag != "":
		rows, err := a.store.ListPlacesByTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("list places by tag: %w", err)
		}
		places := make([]generated.Place, 0, len(rows))
		for _, r := range rows {
			places = append(places, listPlacesByTagRowToPlace(r))
		}
		return places, nil
	default:
		rows, err := a.store.ListPlaces(ctx)
		if err != nil {
			return nil, fmt.Errorf("list places: %w", err)
		}
		places := make([]generated.Place, 0, len(rows))
		for _, r := range rows {
			places = append(places, listPlacesRowToPlace(r))
		}
		return places, nil
	}
}

func (a *App) GetPlace(ctx context.Context, placeID string) (*generated.Place, error) {
	row, err := a.store.GetPlace(ctx, placeID)
	if err != nil {
		return nil, fmt.Errorf("get place: %w", err)
	}
	p := getPlaceRowToPlace(row)
	return &p, nil
}

func (a *App) CreatePlace(ctx context.Context, req generated.CreatePlaceRequest) (*generated.Place, error) {
	if err := ValidateCreatePlace(req.Name, req.Lat, req.Lng); err != nil {
		return nil, fmt.Errorf("validate create place: %w", err)
	}
	placeID := fmt.Sprintf("manual_%d", time.Now().UnixNano())
	row, err := a.store.CreatePlace(ctx, store.CreatePlaceParams{
		PlaceID:        placeID,
		Name:           req.Name,
		Address:        pgtype.Text{String: strVal(req.Address), Valid: req.Address != nil},
		Lat:            pgtype.Float8{Float64: req.Lat, Valid: true},
		Lng:            pgtype.Float8{Float64: req.Lng, Valid: true},
		Rating:         pgtype.Float8{Float64: f64Val(req.Rating), Valid: req.Rating != nil},
		BusinessStatus: pgtype.Text{String: strVal(req.BusinessStatus), Valid: req.BusinessStatus != nil},
		Category:       pgtype.Text{String: strVal(req.Category), Valid: req.Category != nil},
	})
	if err != nil {
		return nil, fmt.Errorf("create place: %w", err)
	}
	p := createPlaceRowToPlace(row)
	return &p, nil
}

func (a *App) UpdatePlace(ctx context.Context, placeID string, req generated.UpdatePlaceRequest) (*generated.Place, error) {
	if err := ValidateCreatePlace(req.Name, req.Lat, req.Lng); err != nil {
		return nil, fmt.Errorf("validate update place: %w", err)
	}
	row, err := a.store.UpdatePlace(ctx, store.UpdatePlaceParams{
		Name:           req.Name,
		Address:        pgtype.Text{String: strVal(req.Address), Valid: req.Address != nil},
		Lat:            pgtype.Float8{Float64: req.Lat, Valid: true},
		Lng:            pgtype.Float8{Float64: req.Lng, Valid: true},
		Rating:         pgtype.Float8{Float64: f64Val(req.Rating), Valid: req.Rating != nil},
		BusinessStatus: pgtype.Text{String: strVal(req.BusinessStatus), Valid: req.BusinessStatus != nil},
		Category:       pgtype.Text{String: strVal(req.Category), Valid: req.Category != nil},
		PlaceID:        placeID,
	})
	if err != nil {
		return nil, fmt.Errorf("update place: %w", err)
	}
	p := updatePlaceRowToPlace(row)
	return &p, nil
}

func (a *App) DeletePlace(ctx context.Context, placeID string) error {
	return a.store.DeletePlace(ctx, placeID)
}
