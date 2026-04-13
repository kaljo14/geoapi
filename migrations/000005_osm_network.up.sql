CREATE TABLE IF NOT EXISTS osm_nodes (
    osm_id     BIGINT PRIMARY KEY,
    lat        FLOAT  NOT NULL,
    lng        FLOAT  NOT NULL,
    geom       GEOMETRY(Point, 4326),
    walk_score FLOAT  NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_osm_nodes_geom ON osm_nodes USING GIST(geom);

CREATE TABLE IF NOT EXISTS osm_edges (
    osm_way_id   BIGINT NOT NULL,
    from_node_id BIGINT NOT NULL,
    to_node_id   BIGINT NOT NULL,
    highway      TEXT,
    geom         GEOMETRY(LineString, 4326),
    walk_score   FLOAT  NOT NULL DEFAULT 0,
    PRIMARY KEY (osm_way_id, from_node_id, to_node_id)
);

CREATE INDEX IF NOT EXISTS idx_osm_edges_geom     ON osm_edges USING GIST(geom);
CREATE INDEX IF NOT EXISTS idx_osm_edges_from_node ON osm_edges(from_node_id);
CREATE INDEX IF NOT EXISTS idx_osm_edges_to_node   ON osm_edges(to_node_id);



  docker exec -i neofyis-geopulse-postgres-1 psql -U geopulse -d geopulse </Users/kaloyanivanov/Personal-SRV/apis/barbershops_plain.sql