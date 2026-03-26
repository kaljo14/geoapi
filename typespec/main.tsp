import "@typespec/http";

using TypeSpec.Http;

@service({
  title: "Geopulse",
})
@server("http://localhost:8080", "Local development server")
namespace Geopulse;

// --- Models ---

model Place {
  place_id: string;
  name: string;
  address?: string;
  lat: float64;
  lng: float64;
  rating?: float64;
  business_status?: string;
  website?: string;
  formatted_phone_number?: string;
  international_phone_number?: string;
  opening_hours?: string;
  reviews?: string;
  editorial_summary?: string;
  photos?: string;
  types?: string;
  price_level?: int32;
  user_ratings_total?: int32;
  utc_offset_minutes?: int32;
  google_maps_url?: string;
  icon_url?: string;
  curbside_pickup?: boolean;
  delivery?: boolean;
  dine_in?: boolean;
  reservable?: boolean;
  takeout?: boolean;
  wheelchair_accessible?: boolean;
  category?: string;
  tags?: string;
  scraped_at?: string;
}

model CreatePlaceRequest {
  name: string;
  lat: float64;
  lng: float64;
  address?: string;
  rating?: float64;
  business_status?: string;
  category?: string;
}

model UpdatePlaceRequest {
  name: string;
  lat: float64;
  lng: float64;
  address?: string;
  rating?: float64;
  business_status?: string;
  category?: string;
}

model SaturationResult {
  lat: float64;
  lng: float64;
  radius_meters: float64;
  category?: string;
  competitor_count: int32;
  density_per_km2: float64;
  score: float64;
}

model HeatmapTile {
  lat: float64;
  lng: float64;
  anchor_score: float64;
  comp_penalty: float64;
  score: float64;
}

model GeoJSONFeatureCollection {
  type: string;
  features: unknown[];
}

model LivezResponse {
  status: string;
}

model ReadyzResponse {
  status: string;
  db: string;
}

model ErrorResponse {
  error: string;
}

model MessageResponse {
  message: string;
}

// --- Routes ---

@route("/api/places")
interface Places {
  @get
  @summary("List places")
  @doc("Returns all places, optionally filtered by category and/or tag.")
  list(@query category?: string, @query tag?: string): Place[];

  @post
  @summary("Create place")
  @doc("Creates a new place manually.")
  create(@body body: CreatePlaceRequest): Place | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

@route("/api/places/{place_id}")
interface PlaceByID {
  @get
  @summary("Get place")
  @doc("Returns a single place by ID.")
  get(@path place_id: string): Place | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @put
  @summary("Update place")
  @doc("Updates a place by ID.")
  update(@path place_id: string, @body body: UpdatePlaceRequest): Place | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };

  @delete
  @summary("Delete place")
  @doc("Deletes a place by ID.")
  delete(@path place_id: string): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 404;
    @body body: ErrorResponse;
  };
}

@route("/api/metro/shapes")
interface MetroShapes {
  @get
  @summary("Metro shapes")
  @doc("Returns metro line shapes as GeoJSON FeatureCollection.")
  get(): GeoJSONFeatureCollection;
}

@route("/api/metro/stops")
interface MetroStops {
  @get
  @summary("Metro stops")
  @doc("Returns metro stops as GeoJSON FeatureCollection.")
  get(): GeoJSONFeatureCollection;
}

@route("/api/saturation")
interface Saturation {
  @get
  @summary("Saturation score")
  @doc("Returns competitor density score at a given point.")
  get(
    @query lat: float64,
    @query lng: float64,
    @query radius?: float64,
    @query category?: string,
  ): SaturationResult | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

@route("/api/heatmap")
interface Heatmap {
  @get
  @summary("Opportunity heatmap")
  @doc("Returns opportunity score grid over a bounding box.")
  get(
    @query min_lat: float64,
    @query min_lng: float64,
    @query max_lat: float64,
    @query max_lng: float64,
    @query cell_size?: float64,
    @query category?: string,
  ): HeatmapTile[] | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

@route("/scrape")
interface Scrape {
  @post
  @summary("Trigger scraper")
  @doc("Starts a background Google Places scraper job.")
  create(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };
}

@route("/enrich")
interface Enrich {
  @post
  @summary("Trigger enricher")
  @doc("Starts a background Google Places details enricher job.")
  create(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };
}

@route("/api/osm")
interface OSM {
  @post
  @route("/import")
  @summary("Trigger OSM network import")
  @doc("Starts a background OpenStreetMap pedestrian network import and walkability scoring job for Sofia.")
  trigger(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };

  @post
  @route("/pois")
  @summary("Trigger OSM POI import")
  @doc("Starts a background OpenStreetMap POI import (amenity, shop, tourism, leisure) and walk score recomputation for Sofia.")
  pois(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };
}

@route("/livez")
interface Livez {
  @get
  @summary("Liveness probe")
  @doc("Always returns 200 OK.")
  get(): LivezResponse;
}

@route("/readyz")
interface Readyz {
  @get
  @summary("Readiness probe")
  @doc("Returns 200 if the service is ready. Returns 503 otherwise.")
  get(): ReadyzResponse | {
    @statusCode statusCode: 503;
    @body body: ReadyzResponse;
  };
}
