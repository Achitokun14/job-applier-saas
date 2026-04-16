-- 000002_add_fulltext_search.down.sql
-- Removes full-text search from the jobs table.

DROP TRIGGER IF EXISTS trg_jobs_search_vector_update ON jobs;
DROP FUNCTION IF EXISTS jobs_search_vector_update();
DROP INDEX IF EXISTS idx_jobs_search_vector;
ALTER TABLE jobs DROP COLUMN IF EXISTS search_vector;
