package main

import (
	"log"
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

	searchLimit, err := getSearchResultLimit(c)
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

func (e *env) handleStocksRanking(c *gin.Context) {
	httputil.SendOK(c)
}

func (e *env) handleStockRanking(c *gin.Context) {
	httputil.SendOK(c)
}

func getSearchResultLimit(c *gin.Context) (int, error) {
	value, ok := c.GetQuery("limit")
	if !ok {
		return defaultSearchLimit, nil
	}

	searchLimit, err := strconv.Atoi(value)
	if err != nil {
		log.Panicln(err)
		return 0, httputil.NewError("Invalid limit", http.StatusBadRequest)
	}

	return searchLimit, nil
}
