package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mimir-news/pkg/httputil"
)

func (e *env) handleStockSearch(c *gin.Context) {
	query, err := httputil.ParseQueryValue(c, "query")
	if err != nil {
		c.Error(err)
		return
	}

	searchLimit, err := getIntParam(c, "limit", defaultSearchLimit)
	if err != nil {
		c.Error(err)
		return
	}

	results, err := e.stockSvc.Search(query, searchLimit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (e *env) handleSuggestStocks(c *gin.Context) {
	excluded := getSymbolsFromQuery(c, "exclude")
	limit, err := getIntParam(c, "limit", defaultSuggestionLimit)
	if err != nil {
		c.Error(err)
		return
	}

	results, err := e.stockSvc.GetSuggestions(excluded, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (e *env) handleStocksRanking(c *gin.Context) {
	err := e.stockSvc.RankStocks()
	if err != nil {
		c.Error(err)
		return
	}

	httputil.SendOK(c)
}

func (e *env) handleStockRanking(c *gin.Context) {
	symbol := c.Param("symbol")
	err := e.stockSvc.RankStock(symbol)
	if err != nil {
		c.Error(err)
		return
	}

	httputil.SendOK(c)
}

func getIntParam(c *gin.Context, name string, defaultValue int) (int, error) {
	value, ok := c.GetQuery(name)
	if !ok {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, httputil.NewError("Invalid "+name, http.StatusBadRequest)
	}

	return intValue, nil
}

func getSymbolsFromQuery(c *gin.Context, name string) []string {
	symbols, ok := c.GetQueryArray(name)
	if !ok {
		return []string{}
	}

	return symbols
}
