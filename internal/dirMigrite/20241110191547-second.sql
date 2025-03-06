
-- +migrate Up

INSERT INTO products (name, hobby, price) VALUES
                                              ('Alice', 'Gardening', 100),
                                              ('Bob', 'Painting', 200),
                                              ('Charlie', 'Cooking', 150),
                                              ('Diana', 'Photography', 300),
                                              ('Eve', 'Reading', 50);

-- +migrate Down
DELETE FROM products WHERE name IN ('Alice', 'Bob', 'Charlie', 'Diana', 'Eve');
