package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/backend/internal/app"
	"github.com/neofyis/geopulse/backend/internal/generated"
)

type Handler struct {
	app    app.Service
	logger *zap.Logger
}

func NewHandler(a app.Service, logger *zap.Logger) *Handler {
	return &Handler{app: a, logger: logger}
}

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

// ExportPlaces streams all places as CSV (full columns).
func (h *Handler) ExportPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="places.csv"`)
	if err := h.app.ExportCSV(r.Context(), w, false); err != nil {
		h.logger.Error("export places failed", zap.Error(err))
	}
}

// ExportPlacesSimple streams a simple 3-column CSV (name, place_id, website).
func (h *Handler) ExportPlacesSimple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="places-simple.csv"`)
	if err := h.app.ExportCSV(r.Context(), w, true); err != nil {
		h.logger.Error("export places simple failed", zap.Error(err))
	}
}

// ImportPlaces accepts a multipart CSV upload and bulk-inserts places.
func (h *Handler) ImportPlaces(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	count, err := h.app.ImportCSV(r.Context(), file)
	if err != nil {
		h.logger.Error("import places failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"imported": count})
}
