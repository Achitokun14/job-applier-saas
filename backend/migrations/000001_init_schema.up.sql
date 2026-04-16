-- 000001_init_schema.up.sql
-- Creates all initial tables for the job-applier-saas application.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    email_verified BOOLEAN DEFAULT FALSE,
    two_factor_secret VARCHAR(500),
    two_factor_enabled BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Resumes table
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    personal_info TEXT,
    education TEXT,
    experience TEXT,
    skills TEXT,
    projects TEXT,
    achievements TEXT,
    certifications TEXT,
    languages TEXT,
    pdf_path VARCHAR(500)
);

CREATE INDEX IF NOT EXISTS idx_resumes_user_id ON resumes(user_id);

-- Jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    external_id VARCHAR(255),
    title VARCHAR(500),
    company VARCHAR(500),
    location VARCHAR(500),
    description TEXT,
    url TEXT,
    source VARCHAR(100),
    remote BOOLEAN DEFAULT FALSE,
    salary VARCHAR(255)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_external_id ON jobs(external_id);
CREATE INDEX IF NOT EXISTS idx_jobs_title_company ON jobs(title, company);
CREATE INDEX IF NOT EXISTS idx_jobs_source_created_at ON jobs(source, created_at);

-- Applications table
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'applied',
    resume_pdf VARCHAR(500),
    cover_pdf VARCHAR(500),
    notes TEXT,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_user_id_status ON applications(user_id, status);
CREATE INDEX IF NOT EXISTS idx_applications_job_id ON applications(job_id);

-- Settings table
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    llm_provider VARCHAR(50) DEFAULT 'openai',
    llm_model VARCHAR(100) DEFAULT 'gpt-4o-mini',
    llm_api_key VARCHAR(500),
    job_search_remote BOOLEAN DEFAULT TRUE,
    job_search_hybrid BOOLEAN DEFAULT TRUE,
    job_search_onsite BOOLEAN DEFAULT FALSE,
    experience_level VARCHAR(50) DEFAULT 'mid_senior',
    job_types VARCHAR(100) DEFAULT 'full_time',
    positions TEXT,
    locations TEXT,
    distance INTEGER DEFAULT 50,
    company_blacklist TEXT,
    title_blacklist TEXT
);

CREATE INDEX IF NOT EXISTS idx_settings_user_id ON settings(user_id);

-- Refresh tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);

-- Password reset tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
