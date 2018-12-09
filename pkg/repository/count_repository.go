package repository

import (
	"database/sql"
	"errors"

	"github.com/mimir-news/stock-search/pkg/domain"
)

// Common errors.
var (
	ErrNoSuchStock = errors.New("No such stock.")
)

// CountRepo handles volume counting of stocks.
type CountRepo interface {
	CountOne(symbol string) (domain.Stock, error)
	CountAll() ([]domain.Stock, error)
}

type pgCountRepo struct {
	db *sql.DB
}

const countStockQuery = `
	SELECT symbol, COUNT(*) FROM tweet_symbol
	WHERE symbol = $1
	GROUP BY symbol`

// CountOne counts the total tweet volume of a single stock.
func (cr *pgCountRepo) CountOne(symbol string) (domain.Stock, error) {
	var s domain.Stock
	err := cr.db.QueryRow(countStockQuery, symbol).Scan(&s.Symbol, &s.Count)
	if err == sql.ErrNoRows {
		return domain.Stock{}, ErrNoSuchStock
	} else if err != nil {
		return domain.Stock{}, err
	}

	return s, nil
}

const countStocksQuery = `
	SELECT symbol, COUNT(*) FROM tweet_symbol
	GROUP BY symbol`

// CountAll counts the total tweet volume of all stocks in the system.
func (cr *pgCountRepo) CountAll() ([]domain.Stock, error) {
	rows, err := cr.db.Query(countStockQuery)
	if err != nil {
		return nil, err
	}

	return mapRowsToCountedStocks(rows)
}

func mapRowsToCountedStocks(rows *sql.Rows) ([]domain.Stock, error) {
	stocks := make([]domain.Stock, 0)

	for rows.Next() {
		var s domain.Stock
		err := rows.Scan(&s.Symbol, &s.Count)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

// MockCountRepo mock implementation of CountRepo.
type MockCountRepo struct {
	CountOneArg   string
	CountOneStock domain.Stock
	CountOneErr   error

	CountAllStocks []domain.Stock
	CountAllErr    error
}

// CountRepo mock CountOne implementation.
func (cr *MockCountRepo) CountRepo(symbol string) (domain.Stock, error) {
	cr.CountOneArg = symbol
	return cr.CountOneStock, cr.CountOneErr
}

// CountAll mock CountAll implementation.
func (cr *MockCountRepo) CountAll() ([]domain.Stock, error) {
	return cr.CountAllStocks, cr.CountAllErr
}
