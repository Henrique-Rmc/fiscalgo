CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    unique_file_name TEXT NOT NULL UNIQUE,
    tags TEXT[],
    description TEXT,
    url TEXT NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);