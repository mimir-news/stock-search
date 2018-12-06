package repository

import (
	"github.com/mimir-news/pkg/schema/stock"
)

// StockRepo handles storing and retrival of stocks.
type StockRepo interface {
	Save(s stock.Stock) error
	Search(query string, limit int) ([]stock.Stock, error)
}

// MockStockRepo mock implementation for stock repo.
type MockStockRepo struct {
	SaveArg stock.Stock
	SaveErr error

	SearchArgQuery string
	SearchArgLimit int
	SearchStocks   []stock.Stock
	SearchErr      error
}

// UnsetArgs sets all repo arguments to their default value.
func (sr *MockStockRepo) UnsetArgs() {
	sr.SaveArg = stock.Stock{}
	sr.SearchArgQuery = ""
	sr.SearchArgLimit = 0
}

// Save mock implementation of saving a stock.
func (sr *MockStockRepo) Save(s stock.Stock) error {
	sr.SaveArg = s
	return sr.SaveErr
}

// Search mock implemntation of searching for stocks.
func (sr *MockStockRepo) Search(query string, limit int) ([]stock.Stock, error) {
	sr.SearchArgQuery = query
	sr.SearchArgLimit = limit
	return sr.SearchStocks, sr.SearchErr
}
