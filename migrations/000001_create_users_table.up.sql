CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    cpf VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    occupation VARCHAR(255) NOT NULL,
    professional_registry VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clients(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    cpf VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    asks_invoice BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS revenues (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    client_id UUID NOT NULL REFERENCES clients(id),
    beneficiary_cpf_cnpj VARCHAR(255) NOT NULL,
    value DECIMAL(10, 2) NOT NULL,
    total_paid DECIMAL(10,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    issue_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
/*   
Se Issue date for vazio, é o mesmo do Created_at
*/
CREATE TABLE IF NOT EXISTS invoices(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    description TEXT NOT NULL,
    value DECIMAL(10, 2) NOT NULL,
    expense_category VARCHAR(255) NOT NULL,
    access_key VARCHAR(44),
    image_url VARCHAR(255),
    issue_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clients_user_id ON clients (user_id);
CREATE INDEX idx_revenues_user_id ON revenues (user_id);
CREATE INDEX idx_revenues_client_id ON revenues (client_id);


INSERT INTO users(id, name, email, cpf, password_hash, occupation, professional_registry)
VALUES(
    '6daa7ce0-6594-43ed-b583-c74bd6aa1a13', 
    'joao',
    'joao@email',
    '07338057308',
    'hash_da_senha_gerado_pelo_bcrypt_aqui', 
    'dentista',
    'CRM-3333'
);

INSERT INTO clients(id,user_id, name, cpf, phone, email, asks_invoice)
VALUES(
    'ee31d0ea-14ce-45fd-b7d4-88beffd0c58c', 
    '6daa7ce0-6594-43ed-b583-c74bd6aa1a13', 
    'cliente Maria', 
    '07665846235', 
    '887797451', 
    'clienteMaria@email', 
    FALSE
);

INSERT INTO clients(id,user_id, name, cpf, phone, email, asks_invoice)
VALUES(
    '142c5b1b-4f27-4594-9be6-098e5f8a1216', 
    '6daa7ce0-6594-43ed-b583-c74bd6aa1a13', 
    'cliente Joao', 
    '08457498536', 
    '887797455', 
    'clienteJoaoa@email', 
    FALSE
);
