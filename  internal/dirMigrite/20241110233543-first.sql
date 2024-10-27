
-- +migrate Up
INSERT INTO products (name, hobby, price) VALUES
('Prigojin', 'runer', 1);
-- +migrate Down
DELETE FROM products WHERE name IN ('Prigojin');
