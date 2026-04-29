-- 003_seed_sales.sql
-- Generates 7 days of realistic mock sales so the forecast algorithm
-- has data to average over on first boot.
--
-- Pattern:
--   - Stores open 10:00–22:00
--   - Lunch peak 12–14, dinner peak 18–20
--   - Weekends ~20% busier
--   - Each store has a slight volume multiplier (Times Square busiest)

DO $$
DECLARE
    v_day        DATE;
    v_store_id   INT;
    v_product_id INT;
    v_hour       INT;
    v_base       INT;
    v_qty        INT;
    v_multiplier NUMERIC;
BEGIN
    FOR v_day IN
        SELECT generate_series(
            (CURRENT_DATE - INTERVAL '7 days')::date,
            (CURRENT_DATE - INTERVAL '1 day')::date,
            '1 day'::interval
        )::date
    LOOP
        FOR v_store_id IN SELECT id FROM stores LOOP
            FOR v_product_id IN SELECT id FROM products LOOP
                FOR v_hour IN 10..22 LOOP

                    -- Base units per hour (tuned per product)
                    v_base := CASE v_product_id
                        WHEN 1 THEN 8   -- bucket: slower, high value
                        WHEN 2 THEN 12  -- zinger: most popular
                        WHEN 3 THEN 10  -- popcorn: steady
                        WHEN 4 THEN 9   -- strips: steady
                        WHEN 5 THEN 6   -- coleslaw: side dish
                        ELSE 8
                    END;

                    -- Time-of-day multiplier
                    v_multiplier := CASE
                        WHEN v_hour IN (12, 13) THEN 2.2   -- lunch peak
                        WHEN v_hour = 14        THEN 1.6
                        WHEN v_hour IN (18, 19) THEN 2.0   -- dinner peak
                        WHEN v_hour = 20        THEN 1.5
                        WHEN v_hour IN (10, 11) THEN 0.7   -- morning slow
                        WHEN v_hour = 22        THEN 0.5   -- closing
                        ELSE 1.0
                    END;

                    -- Store volume multiplier
                    v_multiplier := v_multiplier * CASE v_store_id
                        WHEN 1 THEN 1.4   -- Times Square: tourist traffic
                        WHEN 2 THEN 1.1   -- Hollywood: moderate
                        WHEN 3 THEN 1.0   -- Chicago: baseline
                        ELSE 1.0
                    END;

                    -- Weekend bump
                    IF EXTRACT(DOW FROM v_day) IN (0, 6) THEN
                        v_multiplier := v_multiplier * 1.2;
                    END IF;

                    -- Add ±15% random noise so averages aren't perfectly flat
                    v_qty := GREATEST(1, ROUND(
                        v_base * v_multiplier * (0.85 + RANDOM() * 0.30)
                    )::int);

                    INSERT INTO sales (store_id, product_id, sold_at, quantity)
                    VALUES (
                        v_store_id,
                        v_product_id,
                        (v_day + (v_hour || ' hours')::interval),
                        v_qty
                    );

                END LOOP; -- hour
            END LOOP; -- product
        END LOOP; -- store
    END LOOP; -- day
END $$;


INSERT INTO forecasts (store_id, product_id, forecast_date, hour, quantity)
SELECT
    s.store_id,
    s.product_id,
    CURRENT_DATE + 1                    AS forecast_date,
    EXTRACT(HOUR FROM s.sold_at)::int   AS hour,
    ROUND(AVG(s.quantity)::numeric, 2)  AS quantity
FROM sales s
WHERE s.sold_at >= CURRENT_DATE - INTERVAL '7 days'
  AND s.sold_at <  CURRENT_DATE
GROUP BY s.store_id, s.product_id, EXTRACT(HOUR FROM s.sold_at)
ON CONFLICT (store_id, product_id, forecast_date, hour) DO UPDATE
    SET quantity = EXCLUDED.quantity;
