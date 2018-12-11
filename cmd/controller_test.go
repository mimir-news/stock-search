package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mimir-news/pkg/httputil/auth"
	"github.com/mimir-news/pkg/id"
	"github.com/mimir-news/pkg/schema/stock"
	"github.com/mimir-news/stock-search/pkg/domain"
	"github.com/mimir-news/stock-search/pkg/repository"
	"github.com/mimir-news/stock-search/pkg/service"
	"github.com/stretchr/testify/assert"
)

var testAdminID = id.New()

func TestHandleStockSearch(t *testing.T) {
	assert := assert.New(t)

	userID := id.New()
	clientID := id.New()
	query := "A"

	expectedStocks := []domain.Stock{
		domain.Stock{Symbol: "AAPL"},
		domain.Stock{Symbol: "AMD"},
	}

	stockRepo := &repository.MockStockRepo{
		SearchStocks: expectedStocks,
	}

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, stockRepo, nil), conf)
	token := getTestToken(conf, userID, clientID)

	req := createTestGetRequest(clientID, token, "/v1/stocks?query="+query)
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(query, stockRepo.SearchArgQuery)
	assert.Equal(defaultSearchLimit, stockRepo.SearchArgLimit)
	var searchResults []stock.Stock
	err := json.NewDecoder(res.Body).Decode(&searchResults)
	assert.NoError(err)
	assert.Equal(len(expectedStocks), len(searchResults))
	for i, s := range searchResults {
		assert.Equal(expectedStocks[i].Symbol, s.Symbol)
	}

	stockRepo.UnsetArgs()
	req = createTestGetRequest(clientID, token, "/v1/stocks?limit=5&query="+query)
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(query, stockRepo.SearchArgQuery)
	assert.Equal(5, stockRepo.SearchArgLimit)

	stockRepo.UnsetArgs()
	req = createTestGetRequest(clientID, token, "/v1/stocks?limit=5")
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	assert.Equal("", stockRepo.SearchArgQuery)
	assert.Equal(0, stockRepo.SearchArgLimit)

	stockRepo.UnsetArgs()
	stockRepo.SearchErr = errors.New("mock error")
	req = createTestGetRequest(clientID, token, "/v1/stocks?query="+query)
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
	assert.Equal(query, stockRepo.SearchArgQuery)
	assert.Equal(defaultSearchLimit, stockRepo.SearchArgLimit)

}

func TestHandleSuggestStocks(t *testing.T) {
	assert := assert.New(t)

	userID := id.New()
	clientID := id.New()

	expectedStocks := []domain.Stock{
		domain.Stock{Symbol: "AAPL"},
		domain.Stock{Symbol: "AMD"},
	}

	stockRepo := &repository.MockStockRepo{
		FindMostCommonStocks: expectedStocks,
	}

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, stockRepo, nil), conf)
	token := getTestToken(conf, userID, clientID)

	req := createTestGetRequest(clientID, token, "/v1/stocks/suggestions?exclude=A,B")
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, stockRepo.FindMostCommonInvocations)
	assert.Equal(defaultSuggestionLimit, stockRepo.FindMostCommonArgLimit)
	var suggestions []stock.Stock
	err := json.NewDecoder(res.Body).Decode(&suggestions)
	assert.NoError(err)
	assert.Equal(len(expectedStocks), len(suggestions))
	for i, s := range suggestions {
		assert.Equal(expectedStocks[i].Symbol, s.Symbol)
	}

	stockRepo.UnsetArgs()
	req = createTestGetRequest(clientID, token, "/v1/stocks/suggestions?exclude=A,B&limit=10")
	res = performTestRequest(server.Handler, req)

	assert.Equal(1, stockRepo.FindMostCommonInvocations)
	assert.Equal(10, stockRepo.FindMostCommonArgLimit)
	err = json.NewDecoder(res.Body).Decode(&suggestions)
	assert.NoError(err)
	assert.Equal(len(expectedStocks), len(suggestions))
	for i, s := range suggestions {
		assert.Equal(expectedStocks[i].Symbol, s.Symbol)
	}
}

