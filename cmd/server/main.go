package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	"github.com/kfc-forecast/config"
	"github.com/kfc-forecast/internal/forecast"
	"github.com/kfc-forecast/internal/httphandler"
	"github.com/kfc-forecast/internal/postgres"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	log.Printf("DSN: %s", cfg.DB.DSN)

	pool, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	repo := postgres.New(pool)
	svc := forecast.New(repo, cfg.Forecast.LookbackDays)
	h := httphandler.New(svc)

	c := cron.New()
	_, err = c.AddFunc(cfg.Scheduler.RunAt, func() {
		tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		log.Printf("scheduler: generating forecast for %s", tomorrow.Format("2006-01-02"))
		if genErr := svc.GenerateForDate(tomorrow); genErr != nil {
			log.Printf("scheduler error: %v", genErr)
		}
	})
	if err != nil {
		log.Fatalf("cron setup: %v", err)
	}
	c.Start()
	defer c.Stop()

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Printf("server listening on %s", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, corsMiddleware(mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
