-- Ultimate predictive foot traffic: combines space syntax integration,
-- residential population density, POI attraction, and transit proximity.
--
-- Extends migration 000041 (calibrated_foot_traffic_tiles) which uses
-- only space syntax. This view adds three additional signals and produces
-- a separate materialized view for side-by-side comparison.
--
-- Signal weights (geometric mean exponents):
--   syntax 0.50, population 0.25, POI 0.15, transit 0.10
--
-- Formula: predicted_hourly = 80 * (syntax^0.50 * pop^0.25 * poi^0.15 * transit^0.10) ^ 1.5
--
-- All spatial operations use GEOMETRY with GIST indexes for speed.
-- Distance decay uses degree-based radii (at Sofia 42.7°N):
--   0.0024° ≈ 200m (avg of lat/lon at this latitude)
--   0.0018° ≈ 150m
--   0.0036° ≈ 300m
-- Since all segments are scored with the same distance function, the
-- PERCENT_RANK ordering is identical to meter-based ranking.
-- The final predicted_hourly is calibrated by the composite formula,
-- not by individual distance values.

CREATE MATERIALIZED VIEW ultimate_foot_traffic_tiles AS
WITH syntax_ranked AS (
    SELECT
        t.id,
        t.neighborhood,
        t.score         AS integration,
        t.choice,
        t.segment_length,
        PERCENT_RANK() OVER (ORDER BY t.score) AS syntax_pctl,
        ST_Centroid(t.geom) AS centroid,
        t.geom
    FROM sofiaplan_pedestrian_syntax_tiles t
),
-- Population: sum residential load within ~200m, linear distance decay
pop_raw AS (
    SELECT
        s.id,
        COALESCE(SUM(
            COALESCE((r.properties->>'ppl_30kvm')::numeric, 0)
            * GREATEST(0, 1.0 - ST_Distance(s.centroid, ST_Centroid(r.geom)) / 0.0024)
        ), 0) AS pop_score
    FROM syntax_ranked s
    LEFT JOIN sofiaplan_residential_load r
        ON r.geom IS NOT NULL
        AND ST_DWithin(s.centroid, r.geom, 0.0024)
    GROUP BY s.id
),
pop_ranked AS (
    SELECT id,
           GREATEST(0.01, PERCENT_RANK() OVER (ORDER BY pop_score)) AS pop_pctl
    FROM pop_raw
),
-- OSM POI attraction: weighted count within ~150m
osm_poi_scores AS (
    SELECT
        s.id,
        COALESCE(SUM(
            CASE
                WHEN o.category IN ('supermarket','marketplace','department_store') THEN 5.0
                WHEN o.category IN ('restaurant','cafe','fast_food','bar','pub',
                    'university','school','hospital','cinema','theatre') THEN 3.0
                WHEN o.category IN ('bank','pharmacy','clothes','convenience',
                    'bakery','hairdresser','beauty','gym','fitness_centre') THEN 2.0
                ELSE 1.0
            END
            * GREATEST(0, 1.0 - ST_Distance(s.centroid, o.geom) / 0.0018)
        ), 0) AS osm_score
    FROM syntax_ranked s
    LEFT JOIN osm_pois o
        ON o.geom IS NOT NULL
        AND ST_DWithin(s.centroid, o.geom, 0.0018)
    GROUP BY s.id
),
-- Google Places attraction: rating-weighted within ~150m
places_scores AS (
    SELECT
        s.id,
        COALESCE(SUM(
            LEAST(COALESCE(p.user_ratings_total, 100)::numeric / 500.0, 3.0)
            * GREATEST(0, 1.0 - ST_Distance(s.centroid, p.location) / 0.0018)
        ), 0) AS places_score
    FROM syntax_ranked s
    LEFT JOIN places p
        ON p.location IS NOT NULL
        AND p.business_status = 'OPERATIONAL'
        AND ST_DWithin(s.centroid, p.location, 0.0018)
    GROUP BY s.id
),
poi_ranked AS (
    SELECT o.id,
           GREATEST(0.01, PERCENT_RANK() OVER (
               ORDER BY o.osm_score + p.places_score
           )) AS poi_pctl
    FROM osm_poi_scores o
    JOIN places_scores p ON p.id = o.id
),
-- Transit: count GTFS stops within ~300m with distance decay
transit_raw AS (
    SELECT
        s.id,
        COALESCE(SUM(
            GREATEST(0, 1.0 - ST_Distance(s.centroid, g.geom) / 0.0036)
        ), 0) AS transit_score
    FROM syntax_ranked s
    LEFT JOIN gtfs_stops g
        ON g.geom IS NOT NULL
        AND ST_DWithin(s.centroid, g.geom, 0.0036)
    GROUP BY s.id
),
transit_ranked AS (
    SELECT id,
           GREATEST(0.01, PERCENT_RANK() OVER (ORDER BY transit_score)) AS transit_pctl
    FROM transit_raw
)
SELECT
    s.id,
    s.neighborhood,
    s.integration,
    s.choice,
    s.segment_length,
    s.syntax_pctl,
    p.pop_pctl,
    o.poi_pctl,
    t.transit_pctl,
    LEAST(200, GREATEST(0, ROUND(
        80.0 * POWER(
            POWER(s.syntax_pctl, 0.50)
            * POWER(p.pop_pctl, 0.25)
            * POWER(o.poi_pctl, 0.15)
            * POWER(t.transit_pctl, 0.10),
            1.5
        )
    )))::int AS predicted_hourly,
    s.geom
FROM syntax_ranked s
JOIN pop_ranked p ON p.id = s.id
JOIN poi_ranked o ON o.id = s.id
JOIN transit_ranked t ON t.id = s.id;

CREATE INDEX idx_ultimate_ft_geom
    ON ultimate_foot_traffic_tiles USING GIST(geom);
