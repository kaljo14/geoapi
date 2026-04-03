-- Update flood risk tile views to expose zone_id and score for frontend
-- Must DROP first because CREATE OR REPLACE cannot rename columns
DROP VIEW IF EXISTS sofiaplan_flood_risk_low_tiles;
CREATE VIEW sofiaplan_flood_risk_low_tiles AS
SELECT
    id,
    COALESCE(properties->>'apsfr', '')   AS label,
    COALESCE(properties->>'apsfr', '')   AS zone_id,
    1                                    AS score,
    1                                    AS risk_level,
    geom
FROM sofiaplan_flood_risk_low;

DROP VIEW IF EXISTS sofiaplan_flood_risk_medium_tiles;
CREATE VIEW sofiaplan_flood_risk_medium_tiles AS
SELECT
    id,
    COALESCE(properties->>'eu_cd_hp', '') AS label,
    COALESCE(properties->>'eu_cd_hp', '') AS zone_id,
    2                                     AS score,
    2                                     AS risk_level,
    geom
FROM sofiaplan_flood_risk_medium;

DROP VIEW IF EXISTS sofiaplan_flood_risk_high_tiles;
CREATE VIEW sofiaplan_flood_risk_high_tiles AS
SELECT
    id,
    COALESCE(properties->>'apsfr', '')   AS label,
    COALESCE(properties->>'apsfr', '')   AS zone_id,
    3                                    AS score,
    3                                    AS risk_level,
    geom
FROM sofiaplan_flood_risk_high;
