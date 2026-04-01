package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/generated"
	"github.com/neofyis/geopulse/internal/store"
)

func (a *App) ListNeighborhoods(ctx context.Context) ([]string, error) {
	return a.store.ListNeighborhoodNames(ctx)
}

const sofiaplanBaseURL = "https://api.sofiaplan.bg/datasets"

// sofiaplanDataset maps a dataset ID to its target PostGIS table.
type sofiaplanDataset struct {
	id        int
	tableName string
}

var sofiaplanDatasets = []sofiaplanDataset{
	{292, "sofiaplan_zoning"},
	{502, "sofiaplan_zoning_params"},
	{635, "sofiaplan_income"},
	{636, "sofiaplan_business_turnover"},
	{624, "sofiaplan_property_prices"},
	{39, "sofiaplan_pedestrian_syntax"},
	{437, "sofiaplan_metro_catchments"},
	{640, "sofiaplan_population_grid"},
	{627, "sofiaplan_development_potential"},
	{297, "sofiaplan_neighborhoods"},
	{218, "sofiaplan_census_addresses"},
	{3, "sofiaplan_demographic_forecast"},
	{360, "sofiaplan_demographic_forecast_ge"},
	{622, "sofiaplan_population_potential"},
	{621, "sofiaplan_residential_load"},
	// Accessibility & Transport datasets
	{96, "sofiaplan_transit_access_ge"},
	{279, "sofiaplan_transit_access_district"},
	{289, "sofiaplan_metro_access_800m"},
	{282, "sofiaplan_metro_access_1200m"},
	{268, "sofiaplan_bus_lines"},
	{333, "sofiaplan_bus_lines_alt"},
	{223, "sofiaplan_trolleybus_lines"},
	{472, "sofiaplan_tram_lines"},
	{254, "sofiaplan_tram_lines_alt"},
	{602, "sofiaplan_railway_stations"},
	// Parking zones from ЦГМ
	{470, "sofiaplan_parking_green"}, // Зелена зона за паркиране
	{291, "sofiaplan_parking_blue"},  // Синя зона за паркиране
	// Cycling network
	{606, "sofiaplan_cycling_network"},     // Built cycling network (primary)
	{290, "sofiaplan_cycling_network_alt"}, // Built cycling network (alternate)
	{146, "sofiaplan_cycling_planned"},     // Planned cycling extensions
	// Health services
	{597, "sofiaplan_health_service_concentration"},         // Health service concentration
	{598, "sofiaplan_health_infrastructure_concentration"}, // Health infrastructure concentration by GE
	// Building analysis by GE
	{632, "sofiaplan_building_density_ge"},    // Building density, intensity, enclosure ratio, floor count by GE
	{633, "sofiaplan_building_footprint_ge"},  // Building footprint (ЗП) and total floor area (РЗП) by GE
	{626, "sofiaplan_residential_typology_ge"}, // Residential building typology (single/multi/panel) by GE
	{455, "sofiaplan_urban_morphology_ge"},    // Urban morphology classification by GE
}

// batchSize controls how many features are inserted per SQL statement.
// Kept small enough to avoid memory pressure on the 398MB pedestrian dataset.
const batchSize = 500

// ImportSofiaplan starts a background import. If layer is non-empty only that
// table is imported; otherwise every dataset in sofiaplanDatasets is imported.
func (a *App) ImportSofiaplan(_ context.Context, layer string) error {
	go func() {
		if layer != "" {
			a.logger.Info("sofiaplan single-layer import started", zap.String("layer", layer))
		} else {
			a.logger.Info("sofiaplan full import started")
		}
		if err := a.runSofiaplanImport(context.Background(), layer); err != nil {
			a.logger.Error("sofiaplan import failed", zap.Error(err))
		}
	}()
	return nil
}

func (a *App) runSofiaplanImport(ctx context.Context, layer string) error {
	for _, ds := range sofiaplanDatasets {
		if layer != "" && ds.tableName != layer {
			continue
		}
		if err := a.importDataset(ctx, ds); err != nil {
			a.logger.Error("sofiaplan dataset import failed",
				zap.Int("dataset_id", ds.id),
				zap.String("table", ds.tableName),
				zap.Error(err),
			)
			// continue with remaining datasets
		}
	}
	a.logger.Info("sofiaplan import complete")
	return nil
}

