-- Pedestrian sensor locations
CREATE TABLE pedestrian_sensors (
    device_id   TEXT PRIMARY KEY,
    lat         DOUBLE PRECISION NOT NULL,
    lng         DOUBLE PRECISION NOT NULL,
    geom        GEOMETRY(Point, 4326) NOT NULL
);
CREATE INDEX idx_pedestrian_sensors_geom ON pedestrian_sensors USING GIST (geom);

-- Raw 10-minute interval readings
CREATE TABLE pedestrian_readings (
    id          BIGSERIAL PRIMARY KEY,
    device_id   TEXT NOT NULL REFERENCES pedestrian_sensors(device_id),
    recorded_at TIMESTAMPTZ NOT NULL,
    count_left  INT NOT NULL DEFAULT 0,
    count_right INT NOT NULL DEFAULT 0,
    UNIQUE (device_id, recorded_at)
);
CREATE INDEX idx_pedestrian_readings_device ON pedestrian_readings (device_id);
CREATE INDEX idx_pedestrian_readings_time   ON pedestrian_readings (recorded_at);

-- Aggregated tile view for Martin (auto-discovered)
CREATE OR REPLACE VIEW pedestrian_sensors_tiles AS
SELECT
    s.device_id,
    COALESCE(SUM(r.count_left + r.count_right), 0)::INT AS total_pedestrians,
    COALESCE(SUM(r.count_left), 0)::INT  AS total_left,
    COALESCE(SUM(r.count_right), 0)::INT AS total_right,
    COALESCE(COUNT(r.id), 0)::INT        AS reading_count,
    s.geom
FROM pedestrian_sensors s
LEFT JOIN pedestrian_readings r ON r.device_id = s.device_id
GROUP BY s.device_id, s.geom;
