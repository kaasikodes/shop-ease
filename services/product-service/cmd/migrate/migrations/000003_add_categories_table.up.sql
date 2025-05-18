CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    entityId INT NOT NULL,
    tokenType VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    expiresAt TIMESTAMP NOT NULL,
    createdAt TIMESTAMP DEFAULT NOW(),
    updatedAt TIMESTAMP DEFAULT NOW(),
    UNIQUE KEY (value(255)) -- Specify a key length for TEXT, so only yhe first 255 characters are indexed
);
