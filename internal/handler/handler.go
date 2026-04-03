package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/app"
	"github.com/neofyis/geopulse/internal/generated"
)

// errInternal is a generic error returned to clients to avoid leaking internal details.
var errInternal = errors.New("internal server error")

type Handler struct {
	app    app.Service
	logger *zap.Logger
}

func NewHandler(a app.Service, logger *zap.Logger) *Handler {
	return &Handler{app: a, logger: logger}
}

// --- Places ---

func (h *Handler) PlacesList(ctx context.Context, req generated.PlacesListRequestObject) (generated.PlacesListResponseObject, error) {
	category := ""
	if req.Params.Category != nil {
		category = *req.Params.Category
	}
	tag := ""
	if req.Params.Tag != nil {
		tag = *req.Params.Tag
	}
	places, err := h.app.ListPlaces(ctx, category, tag)
	if err != nil {
		h.logger.Error("list places failed", zap.Error(err))
		return nil, errInternal
	}
	return generated.PlacesList200JSONResponse(places), nil
}

func (h *Handler) PlacesCreate(ctx context.Context, req generated.PlacesCreateRequestObject) (generated.PlacesCreateResponseObject, error) {
	place, err := h.app.CreatePlace(ctx, *req.Body)
	if err != nil {
		h.logger.Error("create place failed", zap.Error(err))
		return generated.PlacesCreate400JSONResponse{Error: err.Error()}, nil
	}
	return generated.PlacesCreate200JSONResponse(*place), nil
}

func (h *Handler) PlaceByIDGet(ctx context.Context, req generated.PlaceByIDGetRequestObject) (generated.PlaceByIDGetResponseObject, error) {
	place, err := h.app.GetPlace(ctx, req.PlaceId)
	if err != nil {
		h.logger.Warn("get place failed", zap.String("place_id", req.PlaceId), zap.Error(err))
		return generated.PlaceByIDGet404JSONResponse{Error: err.Error()}, nil
	}
	return generated.PlaceByIDGet200JSONResponse(*place), nil
}

func (h *Handler) PlaceByIDUpdate(ctx context.Context, req generated.PlaceByIDUpdateRequestObject) (generated.PlaceByIDUpdateResponseObject, error) {
	place, err := h.app.UpdatePlace(ctx, req.PlaceId, *req.Body)
	if err != nil {
		h.logger.Warn("update place failed", zap.String("place_id", req.PlaceId), zap.Error(err))
		return generated.PlaceByIDUpdate404JSONResponse{Error: err.Error()}, nil
	}
	return generated.PlaceByIDUpdate200JSONResponse(*place), nil
}

func (h *Handler) PlaceByIDDelete(ctx context.Context, req generated.PlaceByIDDeleteRequestObject) (generated.PlaceByIDDeleteResponseObject, error) {
	if err := h.app.DeletePlace(ctx, req.PlaceId); err != nil {
		h.logger.Warn("delete place failed", zap.String("place_id", req.PlaceId), zap.Error(err))
		return generated.PlaceByIDDelete404JSONResponse{Error: err.Error()}, nil
	}
	return generated.PlaceByIDDelete204Response{}, nil
}

// --- Sofiaplan ---

