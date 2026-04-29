package forecast

import (
	"time"

	"github.com/kfc-forecast/internal/core"
)

type repo interface {
	ListStores() ([]core.Store, error)
	AvgSalesForDate(days int, before time.Time) ([]core.ForecastEntry, error)
	UpsertForecasts(entries []core.ForecastEntry) error
	GetForecast(storeID int, date time.Time) ([]core.ForecastEntry, error)
}

type Service struct {
	repo         repo
	lookbackDays int
}

func New(r repo, lookbackDays int) *Service {
	return &Service{repo: r, lookbackDays: lookbackDays}
}

func (s *Service) GenerateForDate(date time.Time) error {
	before := date.Truncate(24 * time.Hour)
	entries, err := s.repo.AvgSalesForDate(s.lookbackDays, before)
	if err != nil {
		return err
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
