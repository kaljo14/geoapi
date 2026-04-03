ALTER TABLE places ADD COLUMN IF NOT EXISTS estimated_monthly_visitors INTEGER;
ALTER TABLE places ADD COLUMN IF NOT EXISTS visitor_location_score NUMERIC(5,3);
