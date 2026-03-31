-- Tile-serving views for health service datasets.
-- Martin auto-discovers these views and serves them as vector tile sources.

-- Health service concentration: numpoints = number of health service points, regname = locality name
CREATE OR REPLACE VIEW sofiaplan_health_service_concentration_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon', '')    AS label,
    (properties->>'rajon')                                        AS district,
    COALESCE((properties->>'numpoints')::numeric, 0)              AS score,
    geom
FROM sofiaplan_health_service_concentration;

-- Health infrastructure concentration by GE: numpoints = number of health infrastructure points, regname = locality name
CREATE OR REPLACE VIEW sofiaplan_health_infrastructure_concentration_tiles AS
SELECT
    id,
    COALESCE(properties->>'regname', properties->>'rajon', '')    AS label,
    (properties->>'rajon')                                        AS district,
    COALESCE((properties->>'numpoints')::numeric, 0)              AS score,
    geom
FROM sofiaplan_health_infrastructure_concentration;
