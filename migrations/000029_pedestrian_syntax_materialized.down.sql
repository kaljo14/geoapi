-- Restore the live view from migration 000022.
DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_syntax_tiles;

CREATE VIEW sofiaplan_pedestrian_syntax_tiles AS
SELECT
    id,
    score,
    choice,
    connectivity,
    segment_length,
    neighborhood,
    PERCENT_RANK() OVER (
        PARTITION BY neighborhood_id
        ORDER BY score
    ) AS local_percentile,
    geom::geometry(Geometry, 4326) AS geom
FROM (
    SELECT
        p.id,
        COALESCE((p.properties->>'t1024_inte')::numeric, 0)  AS score,
        COALESCE((p.properties->>'t1024_choi')::numeric, 0)  AS choice,
        COALESCE((p.properties->>'connectivi')::numeric, 0)  AS connectivity,
        COALESCE((p.properties->>'segment_le')::numeric, 0)  AS segment_length,
        n.id                                                  AS neighborhood_id,
        (n.properties->>'kvname')                            AS neighborhood,
        p.geom
    FROM sofiaplan_pedestrian_syntax p
    LEFT JOIN sofiaplan_neighborhoods n
        ON ST_Contains(n.geom, ST_Centroid(p.geom))
) matched;
