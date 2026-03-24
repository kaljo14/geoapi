package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"

	"github.com/neofyis/geopulse/bff/internal/generated"
	"github.com/neofyis/geopulse/bff/internal/mocks"
	"github.com/neofyis/geopulse/bff/internal/store"
)

func TestListPlaces_NoFilter_CallsListPlaces(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().ListPlaces(mock.Anything).
		Return([]store.ListPlacesRow{{PlaceID: "p1", Name: "A"}}, nil)

	places, err := newApp(ms).ListPlaces(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 || places[0].PlaceId != "p1" {
		t.Errorf("unexpected places: %v", places)
	}
}

func TestListPlaces_CategoryOnly_CallsListPlacesByCategory(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().
		ListPlacesByCategory(mock.Anything, pgtype.Text{String: "cafe", Valid: true}).
		Return([]store.ListPlacesByCategoryRow{{PlaceID: "p2", Name: "B"}}, nil)

	places, err := newApp(ms).ListPlaces(context.Background(), "cafe", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 || places[0].PlaceId != "p2" {
		t.Errorf("unexpected places: %v", places)
	}
}

func TestListPlaces_TagOnly_CallsListPlacesByTag(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().ListPlacesByTag(mock.Anything, "rooftop").
		Return(nil, nil)

	_, err := newApp(ms).ListPlaces(context.Background(), "", "rooftop")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListPlaces_BothFilters_CallsListPlacesByCategoryAndTag(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().
		ListPlacesByCategoryAndTag(mock.Anything, store.ListPlacesByCategoryAndTagParams{
			Category: pgtype.Text{String: "cafe", Valid: true},
			Column2:  "rooftop",
		}).
		Return(nil, nil)

	_, err := newApp(ms).ListPlaces(context.Background(), "cafe", "rooftop")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListPlaces_StoreError_ReturnsError(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().ListPlaces(mock.Anything).Return(nil, errors.New("db down"))

	_, err := newApp(ms).ListPlaces(context.Background(), "", "")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreatePlace_GeneratesManualPlaceID(t *testing.T) {
	ms := mocks.NewMockStore(t)
	ms.EXPECT().
		CreatePlace(mock.Anything, mock.MatchedBy(func(arg store.CreatePlaceParams) bool {
			return strings.HasPrefix(arg.PlaceID, "manual_")
		})).
		RunAndReturn(func(_ context.Context, arg store.CreatePlaceParams) (store.CreatePlaceRow, error) {
			return store.CreatePlaceRow{PlaceID: arg.PlaceID, Name: arg.Name}, nil
		})

	cat := "cafe"
	place, err := newApp(ms).CreatePlace(context.Background(), generated.CreatePlaceRequest{
		Name:     "Test",
		Lat:      42.7,
		Lng:      23.3,
		Category: &cat,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(place.PlaceId, "manual_") {
		t.Errorf("expected place_id with 'manual_' prefix, got %q", place.PlaceId)
	}
}
