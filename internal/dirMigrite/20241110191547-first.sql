
-- +migrate Up

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(name) >= 2),
    hobby VARCHAR(255) NOT NULL CHECK (LENGTH(hobby) >= 2),
    price INT NOT NULL CHECK(price !=0)
    );

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(name) >= 2),
    email VARCHAR(255) UNIQUE NOT NULL CHECK (POSITION('@' IN email) > 1),
    password TEXT NOT NULL CHECK (LENGTH(password) >= 6),
    registeredAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    id_user INT NOT NULL,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
    );

-- +migrate Down
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS refresh_tokens;