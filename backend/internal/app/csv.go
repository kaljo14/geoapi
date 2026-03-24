package app

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"github.com/neofyis/geopulse/backend/internal/store"
)

var fullExportHeaders = []string{
	"place_id", "name", "address", "lat", "lng", "rating", "business_status",
	"website", "formatted_phone_number", "international_phone_number",
	"opening_hours", "reviews", "editorial_summary", "photos", "types",
	"price_level", "user_ratings_total", "utc_offset_minutes",
	"google_maps_url", "icon_url", "curbside_pickup", "delivery", "dine_in",
	"reservable", "takeout", "wheelchair_accessible", "category", "tags", "scraped_at",
}

var simpleExportHeaders = []string{"name", "place_id", "website"}

func (a *App) ExportCSV(ctx context.Context, w io.Writer, simple bool) error {
	cw := csv.NewWriter(w)

	// Flush the header row first so the HTTP Content-Type header is committed
	// before any database call that could fail.
	if simple {
		_ = cw.Write(simpleExportHeaders)
	} else {
		_ = cw.Write(fullExportHeaders)
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	if a.store == nil {
		return fmt.Errorf("database not configured")
	}

	rows, err := a.store.ListPlacesForExport(ctx)
	if err != nil {
		return fmt.Errorf("list places for export: %w", err)
	}

	if simple {
		for _, r := range rows {
			_ = cw.Write([]string{r.Name, r.PlaceID, r.Website.String})
		}
	} else {
		for _, r := range rows {
			_ = cw.Write([]string{
				r.PlaceID, r.Name, r.Address.String,
				f64str(r.Lat.Float64), f64str(r.Lng.Float64), f64str(r.Rating.Float64),
				r.BusinessStatus.String, r.Website.String,
				r.FormattedPhoneNumber.String, r.InternationalPhoneNumber.String,
				r.OpeningHours.String, r.Reviews.String, r.EditorialSummary.String,
				r.Photos.String, r.Types.String,
				i32str(r.PriceLevel.Int32), i32str(r.UserRatingsTotal.Int32), i32str(r.UtcOffsetMinutes.Int32),
				r.GoogleMapsUrl.String, r.IconUrl.String,
				boolstr(r.CurbsidePickup.Bool), boolstr(r.Delivery.Bool), boolstr(r.DineIn.Bool),
				boolstr(r.Reservable.Bool), boolstr(r.Takeout.Bool), boolstr(r.WheelchairAccessible.Bool),
				r.Category.String, r.Tags.String,
				timestr(r.ScrapedAt),
			})
		}
	}
	cw.Flush()
	return cw.Error()
}

func (a *App) ImportCSV(ctx context.Context, r io.Reader) (int, error) {
	cr := csv.NewReader(r)
	headers, err := cr.Read()
	if err != nil {
		return 0, fmt.Errorf("read csv headers: %w", err)
	}

	records, err := cr.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("read csv records: %w", err)
	}

	imported := 0
	switch len(headers) {
	case 3:
		// simple: name, place_id, website
		for _, rec := range records {
			if len(rec) < 3 {
				continue
			}
			_, err := a.store.UpsertPlace(ctx, store.UpsertPlaceParams{
				PlaceID: rec[1],
				Name:    rec[0],
				Website: pgtype.Text{String: rec[2], Valid: rec[2] != ""},
			})
			if err != nil {
				a.logger.Warn("import: upsert failed", zap.String("place_id", rec[1]), zap.Error(err))
				continue
			}
			imported++
		}
	case 8:
		// basic: place_id, name, address, lat, lng, rating, business_status, category
		for _, rec := range records {
			if len(rec) < 8 {
				continue
			}
			lat, _ := strconv.ParseFloat(rec[3], 64)
			lng, _ := strconv.ParseFloat(rec[4], 64)
			rating, _ := strconv.ParseFloat(rec[5], 64)
			_, err := a.store.UpsertPlace(ctx, store.UpsertPlaceParams{
				PlaceID:        rec[0],
				Name:           rec[1],
				Address:        pgtype.Text{String: rec[2], Valid: rec[2] != ""},
				Lat:            pgtype.Float8{Float64: lat, Valid: lat != 0},
				Lng:            pgtype.Float8{Float64: lng, Valid: lng != 0},
				Rating:         pgtype.Float8{Float64: rating, Valid: rating != 0},
				BusinessStatus: pgtype.Text{String: rec[6], Valid: rec[6] != ""},
				Category:       pgtype.Text{String: rec[7], Valid: rec[7] != ""},
			})
			if err != nil {
				a.logger.Warn("import: upsert failed", zap.String("place_id", rec[0]), zap.Error(err))
				continue
			}
			imported++
		}
	default:
		// full: all 29 cols matching fullExportHeaders order
		for _, rec := range records {
			if len(rec) < len(fullExportHeaders) {
				continue
			}
			lat, _ := strconv.ParseFloat(rec[3], 64)
			lng, _ := strconv.ParseFloat(rec[4], 64)
			rating, _ := strconv.ParseFloat(rec[5], 64)
			priceLevel, _ := strconv.ParseInt(rec[15], 10, 32)
			userRatings, _ := strconv.ParseInt(rec[16], 10, 32)
			utcOffset, _ := strconv.ParseInt(rec[17], 10, 32)
			_, err := a.store.UpsertPlace(ctx, store.UpsertPlaceParams{
				PlaceID:                  rec[0],
				Name:                     rec[1],
				Address:                  pgtype.Text{String: rec[2], Valid: rec[2] != ""},
				Lat:                      pgtype.Float8{Float64: lat, Valid: lat != 0},
				Lng:                      pgtype.Float8{Float64: lng, Valid: lng != 0},
				Rating:                   pgtype.Float8{Float64: rating, Valid: rating != 0},
				BusinessStatus:           pgtype.Text{String: rec[6], Valid: rec[6] != ""},
				Website:                  pgtype.Text{String: rec[7], Valid: rec[7] != ""},
				FormattedPhoneNumber:     pgtype.Text{String: rec[8], Valid: rec[8] != ""},
				InternationalPhoneNumber: pgtype.Text{String: rec[9], Valid: rec[9] != ""},
				OpeningHours:             pgtype.Text{String: rec[10], Valid: rec[10] != ""},
				Reviews:                  pgtype.Text{String: rec[11], Valid: rec[11] != ""},
				EditorialSummary:         pgtype.Text{String: rec[12], Valid: rec[12] != ""},
				Photos:                   pgtype.Text{String: rec[13], Valid: rec[13] != ""},
				Types:                    pgtype.Text{String: rec[14], Valid: rec[14] != ""},
				PriceLevel:               pgtype.Int4{Int32: int32(priceLevel), Valid: priceLevel != 0},
				UserRatingsTotal:         pgtype.Int4{Int32: int32(userRatings), Valid: userRatings != 0},
				UtcOffsetMinutes:         pgtype.Int4{Int32: int32(utcOffset), Valid: utcOffset != 0},
				GoogleMapsUrl:            pgtype.Text{String: rec[18], Valid: rec[18] != ""},
				IconUrl:                  pgtype.Text{String: rec[19], Valid: rec[19] != ""},
				Category:                 pgtype.Text{String: rec[26], Valid: rec[26] != ""},
				Tags:                     pgtype.Text{String: rec[27], Valid: rec[27] != ""},
			})
			if err != nil {
				a.logger.Warn("import: upsert failed", zap.String("place_id", rec[0]), zap.Error(err))
				continue
			}
			imported++
		}
	}
	return imported, nil
}