func TestHandleStockRanking(t *testing.T) {
	assert := assert.New(t)

	clientID := id.New()
	symbol := "AAPL"

	coutedStock := domain.Stock{Symbol: symbol, Count: 10}

	stockRepo := &repository.MockStockRepo{}
	countRepo := &repository.MockCountRepo{
		CountOneStock: coutedStock,
	}

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, stockRepo, countRepo), conf)
	token := getTestToken(conf, testAdminID, clientID)

	req := createTestPutRequest(clientID, token, "/v1/stocks/"+symbol)
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(symbol, countRepo.CountOneArg)
	savedStock := stockRepo.SaveArg
	assert.Equal(symbol, savedStock.Symbol)
	assert.Equal(coutedStock.Count, savedStock.Count)

	countRepo.CountOneErr = repository.ErrNoSuchStock
	stockRepo.UnsetArgs()
	req = createTestPutRequest(clientID, token, "/v1/stocks/MISSING")
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusNotFound, res.Code)
	assert.Equal("MISSING", countRepo.CountOneArg)
	savedStock = stockRepo.SaveArg
	assert.Equal("", savedStock.Symbol)
	assert.Equal(int64(0), savedStock.Count)

	countRepo.UnsetArgs()
	wrongToken := getTestToken(conf, id.New(), clientID)
	req = createTestPutRequest(clientID, wrongToken, "/v1/stocks/MISSING")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusForbidden, res.Code)
	assert.Equal(0, countRepo.CountOneInvocations)

}

func TestHandleStocksRanking(t *testing.T) {
	assert := assert.New(t)

	clientID := id.New()

	coutedStocks := []domain.Stock{
		domain.Stock{Symbol: "AAPL", Count: 10},
		domain.Stock{Symbol: "GOOG", Count: 20},
	}

	stockRepo := &repository.MockStockRepo{}
	countRepo := &repository.MockCountRepo{
		CountAllStocks: coutedStocks,
	}

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, stockRepo, countRepo), conf)
	token := getTestToken(conf, testAdminID, clientID)

	req := createTestPutRequest(clientID, token, "/v1/stocks")
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, countRepo.CountAllInvocations)
	savedStock := stockRepo.SaveArg
	assert.Equal(len(coutedStocks), stockRepo.SaveInvocations)
	assert.Equal("GOOG", savedStock.Symbol)
	assert.Equal(int64(20), savedStock.Count)

	countRepo.UnsetArgs()
	wrongToken := getTestToken(conf, id.New(), clientID)
	req = createTestPutRequest(clientID, wrongToken, "/v1/stocks")
	res = performTestRequest(server.Handler, req)
	assert.Equal(http.StatusForbidden, res.Code)
	assert.Equal(0, countRepo.CountAllInvocations)

}

func performTestRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getTestEnv(conf config, stockRepo repository.StockRepo, countRepo repository.CountRepo) *env {
	return &env{
		stockSvc: service.NewStockService(stockRepo, countRepo),
		adminID:  conf.adminID,
	}
}

func getTestConfig() config {
	return config{
		tokenSecret:          "my-secret",
		tokenVerificationKey: "my-verification-key",
		adminID:              testAdminID,
	}
}

func getTestSigner(conf config) auth.Signer {
	return auth.NewSigner(conf.tokenSecret, conf.tokenVerificationKey, 24*time.Hour)
}

func getTestToken(conf config, userID, clientID string) string {
	signer := getTestSigner(conf)

	token, err := signer.New(id.New(), userID, clientID)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func createTestPutRequest(clientID, token, route string) *http.Request {
	return createTestRequest(clientID, token, route, http.MethodPut)
}

func createTestGetRequest(clientID, token, route string) *http.Request {
	return createTestRequest(clientID, token, route, http.MethodGet)
}

func createTestRequest(clientID, token, route, method string) *http.Request {
	req, err := http.NewRequest(method, route, nil)
	if err != nil {
		log.Fatal(err)
	}

	if clientID != "" {
		req.Header.Set(auth.ClientIDKey, clientID)
	}
	if token != "" {
		bearerToken := auth.AuthTokenPrefix + token
		req.Header.Set(auth.AuthHeaderKey, bearerToken)
	}

	return req
}
