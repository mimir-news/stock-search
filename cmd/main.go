package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mimir-news/pkg/dbutil"
	"github.com/mimir-news/pkg/httputil"
	"github.com/mimir-news/pkg/httputil/auth"
)

func main() {
	conf := getConfig()
	e := setupEnv(conf)
	defer e.close()
	server := newServer(e, conf)

	log.Printf("Starting %s on port: %s\n", ServiceName, conf.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func newServer(e *env, conf config) *http.Server {
	r := newRouter(e, conf)

	r.GET("/v1/stocks", e.handleStockSearch)
	r.GET("/v1/stocks/suggestions", e.handleSuggestStocks)
	r.PUT("/v1/stocks", e.handleStocksRanking)
	r.PUT("/v1/stocks/:symbol", e.handleStockRanking)

	return &http.Server{
		Addr:    ":" + conf.port,
		Handler: r,
	}
}

func newRouter(e *env, cfg config) *gin.Engine {
	authOpts := auth.NewOptions(cfg.JWTCredentials, unsecuredRoutes...)
	r := httputil.NewRouter(ServiceName, ServiceVersion, e.healthCheck)
	r.Use(auth.RequireToken(authOpts))

	return r
}

func (e *env) healthCheck() error {
	return dbutil.IsConnected(e.db)
}
