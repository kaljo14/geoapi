-- Revert to original minimal flood risk tile views
CREATE OR REPLACE VIEW sofiaplan_flood_risk_low_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district, 1 AS risk_level, geom FROM sofiaplan_flood_risk_low;

CREATE OR REPLACE VIEW sofiaplan_flood_risk_medium_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district, 2 AS risk_level, geom FROM sofiaplan_flood_risk_medium;

CREATE OR REPLACE VIEW sofiaplan_flood_risk_high_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district, 3 AS risk_level, geom FROM sofiaplan_flood_risk_high;
