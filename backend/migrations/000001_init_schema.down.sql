-- 000001_init_schema.down.sql
-- Drops all tables in reverse order of creation.

DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS resumes;
DROP TABLE IF EXISTS users;
