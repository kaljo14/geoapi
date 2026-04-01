-- Drop views first, then tables (reverse creation order)

DROP VIEW IF EXISTS sofiaplan_urban_morphology_ge_tiles;
DROP VIEW IF EXISTS sofiaplan_residential_typology_ge_tiles;
DROP VIEW IF EXISTS sofiaplan_building_footprint_ge_tiles;
DROP VIEW IF EXISTS sofiaplan_building_density_ge_tiles;

DROP TABLE IF EXISTS sofiaplan_urban_morphology_ge;
DROP TABLE IF EXISTS sofiaplan_residential_typology_ge;
DROP TABLE IF EXISTS sofiaplan_building_footprint_ge;
DROP TABLE IF EXISTS sofiaplan_building_density_ge;
