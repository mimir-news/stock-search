package domain

import (
	"github.com/mimir-news/pkg/schema/stock"
)

// Stock holds stock data.
type Stock struct {
	Name   string
	Symbol string
	Count  int64
}

// NewDomainStock converts a stock to the internal domain structure.
func NewDomainStock(s stock.Stock, count int64) Stock {
	return Stock{
		Name:   s.Name,
		Symbol: s.Symbol,
		Count:  count,
	}
}

// ToDTO converts a stock to a DTO.
func (s *Stock) ToDTO() stock.Stock {
	return stock.Stock{
		Name:   s.Name,
		Symbol: s.Symbol,
	}
}
