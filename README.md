# Sales Forecast

Go + React app that predicts how many units each KFC store will sell per product per hour the next day.

---

## Running the app

Use Docker Compose to containerize and run everything:

```bash
docker compose up --build
```

## How it works

The backend is a Go HTTP server. The core of it is a forecast service that does the prediction logic — it averages the last 7 days (you can config the intervals in the yml) of sales per store, product, and hour, then writes the results to a forecasts table.

Every day at 02:00 a scheduler fires a job that generates predictions for the next day and inserts them into the database. The frontend lets you pick a store and a date to view those predictions as cha
rts.
---

## Database

The app has four tables:

- **stores** — restaurant locations
- **products** 
- **sales** — historical sales records, one row per transaction with a timestamp and quantity
- **forecasts** — the prediction output, one row per store + product + hour + date combination

The forecast table uses a composite primary key on `(store_id, product_id, forecast_date, hour)` so re-running the job for the same date just overwrites the existing rows rather than duplicating them.

---

## What data is available

The seed data covers 7 days of sales history. This means the only date that will have a forecast generated is tomorrow relative to when the scheduler runs — so in practice with this seed you will only see predictions around **April 30** (one day ahead of the seeded data window).

Dates further out won't have forecasts because there is no sales history to average from. If you want to test more dates, you would need to generate additional seed data covering those time ranges.
