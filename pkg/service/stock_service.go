package service

import (
	"net/http"

	"github.com/mimir-news/pkg/httputil"
	"github.com/mimir-news/pkg/schema/stock"
	"github.com/mimir-news/stock-search/pkg/domain"
	"github.com/mimir-news/stock-search/pkg/repository"
)

// StockService service for interacting with stocks.
type StockService interface {
	RankStocks() error
	RankStock(symbol string) error
	Search(query string, limit int) ([]stock.Stock, error)
	GetSuggestions(excluded []string, limit int) ([]stock.Stock, error)
}

// NewStockService creates a StockService using the default implementation.
func NewStockService(stockRepo repository.StockRepo, countRepo repository.CountRepo) StockService {
	return &stockSvc{
		stockRepo: stockRepo,
		countRepo: countRepo,
	}
}

type stockSvc struct {
	stockRepo repository.StockRepo
	countRepo repository.CountRepo
}

// Search attempts to match a query against the stored list of stocks.
func (svc *stockSvc) Search(query string, limit int) ([]stock.Stock, error) {
	stocks, err := svc.stockRepo.Search(query, limit)
	if err != nil {
		return nil, err
	}

	return mapStocksToDTOs(stocks), nil
}

// RankStocks counts stock mentions and updates all stocks accordingly.
func (svc *stockSvc) RankStocks() error {
	countedStocks, err := svc.countRepo.CountAll()
	if err != nil {
		return err
	}

	for _, s := range countedStocks {
		err := svc.stockRepo.Save(s)
		if err != nil {
			return err
		}
	}

	return nil
}

// RankStocks counts a single stocks mentions and updates all it accordingly.
func (svc *stockSvc) RankStock(symbol string) error {
	s, err := svc.countRepo.CountOne(symbol)
	if err == repository.ErrNoSuchStock {
		return httputil.NewError(err.Error(), http.StatusNotFound)
	} else if err != nil {
		return err
	}

	return svc.stockRepo.Save(s)
}

// GetSuggestions gets most common stocks except the specified excluded.
func (svc *stockSvc) GetSuggestions(excluded []string, limit int) ([]stock.Stock, error) {
	stocks, err := svc.stockRepo.FindMostCommon(excluded, limit)
	if err != nil {
		return nil, err
	}

	return mapStocksToDTOs(stocks), nil
}

func mapStocksToDTOs(stocks []domain.Stock) []stock.Stock {
	dtos := make([]stock.Stock, 0, len(stocks))
	for _, s := range stocks {
		dtos = append(dtos, s.ToDTO())
	}

	return dtos
}
