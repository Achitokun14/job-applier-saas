-- 000002_add_fulltext_search.up.sql
-- Adds PostgreSQL full-text search to the jobs table.

-- Add tsvector column for full-text search
ALTER TABLE jobs ADD COLUMN search_vector tsvector;

-- Create GIN index on the search_vector column
CREATE INDEX idx_jobs_search_vector ON jobs USING GIN(search_vector);

-- Populate existing rows
UPDATE jobs SET search_vector = to_tsvector('english',
    coalesce(title, '') || ' ' ||
    coalesce(company, '') || ' ' ||
    coalesce(description, '')
);

-- Create trigger function to auto-update search_vector on INSERT/UPDATE
CREATE OR REPLACE FUNCTION jobs_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english',
        coalesce(NEW.title, '') || ' ' ||
        coalesce(NEW.company, '') || ' ' ||
        coalesce(NEW.description, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER trg_jobs_search_vector_update
    BEFORE INSERT OR UPDATE ON jobs
    FOR EACH ROW
    EXECUTE FUNCTION jobs_search_vector_update();
