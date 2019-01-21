package main

import (
	"database/sql"
	"log"

	"github.com/mimir-news/stock-search/pkg/repository"
	"github.com/mimir-news/stock-search/pkg/service"
)

type env struct {
	db       *sql.DB
	stockSvc service.StockService
}

func setupEnv(cfg config) *env {
	db, err := cfg.db.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	stockRepo := repository.NewStockRepo(db)
	countRepo := repository.NewCountRepo(db)

	return &env{
		db:       db,
		stockSvc: service.NewStockService(stockRepo, countRepo),
	}
}

func (e *env) close() {
	err := e.db.Close()
	if err != nil {
		log.Println(err)
	}
}
