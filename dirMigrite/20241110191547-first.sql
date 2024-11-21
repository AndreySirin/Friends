
-- +migrate Up

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hobby VARCHAR(255) NOT NULL,
    price INT
    );

-- +migrate Down
DROP TABLE IF EXISTS products;