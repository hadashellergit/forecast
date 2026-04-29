-- 002_seed_stores_products.sql
-- Static reference data. Safe to re-run (ON CONFLICT DO NOTHING).

INSERT INTO stores (id, name, city) VALUES
    (1, 'KFC Times Square',    'New York'),
    (2, 'KFC Hollywood Blvd',  'Los Angeles'),
    (3, 'KFC Michigan Ave',    'Chicago')
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, name) VALUES
    (1, 'Original Recipe Bucket'),
    (2, 'Zinger Burger'),
    (3, 'Popcorn Chicken'),
    (4, 'Crispy Strips'),
    (5, 'Coleslaw')
ON CONFLICT (id) DO NOTHING;
