CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL UNIQUE,
    occupation VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
/*   
owner_id VARCHAR(255) REFERENCES users(id),
*/

CREATE TABLE images (
    id UUID PRIMARY KEY,
    owner_id UUID,
    unique_file_name TEXT NOT NULL UNIQUE,
    tags TEXT[],
    description TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);