func (h *Handler) GetNeighborhoods(w http.ResponseWriter, r *http.Request) {
	names, err := h.app.ListNeighborhoods(r.Context())
	if err != nil {
		h.logger.Error("list neighborhoods failed", zap.Error(err))
		writeJSONError(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(names)
}

// --- Metro ---

func (h *Handler) MetroShapesGet(ctx context.Context, _ generated.MetroShapesGetRequestObject) (generated.MetroShapesGetResponseObject, error) {
	fc, err := h.app.GetMetroShapes(ctx)
	if err != nil {
		h.logger.Error("get metro shapes failed", zap.Error(err))
		return nil, errInternal
	}
	return generated.MetroShapesGet200JSONResponse(*fc), nil
}

func (h *Handler) MetroStopsGet(ctx context.Context, _ generated.MetroStopsGetRequestObject) (generated.MetroStopsGetResponseObject, error) {
	fc, err := h.app.GetMetroStops(ctx)
	if err != nil {
		h.logger.Error("get metro stops failed", zap.Error(err))
		return nil, errInternal
	}
	return generated.MetroStopsGet200JSONResponse(*fc), nil
}

func (h *Handler) GetTransitStops(w http.ResponseWriter, r *http.Request) {
	fc, err := h.app.GetTransitStops(r.Context())
	if err != nil {
		h.logger.Error("get transit stops failed", zap.Error(err))
		writeJSONError(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(fc)
}

// --- Parking Zones ---

func (h *Handler) ParkingZonesGet(ctx context.Context, _ generated.ParkingZonesGetRequestObject) (generated.ParkingZonesGetResponseObject, error) {
	fc, err := h.app.GetParkingZones(ctx)
	if err != nil {
		h.logger.Error("get parking zones failed", zap.Error(err))
		return nil, errInternal
	}
	return generated.ParkingZonesGet200JSONResponse(*fc), nil
}

// --- Retail Listings ---

func (h *Handler) RetailListingsList(ctx context.Context, _ generated.RetailListingsListRequestObject) (generated.RetailListingsListResponseObject, error) {
	listings, err := h.app.ListRetailListings(ctx)
	if err != nil {
		h.logger.Error("list retail listings failed", zap.Error(err))
		return nil, errInternal
	}
	return generated.RetailListingsList200JSONResponse(listings), nil
}

func (h *Handler) RetailListingsCreate(ctx context.Context, req generated.RetailListingsCreateRequestObject) (generated.RetailListingsCreateResponseObject, error) {
	listing, err := h.app.CreateRetailListing(ctx, *req.Body)
	if err != nil {
		h.logger.Error("create retail listing failed", zap.Error(err))
		return generated.RetailListingsCreate400JSONResponse{Error: err.Error()}, nil
	}
	return generated.RetailListingsCreate200JSONResponse(*listing), nil
}

func (h *Handler) RetailListingByIDGet(ctx context.Context, req generated.RetailListingByIDGetRequestObject) (generated.RetailListingByIDGetResponseObject, error) {
	listing, err := h.app.GetRetailListing(ctx, req.Id)
	if err != nil {
		h.logger.Warn("get retail listing failed", zap.String("id", req.Id), zap.Error(err))
		return generated.RetailListingByIDGet404JSONResponse{Error: err.Error()}, nil
	}
	return generated.RetailListingByIDGet200JSONResponse(*listing), nil
}

func (h *Handler) RetailListingByIDUpdate(ctx context.Context, req generated.RetailListingByIDUpdateRequestObject) (generated.RetailListingByIDUpdateResponseObject, error) {
	listing, err := h.app.UpdateRetailListing(ctx, req.Id, *req.Body)
	if err != nil {
		h.logger.Warn("update retail listing failed", zap.String("id", req.Id), zap.Error(err))
		return generated.RetailListingByIDUpdate404JSONResponse{Error: err.Error()}, nil
	}
	return generated.RetailListingByIDUpdate200JSONResponse(*listing), nil
}

func (h *Handler) RetailListingByIDDelete(ctx context.Context, req generated.RetailListingByIDDeleteRequestObject) (generated.RetailListingByIDDeleteResponseObject, error) {
	if err := h.app.DeleteRetailListing(ctx, req.Id); err != nil {
		h.logger.Warn("delete retail listing failed", zap.String("id", req.Id), zap.Error(err))
		return generated.RetailListingByIDDelete404JSONResponse{Error: err.Error()}, nil
	}
	return generated.RetailListingByIDDelete204Response{}, nil
}

// --- Analytics ---

func (h *Handler) SaturationGet(ctx context.Context, req generated.SaturationGetRequestObject) (generated.SaturationGetResponseObject, error) {
	radius := 500.0
	if req.Params.Radius != nil {
		radius = *req.Params.Radius
	}
	category := ""
	if req.Params.Category != nil {
		category = *req.Params.Category
	}
	result, err := h.app.GetSaturation(ctx, req.Params.Lat, req.Params.Lng, radius, category)
	if err != nil {
		h.logger.Error("get saturation failed", zap.Error(err))
		return generated.SaturationGet400JSONResponse{Error: err.Error()}, nil
	}
	return generated.SaturationGet200JSONResponse(*result), nil
}

func (h *Handler) HeatmapGet(ctx context.Context, req generated.HeatmapGetRequestObject) (generated.HeatmapGetResponseObject, error) {
	cellSize := 0.005
	if req.Params.CellSize != nil {
		cellSize = *req.Params.CellSize
	}
	category := ""
	if req.Params.Category != nil {
		category = *req.Params.Category
	}
	tiles, err := h.app.GetHeatmap(ctx, req.Params.MinLat, req.Params.MinLng, req.Params.MaxLat, req.Params.MaxLng, cellSize, category)
	if err != nil {
		h.logger.Error("get heatmap failed", zap.Error(err))
		return generated.HeatmapGet400JSONResponse{Error: err.Error()}, nil
	}
	return generated.HeatmapGet200JSONResponse(tiles), nil
}

// --- Scrape / Enrich ---

func (h *Handler) ScrapeCreate(ctx context.Context, _ generated.ScrapeCreateRequestObject) (generated.ScrapeCreateResponseObject, error) {
	if err := h.app.StartScraper(ctx); err != nil {
		h.logger.Warn("start scraper failed", zap.Error(err))
		return generated.ScrapeCreate202JSONResponse{Message: err.Error()}, nil
	}
	return generated.ScrapeCreate202JSONResponse{Message: "scraper started"}, nil
}

func (h *Handler) EnrichCreate(ctx context.Context, _ generated.EnrichCreateRequestObject) (generated.EnrichCreateResponseObject, error) {
	if err := h.app.StartEnricher(ctx); err != nil {
		h.logger.Warn("start enricher failed", zap.Error(err))
		return generated.EnrichCreate202JSONResponse{Message: err.Error()}, nil
	}
	return generated.EnrichCreate202JSONResponse{Message: "enricher started"}, nil
}

func (h *Handler) OSMTrigger(ctx context.Context, _ generated.OSMTriggerRequestObject) (generated.OSMTriggerResponseObject, error) {
	if err := h.app.ImportOSMNetwork(ctx); err != nil {
		h.logger.Warn("osm import start failed", zap.Error(err))
		return generated.OSMTrigger202JSONResponse{Message: err.Error()}, nil
	}
	return generated.OSMTrigger202JSONResponse{Message: "osm import started"}, nil
}

func (h *Handler) OSMPois(ctx context.Context, _ generated.OSMPoisRequestObject) (generated.OSMPoisResponseObject, error) {
	if err := h.app.ImportOSMPOIs(ctx); err != nil {
		h.logger.Warn("osm poi import start failed", zap.Error(err))
		return generated.OSMPois202JSONResponse{Message: err.Error()}, nil
	}
	return generated.OSMPois202JSONResponse{Message: "osm poi import started"}, nil
}

// --- SofiaPlan ---

func (h *Handler) SofiaPlanTrigger(ctx context.Context, req generated.SofiaPlanTriggerRequestObject) (generated.SofiaPlanTriggerResponseObject, error) {
	layer := ""
	if req.Params.Layer != nil {
		layer = *req.Params.Layer
	}
	if err := h.app.ImportSofiaplan(ctx, layer); err != nil {
		h.logger.Warn("sofiaplan import start failed", zap.Error(err))
		return generated.SofiaPlanTrigger202JSONResponse{Message: err.Error()}, nil
	}
	msg := "sofiaplan import started"
	if layer != "" {
		msg = "sofiaplan import started for layer: " + layer
	}
	return generated.SofiaPlanTrigger202JSONResponse{Message: msg}, nil
}

func (h *Handler) SofiaPlanContext(ctx context.Context, req generated.SofiaPlanContextRequestObject) (generated.SofiaPlanContextResponseObject, error) {
	result, err := h.app.GetSofiaplanContext(ctx, req.Params.Lat, req.Params.Lng)
	if err != nil {
		h.logger.Error("sofiaplan context failed", zap.Error(err))
		return generated.SofiaPlanContext400JSONResponse{Error: err.Error()}, nil
	}
	return generated.SofiaPlanContext200JSONResponse(*result), nil
}

// --- Health ---

func (h *Handler) LivezGet(_ context.Context, _ generated.LivezGetRequestObject) (generated.LivezGetResponseObject, error) {
	return generated.LivezGet200JSONResponse{Status: "ok"}, nil
}

func (h *Handler) ReadyzGet(ctx context.Context, _ generated.ReadyzGetRequestObject) (generated.ReadyzGetResponseObject, error) {
	if err := h.app.Ready(ctx); err != nil {
		h.logger.Warn("readyz check failed", zap.Error(err))
		return generated.ReadyzGet503JSONResponse{Status: "unavailable", Db: err.Error()}, nil
	}
	return generated.ReadyzGet200JSONResponse{Status: "ok", Db: "ok"}, nil
}

// --- Manual CSV routes (bypass oapi-codegen) ---

func (h *Handler) ExportPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="places.csv"`)
	if err := h.app.ExportCSV(r.Context(), w, false); err != nil {
		h.logger.Error("export places failed", zap.Error(err))
	}
}

func (h *Handler) ExportPlacesSimple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="places-simple.csv"`)
	if err := h.app.ExportCSV(r.Context(), w, true); err != nil {
		h.logger.Error("export places simple failed", zap.Error(err))
	}
}

func (h *Handler) ImportPlaces(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMultipartSize); err != nil {
		writeJSONError(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, "file field required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	count, err := h.app.ImportCSV(r.Context(), file)
	if err != nil {
		h.logger.Error("import places failed", zap.Error(err))
		writeJSONError(w, "import failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"imported": count})
}

// maxMultipartSize is the maximum size for multipart form uploads (10 MB).
const maxMultipartSize = 10 << 20

// writeJSONError writes a JSON error response matching the OpenAPI error format.
func writeJSONError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
