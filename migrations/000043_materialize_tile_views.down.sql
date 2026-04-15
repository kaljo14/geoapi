-- Restore all materialized views back to live views.
-- This reverses 000043_materialize_tile_views.up.sql.

-- =========================================================
-- From 000008: sofiaplan_tile_views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_income_tiles;
CREATE VIEW sofiaplan_income_tiles AS
SELECT id, (properties->>'obns_cyr') AS label, (properties->>'obns_num') AS obns_num,
    COALESCE((properties->>'sr_mes_dohod')::numeric, 0) AS score,
    COALESCE((properties->>'sr_god_dohod')::numeric, 0) AS sr_god_dohod,
    COALESCE((properties->>'percent')::numeric, 0) AS percent,
    COALESCE((properties->>'dohod_ofis')::numeric, 0) AS dohod_ofis, geom
FROM sofiaplan_income;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_zoning_tiles;
CREATE VIEW sofiaplan_zoning_tiles AS
SELECT id, (properties->>'type_') AS label, (properties->>'center') AS center, (properties->>'new_end') AS new_end,
    abs(hashtext(COALESCE(properties->>'type_', ''))) % 100 AS score, geom
FROM sofiaplan_zoning;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_business_turnover_tiles;
CREATE VIEW sofiaplan_business_turnover_tiles AS
SELECT id, (properties->>'obns_cyr') AS label, (properties->>'obns_num') AS obns_num,
    COALESCE((properties->>'oborot_2016')::numeric, 0) AS score,
    COALESCE((properties->>'oborot_2010')::numeric, 0) AS oborot_2010,
    COALESCE((properties->>'prirast')::numeric, 0) AS prirast,
    COALESCE((properties->>'prirast_perc')::numeric, 0) AS prirast_perc,
    COALESCE((properties->>'dial_2016')::numeric, 0) AS dial_2016,
    COALESCE((properties->>'norma_raztej')::numeric, 0) AS norma_raztej, geom
FROM sofiaplan_business_turnover;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_property_prices_tiles;
CREATE VIEW sofiaplan_property_prices_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'kvartal') AS label,
    (properties->>'kvartal') AS kvartal, (properties->>'regname') AS regname, (properties->>'godina') AS godina,
    COALESCE((properties->>'cena_ap_kv_m')::numeric, 0) AS score,
    COALESCE((properties->>'cena_ap')::numeric, 0) AS cena_ap,
    COALESCE((properties->>'cena_ofis')::numeric, 0) AS cena_ofis,
    COALESCE((properties->>'cena_ofis_kv_m')::numeric, 0) AS cena_ofis_kv_m,
    COALESCE((properties->>'naem_ap')::numeric, 0) AS naem_ap,
    COALESCE((properties->>'naem_ap_kv_m')::numeric, 0) AS naem_ap_kv_m,
    COALESCE((properties->>'naem_ofis')::numeric, 0) AS naem_ofis,
    COALESCE((properties->>'naem_ofis_kv_m')::numeric, 0) AS naem_ofis_kv_m,
    COALESCE((properties->>'naem_mag')::numeric, 0) AS naem_mag,
    COALESCE((properties->>'naem_mag_kv_m')::numeric, 0) AS naem_mag_kv_m, geom
FROM sofiaplan_property_prices;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_metro_catchments_tiles;
CREATE VIEW sofiaplan_metro_catchments_tiles AS
SELECT id, (properties->>'name') AS label, COALESCE((properties->>'tobreak')::numeric, 0) AS score,
    COALESCE((properties->>'frombreak')::numeric, 0) AS frombreak, (properties->>'facilityid') AS facilityid, geom
FROM sofiaplan_metro_catchments;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_population_grid_tiles;
CREATE VIEW sofiaplan_population_grid_tiles AS
SELECT id, (properties->>'grd_id') AS label, COALESCE((properties->>'tot_p')::numeric, 0) AS score,
    COALESCE((properties->>'tot_m')::numeric, 0) AS tot_m, COALESCE((properties->>'tot_f')::numeric, 0) AS tot_f,
    COALESCE((properties->>'t_00_14')::numeric, 0) AS t_00_14, COALESCE((properties->>'t_15_64')::numeric, 0) AS t_15_64,
    COALESCE((properties->>'t_65_')::numeric, 0) AS t_65_plus, geom
