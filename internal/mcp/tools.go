package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools adds all GeoPulse tools to the MCP server.
func RegisterTools(s *server.MCPServer, client *Client) {
	s.AddTool(searchPlacesTool(), searchPlacesHandler(client))
	s.AddTool(saturationTool(), saturationHandler(client))
	s.AddTool(heatmapTool(), heatmapHandler(client))
	s.AddTool(locationContextTool(), locationContextHandler(client))
	s.AddTool(retailListingsTool(), retailListingsHandler(client))
}

// --- search_places ---

func searchPlacesTool() mcp.Tool {
	return mcp.NewTool("search_places",
		mcp.WithDescription(
			"Search for businesses/places in Sofia by category or tag. "+
				"Търсене на бизнеси/обекти в София по категория или етикет. "+
				"Returns name, address, rating, lat/lng, business_status.",
		),
		mcp.WithString("category",
			mcp.Description("Business category to filter by, e.g. 'barbershop', 'gym', 'cafe'. Категория на бизнеса."),
		),
		mcp.WithString("tag",
			mcp.Description("Tag to filter by. Етикет за филтриране."),
		),
	)
}

func searchPlacesHandler(client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category := req.GetString("category", "")
		tag := req.GetString("tag", "")

		data, err := client.ListPlaces(category, tag)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(data)
	}
}

// --- get_saturation ---

func saturationTool() mcp.Tool {
	return mcp.NewTool("get_saturation",
		mcp.WithDescription(
			"Get competition saturation at a specific location. "+
				"Получаване на насищане от конкуренция на определена локация. "+
				"Returns competitor_count, density_per_km2, and score within the given radius.",
		),
		mcp.WithNumber("lat",
			mcp.Required(),
			mcp.Description("Latitude (географска ширина). Example: 42.6907"),
		),
		mcp.WithNumber("lng",
			mcp.Required(),
			mcp.Description("Longitude (географска дължина). Example: 23.3196"),
		),
		mcp.WithNumber("radius_meters",
			mcp.Description("Search radius in meters (default 500). Радиус на търсене в метри."),
		),
		mcp.WithString("category",
			mcp.Description("Filter competitors by category, e.g. 'barbershop'. Категория за филтриране."),
		),
	)
}

func saturationHandler(client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lat, err := req.RequireFloat("lat")
		if err != nil {
			return mcp.NewToolResultError("lat is required"), nil
		}
		lng, err := req.RequireFloat("lng")
		if err != nil {
			return mcp.NewToolResultError("lng is required"), nil
		}
		radius := req.GetFloat("radius_meters", 500)
		category := req.GetString("category", "")

		data, err := client.GetSaturation(lat, lng, radius, category)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(data)
	}
}

// --- get_opportunity_heatmap ---

func heatmapTool() mcp.Tool {
	return mcp.NewTool("get_opportunity_heatmap",
		mcp.WithDescription(
			"Get a grid of opportunity scores for an area. Each cell has an anchor_score "+
				"(transit proximity), comp_penalty (existing competitors), and composite score. "+
				"Мрежа от оценки за възможности в район. Higher score = better opportunity.",
		),
		mcp.WithNumber("min_lat", mcp.Required(), mcp.Description("South boundary latitude")),
		mcp.WithNumber("min_lng", mcp.Required(), mcp.Description("West boundary longitude")),
		mcp.WithNumber("max_lat", mcp.Required(), mcp.Description("North boundary latitude")),
		mcp.WithNumber("max_lng", mcp.Required(), mcp.Description("East boundary longitude")),
		mcp.WithNumber("cell_size",
			mcp.Description("Grid cell size in degrees (default 0.005 ≈ 500m). Размер на клетката."),
		),
		mcp.WithString("category",
			mcp.Description("Filter competitors by category. Категория за филтриране."),
		),
	)
}

func heatmapHandler(client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		minLat, err := req.RequireFloat("min_lat")
		if err != nil {
			return mcp.NewToolResultError("min_lat is required"), nil
		}
		minLng, err := req.RequireFloat("min_lng")
		if err != nil {
			return mcp.NewToolResultError("min_lng is required"), nil
		}
		maxLat, err := req.RequireFloat("max_lat")
		if err != nil {
			return mcp.NewToolResultError("max_lat is required"), nil
		}
		maxLng, err := req.RequireFloat("max_lng")
		if err != nil {
			return mcp.NewToolResultError("max_lng is required"), nil
		}
		cellSize := req.GetFloat("cell_size", 0.005)
		category := req.GetString("category", "")

		data, err := client.GetHeatmap(minLat, minLng, maxLat, maxLng, cellSize, category)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(data)
	}
}

// --- get_location_context ---

func locationContextTool() mcp.Tool {
	return mcp.NewTool("get_location_context",
		mcp.WithDescription(
			"Get full urban context for a location in Sofia: zoning, metro catchment, "+
				"demographics, flood risk, building density, walkability, and more from 30+ SofiaPlan layers. "+
				"Пълен градски контекст за локация в София.",
		),
		mcp.WithNumber("lat", mcp.Required(), mcp.Description("Latitude (географска ширина)")),
		mcp.WithNumber("lng", mcp.Required(), mcp.Description("Longitude (географска дължина)")),
	)
}

func locationContextHandler(client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lat, err := req.RequireFloat("lat")
		if err != nil {
			return mcp.NewToolResultError("lat is required"), nil
		}
		lng, err := req.RequireFloat("lng")
		if err != nil {
			return mcp.NewToolResultError("lng is required"), nil
		}

		data, err := client.GetLocationContext(lat, lng)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(data)
	}
}

// --- list_retail_listings ---

func retailListingsTool() mcp.Tool {
	return mcp.NewTool("list_retail_listings",
		mcp.WithDescription(
			"List available commercial retail spaces for rent/sale in Sofia. "+
				"Списък на налични търговски помещения в София. "+
				"Returns title, address, size_sqm, price_eur, lat/lng, listing_url.",
		),
	)
}

func retailListingsHandler(client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, err := client.ListRetailListings()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(data)
	}
}

// jsonResult converts raw JSON to a text tool result.
func jsonResult(data json.RawMessage) (*mcp.CallToolResult, error) {
	pretty, err := json.MarshalIndent(json.RawMessage(data), "", "  ")
	if err != nil {
		return mcp.NewToolResultError("failed to format response"), nil
	}
	return mcp.NewToolResultText(string(pretty)), nil
}
