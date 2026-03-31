-- Tile-serving views for cycling network datasets.
-- Martin auto-discovers these views and serves them as vector tile sources.

-- Built cycling network (primary): type = path type, posoka = direction, length_m = segment length
CREATE OR REPLACE VIEW sofiaplan_cycling_network_tiles AS
SELECT
    id,
    (properties->>'type')                                   AS label,
    (properties->>'type')                                   AS path_type,
    (properties->>'posoka')                                 AS direction,
    COALESCE((properties->>'length')::numeric, 0)           AS length_m,
    geom
FROM sofiaplan_cycling_network;

-- Built cycling network (alternate): same fields, no length
CREATE OR REPLACE VIEW sofiaplan_cycling_network_alt_tiles AS
SELECT
    id,
    (properties->>'type')                                   AS label,
    (properties->>'type')                                   AS path_type,
    (properties->>'posoka')                                 AS direction,
    geom
FROM sofiaplan_cycling_network_alt;

-- Planned cycling extensions: name = street name, priority/project = planning scores
CREATE OR REPLACE VIEW sofiaplan_cycling_planned_tiles AS
SELECT
    id,
    (properties->>'name')                                   AS label,
    COALESCE((properties->>'priority')::numeric, 0)         AS priority,
    COALESCE((properties->>'project')::numeric, 0)          AS project,
    (properties->>'note')                                   AS note,
    geom
FROM sofiaplan_cycling_planned;