FROM sofiaplan_population_grid;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_neighborhoods_tiles;
CREATE VIEW sofiaplan_neighborhoods_tiles AS
SELECT id, (properties->>'kvname') AS label, (properties->>'prefname') AS prefname,
    COALESCE((properties->>'adm_kvid')::numeric, 0) AS adm_kvid,
    COALESCE((properties->>'type_kv')::numeric, 0) AS score, geom
FROM sofiaplan_neighborhoods;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_development_potential_tiles;
CREATE VIEW sofiaplan_development_potential_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon') AS label,
    (properties->>'rajon') AS rajon, (properties->>'zastr_potenc_txt') AS description,
    COALESCE((properties->>'zastr_potencial')::numeric, 0) AS score, geom
FROM sofiaplan_development_potential;

-- =========================================================
-- From 000010: demographics tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_census_addresses_tiles;
CREATE VIEW sofiaplan_census_addresses_tiles AS
SELECT id, (properties->>'nstreetnam') AS label, (properties->>'nnumber') AS street_number,
    (properties->>'ecode_rayon') AS rayon, COALESCE((properties->>'nbroi_lica')::numeric, 0) AS score,
    COALESCE((properties->>'nn_jilisht')::numeric, 0) AS dwellings,
    COALESCE((properties->>'nmale_sum')::numeric, 0) AS male, COALESCE((properties->>'nfemale_su')::numeric, 0) AS female,
    COALESCE((properties->>'nage0_14')::numeric, 0) AS age_0_14, COALESCE((properties->>'nage15_24')::numeric, 0) AS age_15_24,
    COALESCE((properties->>'nage25_34')::numeric, 0) AS age_25_34, COALESCE((properties->>'nage35_44')::numeric, 0) AS age_35_44,
    COALESCE((properties->>'nage45_54')::numeric, 0) AS age_45_54, COALESCE((properties->>'nage55_64')::numeric, 0) AS age_55_64,
    COALESCE((properties->>'nage65_')::numeric, 0) AS age_65_plus, geom
FROM sofiaplan_census_addresses;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_demographic_forecast_tiles;
CREATE VIEW sofiaplan_demographic_forecast_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon') AS label, (properties->>'rajon') AS rajon,
    COALESCE((properties->>'nbroi_lica_sum')::numeric, 0) AS score,
    COALESCE((properties->>'n2017_sum')::numeric, 0) AS pop_2017,
    COALESCE((properties->>'pp2030_sum')::numeric, 0) AS forecast_2030,
    COALESCE((properties->>'pp2040_sum')::numeric, 0) AS forecast_2040,
    COALESCE((properties->>'pp2050_sum')::numeric, 0) AS forecast_2050,
    COALESCE((properties->>'nmale_sum_sum')::numeric, 0) AS male,
    COALESCE((properties->>'nfemale_su_sum')::numeric, 0) AS female,
    COALESCE((properties->>'nage0_14_sum')::numeric, 0) AS age_0_14,
    COALESCE((properties->>'nage65__sum')::numeric, 0) AS age_65_plus, geom
FROM sofiaplan_demographic_forecast;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_demographic_forecast_ge_tiles;
CREATE VIEW sofiaplan_demographic_forecast_ge_tiles AS
SELECT id, COALESCE(properties->>'ime_predl', properties->>'regname') AS label,
    (properties->>'rajon') AS rajon, (properties->>'ident_gr_u') AS ge_id,
    COALESCE((properties->>'nbroi_lica')::numeric, 0) AS score,
    COALESCE((properties->>'n2017')::numeric, 0) AS pop_2017,
    COALESCE((properties->>'pp2030')::numeric, 0) AS forecast_2030,
    COALESCE((properties->>'pp2040')::numeric, 0) AS forecast_2040,
    COALESCE((properties->>'pp2050')::numeric, 0) AS forecast_2050,
    COALESCE((properties->>'nmale_sum')::numeric, 0) AS male,
    COALESCE((properties->>'nfemale_su')::numeric, 0) AS female, geom
