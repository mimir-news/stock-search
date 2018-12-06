package service

import (
	"github.com/mimir-news/pkg/schema/stock"
	"github.com/mimir-news/stock-search/pkg/repository"
)

// StockService service for interacting with stocks.
type StockService interface {
	RankStocks() error
	RankStock(symbol string) error
	Search(query string, limit int) ([]stock.Stock, error)
}

// NewStockService creates a StockService using the default implementation.
func NewStockService(stockRepo repository.StockRepo) StockService {
	return &stockSvc{
		stockRepo: stockRepo,
	}
}

type stockSvc struct {
	stockRepo repository.StockRepo
}

// Search attempts to match a query against the stored list of stocks.
func (s *stockSvc) Search(query string, limit int) ([]stock.Stock, error) {
	return s.stockRepo.Search(query, limit)
}

// RankStocks counts stock mentions and updates all stocks accordingly.
func (s *stockSvc) RankStocks() error {
	return nil
}

// RankStocks counts a single stocks mentions and updates all it accordingly.
func (s *stockSvc) RankStock(symbol string) error {
	return nil
}
