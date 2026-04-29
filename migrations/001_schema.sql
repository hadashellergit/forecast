-- 001_schema.sql
-- Core schema. Run once on first boot.

CREATE TABLE IF NOT EXISTS stores (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- Raw sales data — the source of truth for the algorithm.
-- In production this would be written by a POS integration.
-- For the demo it is populated by the seed script.
CREATE TABLE IF NOT EXISTS sales (
    id         BIGSERIAL PRIMARY KEY,
    store_id   INT       NOT NULL REFERENCES stores(id),
    product_id INT       NOT NULL REFERENCES products(id),
    sold_at    TIMESTAMPTZ NOT NULL,
    quantity   INT       NOT NULL CHECK (quantity > 0)
);

-- Index for the AvgSalesLastNDays query pattern.
CREATE INDEX IF NOT EXISTS sales_store_product_time
    ON sales (store_id, product_id, sold_at);

-- Forecast output table. One row per store+product+date+hour.
-- ON CONFLICT DO UPDATE in the upsert makes regeneration idempotent.
CREATE TABLE IF NOT EXISTS forecasts (
    store_id      INT  NOT NULL REFERENCES stores(id),
    product_id    INT  NOT NULL REFERENCES products(id),
    forecast_date DATE NOT NULL,
    hour          INT  NOT NULL CHECK (hour BETWEEN 0 AND 23),
    quantity      NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (store_id, product_id, forecast_date, hour)
);