FROM sofiaplan_demographic_forecast_ge;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_population_potential_tiles;
CREATE VIEW sofiaplan_population_potential_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon') AS label, (properties->>'rajon') AS rajon,
    (properties->>'posoka') AS direction, COALESCE((properties->>'ppl_30kvm_sum')::numeric, 0) AS score,
    COALESCE((properties->>'ppl_35kvm_sum')::numeric, 0) AS ppl_35,
    COALESCE((properties->>'ppl_40kvm_sum')::numeric, 0) AS ppl_40,
    COALESCE((properties->>'dens30_ppl_kvkm')::numeric, 0) AS density_30,
    COALESCE((properties->>'dens40_ppl_kvkm')::numeric, 0) AS density_40, geom
FROM sofiaplan_population_potential;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_residential_load_tiles;
CREATE VIEW sofiaplan_residential_load_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon') AS label, (properties->>'rajon') AS rajon,
    COALESCE((properties->>'ppl_30kvm')::numeric, 0) AS score,
    COALESCE((properties->>'ppl_35kvm')::numeric, 0) AS ppl_35,
    COALESCE((properties->>'ppl_40kvm')::numeric, 0) AS ppl_40,
    COALESCE((properties->>'rzp_jilsgr')::numeric, 0) AS rzp_residential,
    COALESCE((properties->>'rzp_all')::numeric, 0) AS rzp_total,
    COALESCE((properties->>'dens30_ppl_kvkm')::numeric, 0) AS density_30,
    COALESCE((properties->>'dens40_ppl_kvkm')::numeric, 0) AS density_40, geom
FROM sofiaplan_residential_load;

-- =========================================================
-- From 000012: transport tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_transit_access_ge_tiles;
CREATE VIEW sofiaplan_transit_access_ge_tiles AS
SELECT id, (properties->>'ime_predl') AS label, (properties->>'id') AS ge_id, (properties->>'rajon') AS rajon,
    COALESCE((properties->>'kgr')::numeric, 0) AS score, geom
FROM sofiaplan_transit_access_ge;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_transit_access_district_tiles;
CREATE VIEW sofiaplan_transit_access_district_tiles AS
SELECT id, (properties->>'name') AS label, (properties->>'facilityid') AS district_id,
    COALESCE((properties->>'total_leng')::numeric, 0) AS score, geom
FROM sofiaplan_transit_access_district;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_metro_access_800m_tiles;
CREATE VIEW sofiaplan_metro_access_800m_tiles AS
SELECT id, (properties->>'name') AS label, COALESCE((properties->>'tobreak')::numeric, 800) AS score,
    COALESCE((properties->>'frombreak')::numeric, 0) AS frombreak, (properties->>'facilityid') AS facilityid, geom
FROM sofiaplan_metro_access_800m;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_metro_access_1200m_tiles;
CREATE VIEW sofiaplan_metro_access_1200m_tiles AS
SELECT id, (properties->>'name') AS label, COALESCE((properties->>'tobreak')::numeric, 1200) AS score,
    COALESCE((properties->>'frombreak')::numeric, 0) AS frombreak, (properties->>'facilityid') AS facilityid, geom
FROM sofiaplan_metro_access_1200m;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_bus_lines_tiles;
CREATE VIEW sofiaplan_bus_lines_tiles AS
SELECT id, (properties->>'line_bus') AS label, (properties->>'line_bus') AS route_id, geom
FROM sofiaplan_bus_lines;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_bus_lines_alt_tiles;
CREATE VIEW sofiaplan_bus_lines_alt_tiles AS
SELECT id, (properties->>'ЛИНИЯ') AS label, (properties->>'ЛИНИЯ') AS route_id, geom
FROM sofiaplan_bus_lines_alt;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_trolleybus_lines_tiles;
CREATE VIEW sofiaplan_trolleybus_lines_tiles AS
SELECT id, (properties->>'line_tb') AS label, (properties->>'line_tb') AS route_id, geom
FROM sofiaplan_trolleybus_lines;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_tram_lines_tiles;
CREATE VIEW sofiaplan_tram_lines_tiles AS
SELECT id, (properties->>'line_tram') AS label, (properties->>'line_tram') AS route_id, geom
FROM sofiaplan_tram_lines;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_tram_lines_alt_tiles;
CREATE VIEW sofiaplan_tram_lines_alt_tiles AS
SELECT id, (properties->>'line_tram') AS label, (properties->>'line_tram') AS route_id, geom
FROM sofiaplan_tram_lines_alt;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_railway_stations_tiles;
CREATE VIEW sofiaplan_railway_stations_tiles AS
SELECT id, (properties->>'tradename') AS label,
    COALESCE((properties->>'2019_prist')::numeric, 0) + COALESCE((properties->>'2019_zamin')::numeric, 0) AS score, geom
