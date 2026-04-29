package forecast

import (
	"time"

	"github.com/kfc-forecast/internal/core"
)

type repo interface {
	ListStores() ([]core.Store, error)
	ListProducts() ([]core.Product, error)
	AvgSalesLastNDays(storeID, productID, days int, before time.Time) (map[int]float64, error)
	UpsertForecasts(entries []core.ForecastEntry) error
	GetForecast(storeID int, date time.Time) ([]core.ForecastEntry, error)
}

type Service struct {
	repo        repo
	lookbackDays int 
}

func New(r repo, lookbackDays int) *Service {
	return &Service{repo: r, lookbackDays: lookbackDays}
}

func (s *Service) GenerateForDate(date time.Time) error {
	stores, err := s.repo.ListStores()
	if err != nil {
		return err
	}
	products, err := s.repo.ListProducts()
	if err != nil {
		return err
	}

	before := date.Truncate(24 * time.Hour)

	var entries []core.ForecastEntry

	for _, store := range stores {
		for _, product := range products {
			hourAvg, err := s.repo.AvgSalesLastNDays(store.ID, product.ID, s.lookbackDays, before)
			if err != nil {
				return err
			}

			for hour, avg := range hourAvg {
				entries = append(entries, core.ForecastEntry{
					StoreID:   store.ID,
					StoreName: store.Name,
					ProductID: product.ID,
					Product:   product.Name,
					Date:      before,
					Hour:      hour,
					Quantity:  avg,
				})
			}
		}
	}

	if len(entries) == 0 {
		return nil 
	}
	return s.repo.UpsertForecasts(entries)
}

func (s *Service) GetForecast(storeID int, date time.Time) ([]core.ForecastEntry, error) {
	return s.repo.GetForecast(storeID, date.Truncate(24*time.Hour))
}

func (s *Service) ListStores() ([]core.Store, error) {
	return s.repo.ListStores()
}