func (a *App) importDataset(ctx context.Context, ds sofiaplanDataset) error {
	url := fmt.Sprintf("%s/%d", sofiaplanBaseURL, ds.id)
	a.logger.Info("downloading sofiaplan dataset", zap.Int("id", ds.id), zap.String("url", url))

	fetchCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(fetchCtx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch dataset %d: %w", ds.id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dataset %d: HTTP %d", ds.id, resp.StatusCode)
	}

	// Skip non-GeoJSON responses (e.g. dataset 502 returns an .xlsx file).
	if ct := resp.Header.Get("Content-Type"); !isGeoJSONContentType(ct) {
		a.logger.Warn("skipping non-GeoJSON dataset",
			zap.Int("dataset_id", ds.id),
			zap.String("content_type", ct),
		)
		return nil
	}

	// Truncate before re-import to avoid duplicate rows.
	if err := a.store.TruncateSofiaplanTable(ctx, ds.tableName); err != nil {
		return fmt.Errorf("truncate %s: %w", ds.tableName, err)
	}

	total, err := a.streamAndInsert(ctx, ds, resp.Body)
	if err != nil {
		return err
	}
	a.logger.Info("sofiaplan dataset imported",
		zap.Int("dataset_id", ds.id),
		zap.String("table", ds.tableName),
		zap.Int("features", total),
	)
	return nil
}

// streamAndInsert parses the GeoJSON response body one feature at a time and
// bulk-inserts in batches of batchSize. Returns the total feature count.
func (a *App) streamAndInsert(ctx context.Context, ds sofiaplanDataset, body io.Reader) (int, error) {
	dec := json.NewDecoder(body)

	// Advance to the "features" array.
	if err := seekToFeatures(dec); err != nil {
		return 0, fmt.Errorf("seek to features: %w", err)
	}

	var (
		bgs2005Detected bool
		bgs2005Checked  bool
		insertSQL       string

		props []string
		geoms []string
		total int
	)

	for dec.More() {
		var feat struct {
			Geometry   json.RawMessage `json:"geometry"`
			Properties json.RawMessage `json:"properties"`
		}
		if err := dec.Decode(&feat); err != nil {
			return total, fmt.Errorf("decode feature: %w", err)
		}
		if feat.Geometry == nil || string(feat.Geometry) == "null" {
			continue
		}

		// On first feature, detect CRS by sniffing coordinate magnitude.
		if !bgs2005Checked {
			bgs2005Detected = isBGS2005(feat.Geometry)
			bgs2005Checked = true
			insertSQL = store.BuildSofiaplanInsertSQL(ds.tableName, bgs2005Detected)
			if bgs2005Detected {
				a.logger.Info("detected BGS2005 CRS, will reproject to WGS84",
					zap.Int("dataset_id", ds.id))
			}
		}

		propsStr := "{}"
		if feat.Properties != nil && string(feat.Properties) != "null" {
			propsStr = string(feat.Properties)
		}
		props = append(props, propsStr)
		geoms = append(geoms, string(feat.Geometry))

		if len(props) >= batchSize {
			if err := a.store.BulkInsertSofiaplanFeatures(ctx, insertSQL, props, geoms); err != nil {
				return total, fmt.Errorf("insert batch: %w", err)
			}
			total += len(props)
			props = props[:0]
			geoms = geoms[:0]
		}
	}

	// flush remaining
	if len(props) > 0 {
		if insertSQL == "" {
			insertSQL = store.BuildSofiaplanInsertSQL(ds.tableName, false)
		}
		if err := a.store.BulkInsertSofiaplanFeatures(ctx, insertSQL, props, geoms); err != nil {
			return total, fmt.Errorf("insert final batch: %w", err)
		}
		total += len(props)
	}

	return total, nil
}

// seekToFeatures advances the decoder to just before the first element of the
// "features" array in a GeoJSON FeatureCollection.
func seekToFeatures(dec *json.Decoder) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		if key, ok := tok.(string); ok && key == "features" {
			// consume the opening '['
			if _, err := dec.Token(); err != nil {
				return err
			}
			return nil
		}
	}
}

// isBGS2005 sniffs the first coordinate pair from a GeoJSON geometry.
// BGS2005 (EPSG:7801) uses eastings ~318000 and northings ~4730000.
// WGS84 uses lng ~23 and lat ~42. Threshold: X > 1000.
func isBGS2005(geomJSON json.RawMessage) bool {
	var geom struct {
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(geomJSON, &geom); err != nil {
		return false
	}
	// Try to decode as a flat [x, y] pair first, then nested arrays.
	x := firstCoordinate(geom.Coordinates)
	return x > 1000
}

// firstCoordinate recursively unwraps nested coordinate arrays until it finds
// a number (the X/easting value of the first coordinate pair).
func firstCoordinate(raw json.RawMessage) float64 {
	if len(raw) == 0 {
		return 0
	}
	// Try as number
	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}
	// Try as array — recurse into first element
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil || len(arr) == 0 {
		return 0
	}
	return firstCoordinate(arr[0])
}

