package core

import "time"

type Store struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}


type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ForecastEntry struct {
	StoreID   int       `json:"store_id"`
	StoreName string    `json:"store_name"`
	ProductID int       `json:"product_id"`
	Product   string    `json:"product_name"`
	Date      time.Time `json:"date"`
	Hour      int       `json:"hour"`      
	Quantity  float64   `json:"quantity"`  
}

type ForecastService interface {

	GenerateForDate(date time.Time) error

	GetForecast(storeID int, date time.Time) ([]ForecastEntry, error)

	ListStores() ([]Store, error)
}