FROM sofiaplan_railway_stations;

-- =========================================================
-- From 000021: cycling tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_cycling_network_tiles;
CREATE VIEW sofiaplan_cycling_network_tiles AS
SELECT id, (properties->>'type') AS label, (properties->>'type') AS path_type,
    (properties->>'posoka') AS direction, COALESCE((properties->>'length')::numeric, 0) AS length_m, geom
FROM sofiaplan_cycling_network;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_cycling_network_alt_tiles;
CREATE VIEW sofiaplan_cycling_network_alt_tiles AS
SELECT id, (properties->>'type') AS label, (properties->>'type') AS path_type,
    (properties->>'posoka') AS direction, geom
FROM sofiaplan_cycling_network_alt;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_cycling_planned_tiles;
CREATE VIEW sofiaplan_cycling_planned_tiles AS
SELECT id, (properties->>'name') AS label, COALESCE((properties->>'priority')::numeric, 0) AS priority,
    COALESCE((properties->>'project')::numeric, 0) AS project, (properties->>'note') AS note, geom
FROM sofiaplan_cycling_planned;

-- =========================================================
-- From 000024: health tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_health_service_concentration_tiles;
CREATE VIEW sofiaplan_health_service_concentration_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon', '') AS label,
    (properties->>'rajon') AS district, COALESCE((properties->>'numpoints')::numeric, 0) AS score, geom
FROM sofiaplan_health_service_concentration;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_health_infrastructure_concentration_tiles;
CREATE VIEW sofiaplan_health_infrastructure_concentration_tiles AS
SELECT id, COALESCE(properties->>'regname', properties->>'rajon', '') AS label,
    (properties->>'rajon') AS district, COALESCE((properties->>'numpoints')::numeric, 0) AS score, geom
FROM sofiaplan_health_infrastructure_concentration;

-- =========================================================
-- From 000025: building layer tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_building_density_ge_tiles;
CREATE VIEW sofiaplan_building_density_ge_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district,
    COALESCE((properties->>'zastr_plytnost')::numeric, 0) AS score,
    COALESCE((properties->>'zastr_intenzivnost')::numeric, 0) AS intensity,
    COALESCE((properties->>'zastr_sklyuchenost')::numeric, 0) AS enclosure_ratio,
    COALESCE((properties->>'sredna_etajnost')::numeric, 0) AS avg_floors, geom
FROM sofiaplan_building_density_ge;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_building_footprint_ge_tiles;
CREATE VIEW sofiaplan_building_footprint_ge_tiles AS
SELECT id, COALESCE(properties->>'ge_id', '') AS ge_id, COALESCE(properties->>'funktyp_gen_txt', '') AS label,
    '' AS district, COALESCE((properties->>'rzp')::numeric, 0) AS score,
    COALESCE((properties->>'zp')::numeric, 0) AS zp, COALESCE((properties->>'rzp')::numeric, 0) AS rzp,
    COALESCE((properties->>'ge_sgradi_broi')::numeric, 0) AS avg_floors, geom
FROM sofiaplan_building_footprint_ge;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_residential_typology_ge_tiles;
CREATE VIEW sofiaplan_residential_typology_ge_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district,
    COALESCE((properties->>'dial_ednfa')::numeric, 0) AS single_pct,
    COALESCE((properties->>'dial_mnfam')::numeric, 0) AS multi_pct,
    COALESCE((properties->>'dial_panel')::numeric, 0) AS panel_pct,
    CASE
        WHEN COALESCE((properties->>'dial_panel')::numeric, 0) > COALESCE((properties->>'dial_ednfa')::numeric, 0)
         AND COALESCE((properties->>'dial_panel')::numeric, 0) > COALESCE((properties->>'dial_mnfam')::numeric, 0) THEN 3
        WHEN COALESCE((properties->>'dial_mnfam')::numeric, 0) > COALESCE((properties->>'dial_ednfa')::numeric, 0) THEN 2
        WHEN COALESCE((properties->>'dial_ednfa')::numeric, 0) > 0 THEN 1
        ELSE 0
    END AS score, geom
