package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/mimir-news/pkg/dbutil"
	"github.com/mimir-news/stock-search/pkg/domain"
)

var (
	errInsertStockFailed = errors.New("Inserting stock stock failed")
)

// StockRepo handles storing and retrival of stocks.
type StockRepo interface {
	Save(s domain.Stock) error
	Search(query string, limit int) ([]domain.Stock, error)
	FindMostCommon(excluded []string, limit int) ([]domain.Stock, error)
}

// NewStockRepo created a StockRepo using the default implementation.
func NewStockRepo(db *sql.DB) StockRepo {
	return &pgStockRepo{
		db: db,
	}
}

// pgStockRepo postgres implementation of StockRepo.
type pgStockRepo struct {
	db *sql.DB
}

const saveStockQuery = `
	INSERT INTO stock(symbol, name, is_active, total_count, updated_at)
	VALUES($1, $2, TRUE, $3, $4) ON CONFLICT ON CONSTRAINT stock_pkey 
	DO UPDATE SET total_count = $3, updated_at = $4`

// Save saves a stock.
func (pg *pgStockRepo) Save(s domain.Stock) error {
	res, err := pg.db.Exec(saveStockQuery, s.Symbol, s.Name, s.Count, time.Now().UTC())
	if err != nil {
		return errInsertStockFailed
	}

	return dbutil.AssertRowsAffected(res, 1, errInsertStockFailed)
}

const searchStockQuery = `
	SELECT symbol, name, total_count FROM stock 
	WHERE is_active = TRUE 
	AND (
		LOWER(symbol) LIKE $1 || '%' OR
		LOWER(name) LIKE $1 || '%'
	)
	ORDER BY total_count DESC
	LIMIT $2`

// Search finds stocks mathing a given query.
func (pg *pgStockRepo) Search(query string, limit int) ([]domain.Stock, error) {
	lowerQuery := strings.ToLower(query)
	rows, err := pg.db.Query(searchStockQuery, lowerQuery, limit)
	if err != nil {
		return nil, err
	}

	return mapRowsToStocks(rows)
}

// FindMostCommon finds the most common stocks
// except the ones that contains the symbols provided.
func (pg *pgStockRepo) FindMostCommon(excluded []string, limit int) ([]domain.Stock, error) {
	return nil, nil
}

func mapRowsToStocks(rows *sql.Rows) ([]domain.Stock, error) {
	stocks := make([]domain.Stock, 0)

	for rows.Next() {
		var s domain.Stock
		err := rows.Scan(&s.Symbol, &s.Name, &s.Count)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

// MockStockRepo mock implementation for StockRepo.
type MockStockRepo struct {
	SaveArg         domain.Stock
	SaveErr         error
	SaveInvocations int

	SearchArgQuery    string
	SearchArgLimit    int
	SearchStocks      []domain.Stock
	SearchErr         error
	SearchInvocations int

	FindMostCommonArgExcluded []string
	FindMostCommonArgLimit    int
	FindMostCommonStocks      []domain.Stock
	FindMostCommonErr         error
	FindMostCommonInvocations int
}

// UnsetArgs sets all repo arguments to their default value.
func (sr *MockStockRepo) UnsetArgs() {
	sr.SaveArg = domain.Stock{}
	sr.SaveInvocations = 0

	sr.SearchArgQuery = ""
	sr.SearchArgLimit = 0
	sr.SearchInvocations = 0

	sr.FindMostCommonArgExcluded = nil
	sr.FindMostCommonArgLimit = 0
	sr.FindMostCommonInvocations = 0
}

// Save mock implementation of saving a stock.
func (sr *MockStockRepo) Save(s domain.Stock) error {
	sr.SaveArg = s
	sr.SaveInvocations++
	return sr.SaveErr
}

// Search mock implemntation of searching for stocks.
func (sr *MockStockRepo) Search(query string, limit int) ([]domain.Stock, error) {
	sr.SearchArgQuery = query
	sr.SearchArgLimit = limit
	sr.SearchInvocations++
	return sr.SearchStocks, sr.SearchErr
}

// FindMostCommon mock implementation of finding common stocks.
func (sr *MockStockRepo) FindMostCommon(excluded []string, limit int) ([]domain.Stock, error) {
	sr.FindMostCommonArgExcluded = excluded
	sr.FindMostCommonArgLimit = limit
	sr.FindMostCommonInvocations++
	return sr.FindMostCommonStocks, sr.FindMostCommonErr
}
