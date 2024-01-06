CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR,
    created_at TIMESTAMP DEFAULT Now()
);

CREATE TABLE IF NOT EXISTS hashtags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR,
    created_at TIMESTAMP DEFAULT Now()
);

CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR,
    slug VARCHAR,
    description TEXT,
    created_at TIMESTAMP DEFAULT Now()
);

CREATE TABLE IF NOT EXISTS project_hashtags (
    id BIGSERIAL PRIMARY KEY,
    hashtag_id INT REFERENCES hashtags(id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(id)  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_projects (
    id BIGSERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id)  ON DELETE CASCADE,
    user_id INT REFERENCES users(id)  ON DELETE CASCADE
);
