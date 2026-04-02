package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/neofyis/geopulse/internal/store"
)

type overpassResponse struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	Type  string            `json:"type"`
	ID    int64             `json:"id"`
	Lat   float64           `json:"lat"`
	Lon   float64           `json:"lon"`
	Nodes []int64           `json:"nodes"`
	Tags  map[string]string `json:"tags"`
}

const overpassURL = "https://overpass-api.de/api/interpreter"

// Sofia city centre bounding box: (south,west,north,east)
const overpassQuery = `[out:json][timeout:90][maxsize:134217728];
(
  way["highway"~"^(footway|path|pedestrian|steps|residential|living_street|service)$"]
  (42.67,23.28,42.73,23.38);
);
out body;
>;
out skel qt;`

const overpassPOIQuery = `[out:json][timeout:90][maxsize:134217728];
(
  node["amenity"](42.67,23.28,42.73,23.38);
  node["shop"](42.67,23.28,42.73,23.38);
  node["tourism"](42.67,23.28,42.73,23.38);
  node["leisure"](42.67,23.28,42.73,23.38);
);
out qt;`

func derivePOICategory(amenity, shop, tourism, leisure string) string {
	switch amenity {
	case "restaurant", "cafe", "bar", "fast_food", "food_court", "pub":
		return "food"
	case "school", "university", "kindergarten", "library":
		return "education"
	case "hospital", "clinic", "pharmacy", "dentist", "doctors":
		return "health"
	case "bank", "atm", "post_office":
		return "finance"
	}
	if shop != "" {
		return "retail"
	}
	if tourism != "" {
		return "culture"
	}
	if leisure != "" {
		return "leisure"
	}
	return "other"
}

func (a *App) ImportOSMNetwork(_ context.Context) error {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.logger.Info("osm import started")
		if err := a.runOSMImport(a.ctx); err != nil {
			a.logger.Error("osm import failed", zap.Error(err))
		}
	}()
	return nil
}

func (a *App) ImportOSMPOIs(_ context.Context) error {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.logger.Info("osm poi import started")
		if err := a.runOSMPOIImport(a.ctx); err != nil {
			a.logger.Error("osm poi import failed", zap.Error(err))
		}
	}()
	return nil
}

func (a *App) runOSMImport(ctx context.Context) error {
	fetchCtx, cancel := context.WithTimeout(ctx, 95*time.Second)
	defer cancel()

	resp, err := a.fetchOverpass(fetchCtx, overpassQuery)
	if err != nil {
		return fmt.Errorf("fetch overpass: %w", err)
	}

	nodeIndex, nodes, ways := a.parseOSMResponse(resp)

	if err := a.store.BulkUpsertNodes(ctx, nodes); err != nil {
		return fmt.Errorf("upsert nodes: %w", err)
	}
	a.logger.Info("osm nodes upserted", zap.Int("count", len(nodes)))

	edges := a.buildEdges(ways, nodeIndex)
	if err := a.store.BulkUpsertEdges(ctx, edges); err != nil {
		return fmt.Errorf("upsert edges: %w", err)
	}
	a.logger.Info("osm edges upserted", zap.Int("count", len(edges)))

	scores, err := a.store.ComputeWalkScores(ctx)
	if err != nil {
		return fmt.Errorf("compute walk scores: %w", err)
	}

	if err := a.store.UpdateNodeWalkScores(ctx, scores); err != nil {
		return fmt.Errorf("update walk scores: %w", err)
	}

	if err := a.store.PropagateEdgeScores(ctx); err != nil {
		return fmt.Errorf("propagate edge scores: %w", err)
	}

	a.logger.Info("osm import complete",
		zap.Int("nodes", len(nodes)),
		zap.Int("edges", len(edges)),
		zap.Int("scored", len(scores)),
	)
	return nil
}

func (a *App) runOSMPOIImport(ctx context.Context) error {
	fetchCtx, cancel := context.WithTimeout(ctx, 95*time.Second)
	defer cancel()

	resp, err := a.fetchOverpass(fetchCtx, overpassPOIQuery)
	if err != nil {
		return fmt.Errorf("fetch overpass pois: %w", err)
	}

	var pois []store.OSMPOI
	for _, el := range resp.Elements {
		if el.Type != "node" {
			continue
		}
		amenity := el.Tags["amenity"]
		shop := el.Tags["shop"]
		tourism := el.Tags["tourism"]
		leisure := el.Tags["leisure"]
		pois = append(pois, store.OSMPOI{
			OsmID:    el.ID,
			Lat:      el.Lat,
			Lng:      el.Lon,
			Name:     el.Tags["name"],
			Amenity:  amenity,
			Shop:     shop,
			Tourism:  tourism,
			Leisure:  leisure,
			Category: derivePOICategory(amenity, shop, tourism, leisure),
		})
	}

	if err := a.store.BulkUpsertPOIs(ctx, pois); err != nil {
		return fmt.Errorf("upsert pois: %w", err)
	}
	a.logger.Info("osm pois upserted", zap.Int("count", len(pois)))

	scores, err := a.store.ComputeWalkScores(ctx)
	if err != nil {
		return fmt.Errorf("compute walk scores: %w", err)
	}
	if err := a.store.UpdateNodeWalkScores(ctx, scores); err != nil {
		return fmt.Errorf("update walk scores: %w", err)
	}
	if err := a.store.PropagateEdgeScores(ctx); err != nil {
		return fmt.Errorf("propagate edge scores: %w", err)
	}

	a.logger.Info("osm poi import complete", zap.Int("pois", len(pois)), zap.Int("scored", len(scores)))
	return nil
}

func (a *App) fetchOverpass(ctx context.Context, query string) (*overpassResponse, error) {
	body := strings.NewReader("data=" + url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, overpassURL, body)
	if err != nil {
		return nil, fmt.Errorf("build overpass request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("overpass request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overpass status %d", resp.StatusCode)
	}

	var parsed overpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode overpass response: %w", err)
	}
	return &parsed, nil
}

func (a *App) parseOSMResponse(resp *overpassResponse) (
	nodeIndex map[int64]overpassElement,
	nodes []store.OSMNode,
	ways []overpassElement,
) {
	nodeIndex = make(map[int64]overpassElement, len(resp.Elements))
	for _, el := range resp.Elements {
		switch el.Type {
		case "node":
			nodeIndex[el.ID] = el
			nodes = append(nodes, store.OSMNode{
				OsmID: el.ID,
				Lat:   el.Lat,
				Lng:   el.Lon,
			})
		case "way":
			ways = append(ways, el)
		}
	}
	return nodeIndex, nodes, ways
}

func (a *App) buildEdges(ways []overpassElement, nodeIndex map[int64]overpassElement) []store.OSMEdge {
	var edges []store.OSMEdge
	for _, way := range ways {
		highway := way.Tags["highway"]
		refs := way.Nodes
		for i := 0; i < len(refs)-1; i++ {
			fromNode, okFrom := nodeIndex[refs[i]]
			toNode, okTo := nodeIndex[refs[i+1]]
			if !okFrom || !okTo {
				continue
			}
			edges = append(edges, store.OSMEdge{
				OsmWayID:   way.ID,
				FromNodeID: refs[i],
				ToNodeID:   refs[i+1],
				Highway:    highway,
				FromLat:    fromNode.Lat,
				FromLng:    fromNode.Lon,
				ToLat:      toNode.Lat,
				ToLng:      toNode.Lon,
			})
		}
	}
	return edges
}
