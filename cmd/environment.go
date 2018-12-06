package main

import (
	"database/sql"
	"log"

	"github.com/mimir-news/stock-search/pkg/service"
)

type env struct {
	db       *sql.DB
	stockSvc service.StockService
	adminID  string
}

func setupEnv(conf config) *env {
	db, err := conf.db.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	return &env{
		db:       db,
		stockSvc: service.NewStockService(nil),
		adminID:  conf.adminID,
	}
}

func (e *env) close() {
	err := e.db.Close()
	if err != nil {
		log.Println(err)
	}
}
