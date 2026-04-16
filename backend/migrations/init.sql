-- PostgreSQL initialization script
-- This runs automatically when the database is first created

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255)
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_email ON users(email);

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

CREATE INDEX idx_resumes_user_id ON resumes(user_id);

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

CREATE INDEX idx_jobs_external_id ON jobs(external_id);
CREATE INDEX idx_jobs_title ON jobs(title);
CREATE INDEX idx_jobs_company ON jobs(company);
CREATE INDEX idx_jobs_source ON jobs(source);

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

CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_job_id ON applications(job_id);

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

CREATE INDEX idx_settings_user_id ON settings(user_id);

-- Insert sample data (optional - remove for production)
-- INSERT INTO users (email, password, name) VALUES 
--     ('demo@example.com', '$2a$10$demo_hashed_password', 'Demo User');

-- Grant permissions (adjust as needed)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO jobapplier;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO jobapplier;