// isGeoJSONContentType returns true for content types that may contain GeoJSON.
// Empty content-type is allowed to avoid blocking ambiguous responses.
func isGeoJSONContentType(ct string) bool {
	if ct == "" {
		return true
	}
	for _, ok := range []string{
		"application/json",
		"application/geo+json",
		"application/octet-stream",
		"text/plain",
	} {
		if strings.HasPrefix(ct, ok) {
			return true
		}
	}
	return false
}

// GetSofiaplanContext queries all datasets in parallel and returns the location context.
func (a *App) GetSofiaplanContext(ctx context.Context, lat, lng float64) (*generated.LocationContext, error) {
	type result struct {
		zoning              json.RawMessage
		income              json.RawMessage
		catchment           json.RawMessage
		inCatch             bool
		pedest              json.RawMessage
		bizTurnover         json.RawMessage
		propPrice           json.RawMessage
		population          json.RawMessage
		devPotential        json.RawMessage
		neighborhood        json.RawMessage
		buildingDensityGe   json.RawMessage
		buildingFootprintGe json.RawMessage
		residTypologyGe     json.RawMessage
		urbanMorphologyGe   json.RawMessage
	}

	var (
		wg   sync.WaitGroup
		res  result
		mu   sync.Mutex
		errs []error
	)

	collect := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	collect(func() error {
		r, err := a.store.GetZoningContext(ctx, lng, lat)
		mu.Lock()
		res.zoning = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetIncomeContext(ctx, lng, lat)
		mu.Lock()
		res.income = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, inC, err := a.store.GetMetroCatchmentContext(ctx, lng, lat)
		mu.Lock()
		res.catchment = r
		res.inCatch = inC
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetPedestrianContext(ctx, lng, lat)
		mu.Lock()
		res.pedest = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetBusinessTurnoverContext(ctx, lng, lat)
		mu.Lock()
		res.bizTurnover = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetPropertyPriceContext(ctx, lng, lat)
		mu.Lock()
		res.propPrice = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetPopulationContext(ctx, lng, lat)
		mu.Lock()
		res.population = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetDevelopmentPotentialContext(ctx, lng, lat)
		mu.Lock()
		res.devPotential = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetNeighborhoodContext(ctx, lng, lat)
		mu.Lock()
		res.neighborhood = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetBuildingDensityGeContext(ctx, lng, lat)
		mu.Lock()
		res.buildingDensityGe = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetBuildingFootprintGeContext(ctx, lng, lat)
		mu.Lock()
		res.buildingFootprintGe = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetResidentialTypologyGeContext(ctx, lng, lat)
		mu.Lock()
		res.residTypologyGe = r
		mu.Unlock()
		return err
	})
	collect(func() error {
		r, err := a.store.GetUrbanMorphologyGeContext(ctx, lng, lat)
		mu.Lock()
		res.urbanMorphologyGe = r
		mu.Unlock()
		return err
	})

	wg.Wait()

	if len(errs) > 0 {
		a.logger.Warn("sofiaplan context partial errors", zap.Errors("errors", errs))
	}

	lc := &generated.LocationContext{
		Lat:              lat,
		Lng:              lng,
		InMetroCatchment: res.inCatch,
	}
	if res.zoning != nil {
		lc.ZoningProperties = rawToMap(res.zoning)
	}
	if res.income != nil {
		lc.IncomeProperties = rawToMap(res.income)
	}
	if res.catchment != nil {
		lc.MetroCatchmentProperties = rawToMap(res.catchment)
	}
	if res.pedest != nil {
		lc.PedestrianProperties = rawToMap(res.pedest)
	}
	if res.bizTurnover != nil {
		lc.BusinessTurnoverProperties = rawToMap(res.bizTurnover)
	}
	if res.propPrice != nil {
		lc.PropertyPriceProperties = rawToMap(res.propPrice)
	}
	if res.population != nil {
		lc.PopulationProperties = rawToMap(res.population)
	}
	if res.devPotential != nil {
		lc.DevelopmentPotentialProperties = rawToMap(res.devPotential)
	}
	if res.neighborhood != nil {
		lc.NeighborhoodProperties = rawToMap(res.neighborhood)
	}
	if res.buildingDensityGe != nil {
		lc.BuildingDensityGeProperties = rawToMap(res.buildingDensityGe)
	}
	if res.buildingFootprintGe != nil {
		lc.BuildingFootprintGeProperties = rawToMap(res.buildingFootprintGe)
	}
	if res.residTypologyGe != nil {
		lc.ResidentialTypologyGeProperties = rawToMap(res.residTypologyGe)
	}
	if res.urbanMorphologyGe != nil {
		lc.UrbanMorphologyGeProperties = rawToMap(res.urbanMorphologyGe)
	}
	return lc, nil
}

// rawToMap converts JSONB bytes to a map for the generated API response.
func rawToMap(raw json.RawMessage) *map[string]interface{} {
	m := make(map[string]interface{})
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	return &m
}