FROM sofiaplan_residential_typology_ge;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_urban_morphology_ge_tiles;
CREATE VIEW sofiaplan_urban_morphology_ge_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district,
    CASE
        WHEN COALESCE((properties->>'bui1_perc')::numeric, 0) + COALESCE((properties->>'bui2_perc')::numeric, 0) >= 65
          AND COALESCE((properties->>'gl_dens')::numeric, 0) < 8 THEN 1
        WHEN COALESCE((properties->>'bui4_perc')::numeric, 0) + COALESCE((properties->>'bui5_perc')::numeric, 0)
           + COALESCE((properties->>'bui6_perc')::numeric, 0) + COALESCE((properties->>'bui78_arep')::numeric, 0) >= 10 THEN 3
        WHEN COALESCE((properties->>'gl_dens')::numeric, 0) >= 15 THEN 2
        ELSE 5
    END AS score, COALESCE((properties->>'gl_dens')::numeric, 0) AS density, geom
FROM sofiaplan_urban_morphology_ge;

-- =========================================================
-- From 000027: pedestrian network tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_city_tiles;
CREATE VIEW sofiaplan_pedestrian_city_tiles AS
SELECT id, COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0) AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0) AS segment_length, geom
FROM sofiaplan_pedestrian_city;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_city_alt_tiles;
CREATE VIEW sofiaplan_pedestrian_city_alt_tiles AS
SELECT id, COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0) AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0) AS segment_length, geom
FROM sofiaplan_pedestrian_city_alt;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_municipality_tiles;
CREATE VIEW sofiaplan_pedestrian_municipality_tiles AS
SELECT id, COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0) AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0) AS segment_length, geom
FROM sofiaplan_pedestrian_municipality;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_municipality_alt_tiles;
CREATE VIEW sofiaplan_pedestrian_municipality_alt_tiles AS
SELECT id, COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0) AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0) AS segment_length, geom
FROM sofiaplan_pedestrian_municipality_alt;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_segmented_tiles;
CREATE VIEW sofiaplan_pedestrian_segmented_tiles AS
SELECT id, COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice,
    COALESCE((properties->>'connectivi')::numeric, 0) AS connectivity,
    COALESCE((properties->>'segment_le')::numeric, 0) AS segment_length, geom
FROM sofiaplan_pedestrian_segmented;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_pedestrian_integration_tiles;
CREATE VIEW sofiaplan_pedestrian_integration_tiles AS
SELECT id, COALESCE(properties->>'regname', '') AS label, COALESCE(properties->>'rajon', '') AS district,
    COALESCE((properties->>'t1024_inte')::numeric, 0) AS score,
    COALESCE((properties->>'t1024_choi')::numeric, 0) AS choice, geom
FROM sofiaplan_pedestrian_integration;

-- =========================================================
-- From 000033: flood risk tile views
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_flood_risk_low_tiles;
CREATE VIEW sofiaplan_flood_risk_low_tiles AS
SELECT id, COALESCE(properties->>'apsfr', '') AS label, COALESCE(properties->>'apsfr', '') AS zone_id,
    1 AS score, 1 AS risk_level, geom
FROM sofiaplan_flood_risk_low;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_flood_risk_medium_tiles;
CREATE VIEW sofiaplan_flood_risk_medium_tiles AS
SELECT id, COALESCE(properties->>'eu_cd_hp', '') AS label, COALESCE(properties->>'eu_cd_hp', '') AS zone_id,
    2 AS score, 2 AS risk_level, geom
FROM sofiaplan_flood_risk_medium;

DROP MATERIALIZED VIEW IF EXISTS sofiaplan_flood_risk_high_tiles;
CREATE VIEW sofiaplan_flood_risk_high_tiles AS
SELECT id, COALESCE(properties->>'apsfr', '') AS label, COALESCE(properties->>'apsfr', '') AS zone_id,
    3 AS score, 3 AS risk_level, geom
FROM sofiaplan_flood_risk_high;

-- =========================================================
-- From 000015: parking zones tile view
-- =========================================================

DROP MATERIALIZED VIEW IF EXISTS parking_zones_tiles;
CREATE VIEW parking_zones_tiles AS
SELECT id, 'green'::text AS color, geom FROM sofiaplan_parking_green
UNION ALL
SELECT id, 'blue'::text  AS color, geom FROM sofiaplan_parking_blue;
