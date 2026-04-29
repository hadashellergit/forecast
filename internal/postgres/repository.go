package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-forecast/internal/core"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func readCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func writeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (r *Repository) ListStores() ([]core.Store, error) {
	c, cancel := readCtx()
	defer cancel()
	rows, err := r.pool.Query(c, `SELECT id, name, city FROM stores ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stores []core.Store
	for rows.Next() {
		var s core.Store
		if err := rows.Scan(&s.ID, &s.Name, &s.City); err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}
	return stores, rows.Err()
}

func (r *Repository) AvgSalesForDate(days int, before time.Time) ([]core.ForecastEntry, error) {
	c, cancel := writeCtx()
	defer cancel()
	const q = `
		SELECT s.id, s.name, p.id, p.name,
		       EXTRACT(HOUR FROM sl.sold_at)::int AS hour,
		       AVG(sl.quantity) AS avg_qty
		FROM sales sl
		JOIN stores  s ON s.id = sl.store_id
		JOIN products p ON p.id = sl.product_id
		WHERE sl.sold_at >= $1 AND sl.sold_at < $2
		GROUP BY s.id, s.name, p.id, p.name, hour
		ORDER BY s.id, p.id, hour
	`
	from := before.AddDate(0, 0, -days)
	rows, err := r.pool.Query(c, q, from, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []core.ForecastEntry
	for rows.Next() {
		var e core.ForecastEntry
		e.Date = before
		if err := rows.Scan(&e.StoreID, &e.StoreName, &e.ProductID, &e.Product, &e.Hour, &e.Quantity); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *Repository) UpsertForecasts(entries []core.ForecastEntry) error {
	if len(entries) == 0 {
		return nil
	}
	c, cancel := writeCtx()
	defer cancel()
	const q = `
		INSERT INTO forecasts (store_id, product_id, forecast_date, hour, quantity)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (store_id, product_id, forecast_date, hour)
		DO UPDATE SET quantity = EXCLUDED.quantity
	`
	batch := &pgx.Batch{}
	for _, e := range entries {
		batch.Queue(q, e.StoreID, e.ProductID, e.Date, e.Hour, e.Quantity)
	}

	tx, err := r.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	br := tx.SendBatch(c, batch)
	for range entries {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(c)
}

func (r *Repository) GetForecast(storeID int, date time.Time) ([]core.ForecastEntry, error) {
	c, cancel := readCtx()
	defer cancel()
	const q = `
		SELECT f.store_id, s.name, f.product_id, p.name, f.forecast_date, f.hour, f.quantity
		FROM forecasts f
		JOIN stores   s ON s.id = f.store_id
		JOIN products p ON p.id = f.product_id
		WHERE f.store_id = $1 AND f.forecast_date = $2
		ORDER BY f.hour, p.name
	`
	rows, err := r.pool.Query(c, q, storeID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []core.ForecastEntry
	for rows.Next() {
		var e core.ForecastEntry
		if err := rows.Scan(&e.StoreID, &e.StoreName, &e.ProductID, &e.Product, &e.Date, &e.Hour, &e.Quantity); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
