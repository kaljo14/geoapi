package handler

import (
	"context"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/bff/internal/app"
	"github.com/neofyis/geopulse/bff/internal/generated"
)

type Handler struct {
	app    app.Service
	logger *zap.Logger
}

func NewHandler(a app.Service, logger *zap.Logger) *Handler {
	return &Handler{app: a, logger: logger}
}

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
		return nil, err
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

func (h *Handler) MetroShapesGet(ctx context.Context, _ generated.MetroShapesGetRequestObject) (generated.MetroShapesGetResponseObject, error) {
	fc, err := h.app.GetMetroShapes(ctx)
	if err != nil {
		h.logger.Error("get metro shapes failed", zap.Error(err))
		return nil, err
	}
	return generated.MetroShapesGet200JSONResponse(*fc), nil
}

func (h *Handler) MetroStopsGet(ctx context.Context, _ generated.MetroStopsGetRequestObject) (generated.MetroStopsGetResponseObject, error) {
	fc, err := h.app.GetMetroStops(ctx)
	if err != nil {
		h.logger.Error("get metro stops failed", zap.Error(err))
		return nil, err
	}
	return generated.MetroStopsGet200JSONResponse(*fc), nil
}

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
