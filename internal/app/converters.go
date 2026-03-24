package app

import (
	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/store"
)

func listPlacesRowToPlace(r store.ListPlacesRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Category: pstr(r.Category), Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func listPlacesByCategoryRowToPlace(r store.ListPlacesByCategoryRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Category: pstr(r.Category), Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func listPlacesByCategoryAndTagRowToPlace(r store.ListPlacesByCategoryAndTagRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Category: pstr(r.Category), Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func listPlacesByTagRowToPlace(r store.ListPlacesByTagRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Category: pstr(r.Category), Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func getPlaceRowToPlace(r store.GetPlaceRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Category: pstr(r.Category), Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func createPlaceRowToPlace(r store.CreatePlaceRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Category: pstr(r.Category),
		Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}

func updatePlaceRowToPlace(r store.UpdatePlaceRow) generated.Place {
	return generated.Place{
		PlaceId: r.PlaceID, Name: r.Name, Address: pstr(r.Address),
		Lat: r.Lat.Float64, Lng: r.Lng.Float64, Rating: pf64(r.Rating),
		BusinessStatus: pstr(r.BusinessStatus), Category: pstr(r.Category),
		Website: pstr(r.Website),
		FormattedPhoneNumber: pstr(r.FormattedPhoneNumber), InternationalPhoneNumber: pstr(r.InternationalPhoneNumber),
		OpeningHours: pstr(r.OpeningHours), Reviews: pstr(r.Reviews), EditorialSummary: pstr(r.EditorialSummary),
		Photos: pstr(r.Photos), Types: pstr(r.Types), PriceLevel: pint32(r.PriceLevel),
		UserRatingsTotal: pint32(r.UserRatingsTotal), UtcOffsetMinutes: pint32(r.UtcOffsetMinutes),
		GoogleMapsUrl: pstr(r.GoogleMapsUrl), IconUrl: pstr(r.IconUrl),
		CurbsidePickup: pbool(r.CurbsidePickup), Delivery: pbool(r.Delivery), DineIn: pbool(r.DineIn),
		Reservable: pbool(r.Reservable), Takeout: pbool(r.Takeout), WheelchairAccessible: pbool(r.WheelchairAccessible),
		Tags: pstr(r.Tags), ScrapedAt: ptime(r.ScrapedAt),
	}
}
