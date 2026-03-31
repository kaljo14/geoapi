-- Revert sofiaplan_pedestrian_syntax_tiles to the original (global score only, no neighbourhood join).
DROP VIEW IF EXISTS sofiaplan_pedestrian_syntax_tiles;
CREATE VIEW sofiaplan_pedestrian_syntax_tiles AS
SELECT
    id,
    COALESCE((properties->>'t1024_inte')::numeric, 0)    AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0)    AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0)    AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0)    AS segment_length,
    geom
FROM sofiaplan_pedestrian_syntax;
