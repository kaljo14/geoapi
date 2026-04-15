-- Calibrated foot traffic predictions per street segment.
--
-- Percentile-anchored scaling from space syntax integration scores to
-- estimated pedestrians per hour, calibrated against 36 radar sensors
-- in Lozenets (non-operational sensors removed from DB).
--
-- Sensor stats: min ≈2.5, median ≈15, P75 ≈29, max ≈106 ped/hr
-- Sensors sit at ~P63–P99 of Sofia-wide integration scores
--
-- Formula: predicted_hourly = 50 * (global_percentile ^ 2.5)
--
-- N = 36 sensors, July 2024 data from GATE Institute City Living Lab.

CREATE MATERIALIZED VIEW calibrated_foot_traffic_tiles AS
SELECT
    id,
    neighborhood,
    integration,
    choice,
    segment_length,
    global_percentile,
    GREATEST(0, ROUND(
        50.0 * POWER(global_percentile, 2.5)
    ))::int AS predicted_hourly,
    geom
FROM (
    SELECT
        t.id,
        t.neighborhood,
        t.score         AS integration,
        t.choice,
        t.segment_length,
        PERCENT_RANK() OVER (ORDER BY t.score) AS global_percentile,
        t.geom
    FROM sofiaplan_pedestrian_syntax_tiles t
) ranked;

CREATE INDEX idx_calibrated_ft_geom
    ON calibrated_foot_traffic_tiles USING GIST(geom);
