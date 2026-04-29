package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-forecast/internal/core"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (r *Repository) ListStores() ([]core.Store, error) {
	c, cancel := ctx()
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

func (r *Repository) ListProducts() ([]core.Product, error) {
	c, cancel := ctx()
	defer cancel()
	rows, err := r.pool.Query(c, `SELECT id, name FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []core.Product
	for rows.Next() {
		var p core.Product
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *Repository) AvgSalesLastNDays(storeID, productID, days int, before time.Time) (map[int]float64, error) {
	c, cancel := ctx()
	defer cancel()
	const q = `
		SELECT EXTRACT(HOUR FROM sold_at)::int AS hour, AVG(quantity) AS avg_qty
		FROM sales
		WHERE store_id = $1 AND product_id = $2
		  AND sold_at >= $3 AND sold_at < $4
		GROUP BY hour ORDER BY hour
	`
	from := before.AddDate(0, 0, -days)
	rows, err := r.pool.Query(c, q, storeID, productID, from, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[int]float64, 24)
	for rows.Next() {
		var hour int
		var avg float64
		if err := rows.Scan(&hour, &avg); err != nil {
			return nil, err
		}
		result[hour] = avg
	}
	return result, rows.Err()
}

func (r *Repository) UpsertForecasts(entries []core.ForecastEntry) error {
	c, cancel := ctx()
	defer cancel()
	const q = `
		INSERT INTO forecasts (store_id, product_id, forecast_date, hour, quantity)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (store_id, product_id, forecast_date, hour)
		DO UPDATE SET quantity = EXCLUDED.quantity
	`
	tx, err := r.pool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)
	for _, e := range entries {
		if _, err := tx.Exec(c, q, e.StoreID, e.ProductID, e.Date, e.Hour, e.Quantity); err != nil {
			return err
		}
	}
	return tx.Commit(c)
}

func (r *Repository) GetForecast(storeID int, date time.Time) ([]core.ForecastEntry, error) {
	c, cancel := ctx()
	defer cancel()
	const q = `
		SELECT f.store_id, s.name, f.product_id, p.name, f.forecast_date, f.hour, f.quantity
		FROM forecasts f
		JOIN stores s ON s.id = f.store_id
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
